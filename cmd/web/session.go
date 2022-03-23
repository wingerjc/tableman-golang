package web

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

var timeNow = time.Now

type SessionSet struct {
	sync.RWMutex
	maxSessions int
	maxAge      time.Duration
	sessions    map[string]*Session
}

func NewSessionSet(maxSessions int, maxAge time.Duration) *SessionSet {
	return &SessionSet{
		maxSessions: maxSessions,
		maxAge:      maxAge,
		sessions:    make(map[string]*Session),
	}
}

func (ss *SessionSet) NewSession() string {
	ss.Lock()
	defer ss.Unlock()

	// Prune old sessions or the oldest session in use.
	if len(ss.sessions) >= ss.maxSessions {
		pruneTime := timeNow().Add(-1 * ss.maxAge)
		prune := make([]string, 0)
		min := ""
		minTime := timeNow()
		for k, v := range ss.sessions {
			t := v.Accessed()
			if minTime.After(t) {
				min = k
				minTime = t
			}
			if pruneTime.After(t) {
				prune = append(prune, k)
			}
		}
		for _, k := range prune {
			delete(ss.sessions, k)
		}
		if len(ss.sessions) >= ss.maxSessions && len(min) > 0 {
			delete(ss.sessions, min)
		}
	}

	k := uuid.New().String()
	ss.sessions[k] = NewSession()
	return k
}

func (ss *SessionSet) AddPack(sid string, key string, pack *program.Program) error {
	ss.RLock()
	defer ss.RUnlock()
	s, ok := ss.sessions[sid]
	if !ok {
		return fmt.Errorf("invalid session ID %s", sid)
	}
	s.AddPack(key, pack)
	return nil
}

func (ss *SessionSet) Eval(sid string, key string, expr program.Evallable) (string, error) {
	ss.RLock()
	defer ss.RUnlock()
	s, ok := ss.sessions[sid]
	if !ok {
		return "", fmt.Errorf("invalid session ID %s", sid)
	}
	return s.Eval(key, expr)
}

func (ss *SessionSet) Contains(sid string) bool {
	ss.RLock()
	defer ss.RUnlock()
	_, ok := ss.sessions[sid]
	return ok
}

type Session struct {
	accessMu sync.Mutex
	accessed time.Time
	packMu   sync.RWMutex
	packs    map[string]*program.Program
	history  *program.RollHistory
}

func NewSession() *Session {
	return &Session{
		accessed: timeNow(),
		packs:    make(map[string]*program.Program),
		history:  program.NewRollHistory(),
	}
}

func (s *Session) Touch() {
	s.accessMu.Lock()
	defer s.accessMu.Unlock()
	s.accessed = timeNow()
}

func (s *Session) Accessed() time.Time {
	s.accessMu.Lock()
	defer s.accessMu.Unlock()
	return s.accessed
}

func (s *Session) AddPack(key string, pack *program.Program) {
	s.Touch()
	s.packMu.Lock()
	pack.SetHistory(s.history)
	defer s.packMu.Unlock()
	s.packs[key] = pack
}

func (s *Session) Eval(packKey string, expr program.Evallable) (string, error) {
	s.Touch()
	s.packMu.RLock()
	defer s.packMu.RUnlock()
	p, ok := s.packs[packKey]
	if !ok {
		return "", fmt.Errorf("table set named %s not loaded", packKey)
	}

	res, err := p.Eval(expr)
	if err != nil {
		return "", err
	}
	if res.MatchType(program.IntResult) {
		return fmt.Sprintf("%d", res.IntVal()), nil
	}
	return res.StringVal(), nil
}
