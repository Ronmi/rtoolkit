package session

import (
	"context"
	"net/http"
)

// Middleware is a net/http based middleware to handle session, both fields are required
type Middleware struct {
	Manager *Manager
	Handler http.HandlerFunc
}

func (m *Middleware) Handle(w http.ResponseWriter, r *http.Request) {
	req := r
	sess, err := m.Manager.Start(r)
	if err == nil {
		req = r.WithContext(context.WithValue(r.Context(), "session", sess))
		w.Header().Set("Trailer", "Set-Cookie")
	}

	m.Handler(w, req)
}
