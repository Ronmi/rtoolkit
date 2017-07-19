package session

import (
	"context"
	"net/http"
)

// ContextKey represnets a key used in context
type ContextKey string

const SessionObjectKey = ContextKey("session")

// Middleware add session to context before running your handler/router
type Middleware interface {
	http.Handler
}

// NewMiddleware creates a net/http based middleware.
//
// It allocates session instance before executing handler.
func NewMiddleware(m *Manager, h http.Handler) Middleware {
	return middleware{
		manager: m,
		h:       h,
	}
}

type middleware struct {
	manager *Manager
	h       http.Handler
}

func (m middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := r
	sess, err := m.manager.Start(w, r)
	if err == nil {
		req = r.WithContext(context.WithValue(r.Context(), SessionObjectKey, sess))
		w.Header().Set("Trailer", "Set-Cookie")
	}

	m.h.ServeHTTP(w, req)
}

// FromMiddleware grabs session object passed by middleware
func FromMiddleware(c context.Context) (sess *Session, found bool) {
	sess, found = c.Value(SessionObjectKey).(*Session)
	return
}
