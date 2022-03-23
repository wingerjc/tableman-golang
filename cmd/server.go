package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wingerjc/tableman-golang/cmd/web"
	compiler "github.com/wingerjc/tableman-golang/pkg/compile"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

const (
	SessionIDCookie = "sessionId"
)

type ServerConfig struct {
	Port           string
	CertFile       string
	KeyFile        string
	packConfigPath string
	staticFilePath string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: ":8080",
	}
}

type Server struct {
	cfg         *ServerConfig
	server      *http.Server
	packs       map[string]*program.Program
	loadedPacks []*web.PackEntry
	compiler    *compiler.Compiler
	sessions    *web.SessionSet
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = cfg.Port
	} else {
		port = ":" + port
	}
	result := &Server{
		cfg: cfg,
		server: &http.Server{
			Addr: port,
		},
		sessions: web.NewSessionSet(5000, 2*time.Hour),
		packs:    make(map[string]*program.Program),
	}

	// Register all the endpoints.
	result.register()

	// Load the table packs
	data, err := os.ReadFile(cfg.packConfigPath)
	if err != nil {
		return nil, err
	}
	packs := &web.PackWebConfig{}
	if err := json.Unmarshal(data, packs); err != nil {
		return nil, err
	}
	err = result.load(packs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Server) load(packs *web.PackWebConfig) error {
	s.loadedPacks = make([]*web.PackEntry, 0)
	if s.compiler == nil {
		comp, err := compiler.NewCompiler()
		s.compiler = comp
		if err != nil {
			return err
		}
	}

	for _, p := range packs.Packs {
		prog, err := s.compiler.CompileFile(p.Path)
		if err != nil {
			return err
		}
		s.packs[p.Name] = prog
		s.loadedPacks = append(s.loadedPacks, p)
	}

	return nil
}

func (s *Server) register() {
	mux := http.DefaultServeMux
	// Register routes
	mux.HandleFunc("/hello", s.hello())
	mux.HandleFunc("/session", s.handleSession())
	mux.HandleFunc("/pack", s.handlePacks())
	mux.HandleFunc("/eval", s.handleEval())

	if len(s.cfg.staticFilePath) > 0 {
		pathStr, _ := filepath.Abs(s.cfg.staticFilePath)
		fmt.Printf("Serving static files from %s\n", pathStr)
		fs := http.FileServer(http.Dir(pathStr))
		mux.Handle("/site/", http.StripPrefix("/site/", fs))
	}

	s.server.Handler = mux
}

func (s *Server) hello() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte("Hello World!"))
	}
}

func (s *Server) handleSession() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sid := s.sessionAuth(rw, r)
			res, err := json.Marshal(web.SessionIdentifierDTO{ID: sid})
			if err != nil {
				errOut(rw, err)
				return
			}
			rw.Write(res)
		default:
			rw.WriteHeader(405)
		}
	}
}

func (s *Server) handlePacks() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		// GET returns a list of all packs
		case http.MethodGet:
			result := make([]*web.PackDefDTO, 0)
			for _, k := range s.loadedPacks {
				result = append(result, web.NewPackDefDTO(k.Name, k.Title))
			}
			jsonRes, err := json.Marshal(result)
			if err != nil {
				errOut(rw, err)
				return
			}
			rw.Write(jsonRes)
		// PUT adds the given pack to the current session, idempotent.
		case http.MethodPut:
			sid := s.sessionAuth(rw, r)
			req := &web.LoadPackDTO{}
			if err := decode(r, req); err != nil {
				errOut(rw, err)
				return
			}
			if !s.LoadPack(rw, sid, req.Pack) {
				return
			}
			rw.WriteHeader(200)
		default:
			rw.WriteHeader(405)
		}
	}
}

func (s *Server) handleEval() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			sid := s.sessionAuth(rw, r)
			req := &web.EvalDTO{}
			if err := decode(r, req); err != nil {
				errOut(rw, err)
				return
			}
			result := &web.EvalResultDTO{
				EvalDTO: req,
			}
			expr, err := s.compiler.CompileExpression(req.Expr)
			if err != nil {
				result.CompileError = err.Error()
				rw.WriteHeader(400)
				data, _ := json.Marshal(result)
				rw.Write(data)
				return
			}
			if !s.LoadPack(rw, sid, req.Pack) {
				return
			}
			res, err := s.sessions.Eval(sid, req.Pack, expr)
			if err != nil {
				result.RuntimeError = err.Error()
				rw.WriteHeader(500)
				data, _ := json.Marshal(result)
				rw.Write(data)
				return
			}
			result.Result = res
			rw.WriteHeader(200)
			data, _ := json.Marshal(result)
			rw.Write(data)
		default:
			rw.WriteHeader(405)
		}
	}
}

func (s *Server) LoadPack(rw http.ResponseWriter, sid string, pack string) bool {
	prog, ok := s.packs[pack]
	if !ok {
		rw.WriteHeader(404)
		errFmt(rw, fmt.Sprintf("No pack set named %s", pack))
		return false
	}
	if err := s.sessions.AddPack(sid, pack, prog.Copy()); err != nil {
		errOut(rw, err)
		return false
	}
	return true
}

// sessionAuth will create a new session and store it in the cookie if one does not exist.
func (s *Server) sessionAuth(rw http.ResponseWriter, r *http.Request) string {
	// Check
	c, err := r.Cookie(SessionIDCookie)
	if err == nil && s.sessions.Contains(c.Value) {
		return c.Value
	}

	sid := s.sessions.NewSession()
	http.SetCookie(rw, &http.Cookie{
		Name:  SessionIDCookie,
		Value: sid,
	})
	return sid
}

func errOut(rw http.ResponseWriter, err error) {
	rw.WriteHeader(500)
	errFmt(rw, err.Error())
}

func errFmt(rw http.ResponseWriter, msg string) {
	out, _ := json.Marshal(web.ErrorDTO{Error: msg})
	rw.Write(out)
}

func decode(r *http.Request, i interface{}) error {
	return json.NewDecoder(r.Body).Decode(i)
}

func (s *Server) Run() error {
	fmt.Println("Running web server")
	if len(s.cfg.CertFile) > 0 && len(s.cfg.KeyFile) > 0 {
		return s.server.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
	}
	return s.server.ListenAndServe()
}
