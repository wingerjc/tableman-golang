package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/wingerjc/tableman-golang/cmd/web"
	compiler "github.com/wingerjc/tableman-golang/pkg/compile"
)

type ServerConfig struct {
	Port      string
	CertFile  string
	KeyFile   string
	PackFiles map[string]string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:      ":8080",
		PackFiles: make(map[string]string),
	}
}

type Server struct {
	cfg      *ServerConfig
	server   *http.Server
	compiler *compiler.Compiler
	sessions *web.SessionSet
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	result := &Server{
		cfg: cfg,
		server: &http.Server{
			Addr: cfg.Port,
		},
		sessions: web.NewSessionSet(5000, 2*time.Hour),
	}

	// Register all the endpoints.
	result.register()

	// Load the table packs
	err := result.load()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Server) load() error {
	if s.compiler == nil {
		comp, err := compiler.NewCompiler()
		s.compiler = comp
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) register() {
	mux := http.DefaultServeMux
	// Register routes
	mux.HandleFunc("/hello", s.hello())
	mux.HandleFunc("/session", s.handleSession())

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
			sid := s.sessions.NewSession()
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

func errOut(rw http.ResponseWriter, err error) {
	out, _ := json.Marshal(web.ErrorDTO{Error: err.Error()})
	rw.WriteHeader(500)
	rw.Write(out)
}

func (s *Server) Run() error {
	if len(s.cfg.CertFile) > 0 && len(s.cfg.KeyFile) > 0 {
		return s.server.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
	}
	return s.server.ListenAndServe()
}
