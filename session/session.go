package session

import (
	"errors"
	"math/rand"
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

func generateSeed() string {
	size := store.SeedLength
	arr := []byte(store.SeedChars)
	l := len(arr)

	ret := make([]byte, size)

	for i := 0; i < size; i++ {
		ret[i] = arr[rand.Intn(l)]
	}

	return string(ret)
}

// Manager is main session class
type Manager struct {
	Store       store.Store // default to store.InMemory
	Key         string      // default to "SESSION_ID", this is used in cookie
	ChecksumKey string      // default to "SESSION_CHECK", this is used in cookie
	TTL         int         // default to 7200 (2 hours)
	MakeCookie  CookieMaker // default to DefaultCookieMaker
}

func (m *Manager) init() {
	if m.TTL <= 0 {
		m.TTL = 7200
	}

	if m.Key == "" {
		m.Key = "SESSION_ID"
	}

	if m.ChecksumKey == "" {
		m.Key = "SESSION_CHECK"
	}

	if m.Store == nil {
		m.Store = store.InMemory(m.TTL)
	}

	if m.MakeCookie == nil {
		m.MakeCookie = DefaultCookieMaker
	}
}

// Start begins or resumes a session, returns error if not found, seed
// mismatch or something goes wrong.
//
// It reads session id and seed from cookie, creates one if no relevent
// cookies found, or loads from store if cookies are set.
//
// It's caller's response to decide what to do if session has expired
// or not found.
func (m *Manager) Start(w http.ResponseWriter, r *http.Request) (sess *Session, err error) {
	m.init()
	sid, err := r.Cookie(m.Key)

	if err != nil || sid.Value == "" {
		return newSession(m, w)
	}

	seed, err := r.Cookie(m.Key)

	if err != nil || seed.Value == "" {
		return newSession(m, w)
	}

	return loadSession(sid.Value, seed.Value, m, w)
}

// Session represents session for a specific client
type Session struct {
	id      string
	seed    string
	data    string
	expired bool
	saved   bool
	m       *Manager
}

func newSession(m *Manager, w http.ResponseWriter) (*Session, error) {
	seed := generateSeed()
	id, err := m.Store.Allocate(seed)
	if err != nil {
		return nil, err
	}

	s := &Session{
		id:   id,
		seed: seed,
		m:    m,
	}

	// session id
	c := m.MakeCookie(s.m.Key, s.id, s.m.TTL)
	http.SetCookie(w, c)
	// session checksum
	c = m.MakeCookie(s.m.ChecksumKey, seed, s.m.TTL)
	http.SetCookie(w, c)

	return s, nil
}

func loadSession(id, seed string, m *Manager, w http.ResponseWriter) (*Session, error) {
	expect, data, err := m.Store.Get(id)
	if err != nil {
		return nil, err
	}

	if expect != seed {
		return nil, errors.New("rtoolkit/session: seed mismatch for " + id)
	}

	s := &Session{
		id:   id,
		seed: seed,
		data: data,
		m:    m,
	}
	c := m.MakeCookie(s.m.Key, s.id, s.m.TTL)
	http.SetCookie(w, c)
	return s, nil
}

// ID returns session id
func (s *Session) ID() string {
	return s.id
}

// Destroy sets this session as expired and expires cookie
//
// Calling Destroy() after Save() is a no-op.
// Since it sets cookie, yu SHOULD call it before w.Write().
func (s *Session) Destroy(w http.ResponseWriter) {
	if s.saved {
		return
	}

	s.expired = true
	s.m.Store.Release(s.id)
	c := s.m.MakeCookie(s.m.Key, "", s.m.TTL)
	http.SetCookie(w, c)
	s.saved = true
}

// Save saves data, updates cookie expire time, and delete cookie or session if expired
//
// Calling Save() after Destroy() is a no-op.
// Since it sets cookie, you SHOULD call it before w.Write().
func (s *Session) Save(w http.ResponseWriter) error {
	if s.saved {
		return nil
	}

	err := s.m.Store.Set(s.id, s.seed, s.data)
	if err == nil {
		c := s.m.MakeCookie(s.m.Key, s.id, s.m.TTL)
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
	err := s.m.Store.Set(s.id, s.seed, data)
	if err == nil {
		s.saved = false
		s.expired = false
		s.data = data
	}

	return err
}
