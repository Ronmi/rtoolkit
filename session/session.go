package session

import (
	"net/http"
	"time"

	"github.com/Ronmi/rtoolkit/session/store"
)

// CookieMaker makes a cookie for session, ttl < 0 means it is expired
type CookieMaker func(name, value string, ttl int) *http.Cookie

// DefaultCookieMaker is default cookie maker, forcing httponly
func DefaultCookieMaker(name, value string, ttl int) *http.Cookie {
	ret := &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(time.Duration(ttl) * time.Second),
		MaxAge:   ttl,
		HttpOnly: true,
	}

	return ret
}

// Manager is main session class
type Manager struct {
	Store store.Store // default to store.InMemory
	Key   string      // default to "SESSION_ID", this is used in cookie
	TTL   int         // default to 7200 (2 hours)
}

func (m *Manager) init() {
	if m.TTL <= 0 {
		m.TTL = 7200
	}

	if m.Key == "" {
		m.Key = "SESSION_ID"
	}

	if m.Store == nil {
		m.Store = store.InMemory(m.TTL)
	}
}

// Start begins or resumes a session
func (m *Manager) Start(r *http.Request) (sess *Session, err error) {
	m.init()
	c, err := r.Cookie(m.Key)

	if err != nil || c.Value == "" {
		return newSession(m)
	}

	if _, err = m.Store.Get(c.Value); err != nil {
		return newSession(m)
	}

	return loadSession(c.Value, m)
}

// Session represents session for a specific client
type Session struct {
	id      string
	data    string
	expired bool
	saved   bool
	m       *Manager
}

func newSession(m *Manager) (*Session, error) {
	id, err := m.Store.Allocate()
	if err != nil {
		return nil, err
	}

	return &Session{
		id: id,
		m:  m,
	}, nil
}

func loadSession(id string, m *Manager) (*Session, error) {
	data, err := m.Store.Get(id)
	if err != nil {
		return nil, err
	}

	return &Session{
		id:   id,
		data: data,
		m:    m,
	}, nil
}

// ID returns session id
func (s *Session) ID() string {
	return s.id
}

// Destroy sets this session as expired
func (s *Session) Destroy() {
	s.expired = true
}

// Save saves data, updates cookie expire time, and delete cookie or session if expired
// Since it sets cookie, you have to call it before w.Write()
func (s *Session) Save(w http.ResponseWriter, maker CookieMaker) error {
	if s.expired {
		s.m.Store.Release(s.id)
		c := maker(s.m.Key, "", s.m.TTL)
		http.SetCookie(w, c)
		s.saved = true
		return nil
	}

	err := s.m.Store.Set(s.id, s.data)
	if err == nil {
		c := maker(s.m.Key, s.id, s.m.TTL)
		http.SetCookie(w, c)
		s.saved = true
	}
	return err
}

// Data returns session data
func (s *Session) Data() string {
	return s.data
}

// SetDate updates session data, yu need Save() to save it into session storage
func (s *Session) SetData(data string) error {
	err := s.m.Store.Set(s.id, data)
	if err == nil {
		s.saved = false
		s.expired = false
		s.data = data
	}

	return err
}
