package session

import (
	"context"
	"net/http"
)

// Middleware is something like http.ServeMux, but allocates a session right before executing your handler.
type Middleware interface {
	http.Handler
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// NewMiddleware creates a net/http based middleware.
//
// It allocates session instance before executing handler.
func NewMiddleware(m *Manager, mux *http.ServeMux) Middleware {
	if mux == nil {
		mux = http.NewServeMux()
	}
	return middleware{
		manager:  m,
		ServeMux: mux,
	}
}

type middleware struct {
	manager *Manager
	*http.ServeMux
}

func (m middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, pat := m.Handler(r)
	req := r
	if pat != "" {
		sess, err := m.manager.Start(w, r)
		if err == nil {
			req = r.WithContext(context.WithValue(r.Context(), "session", sess))
			w.Header().Set("Trailer", "Set-Cookie")
		}
	}

	handler.ServeHTTP(w, req)
}
