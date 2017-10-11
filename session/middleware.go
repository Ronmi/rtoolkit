package session

import (
	"context"
	"net/http"

	"github.com/Ronmi/rtoolkit/middleware"
)

// ContextKey represnets a key used in context
type ContextKey string

const SessionObjectKey = ContextKey("session")

// NewMiddleware creates a net/http based middleware.
//
// It allocates session instance before executing handler.
func NewMiddleware(m *Manager, h http.Handler) *middleware.Middleware {
	return &middleware.Middleware{
		Next: h,
		Handler: func(w http.ResponseWriter, r *http.Request) (error, *http.Request) {
			req := r
			sess, err := m.Start(w, r)
			if err == nil {
				req = r.WithContext(context.WithValue(r.Context(), SessionObjectKey, sess))
				w.Header().Set("Trailer", "Set-Cookie")
			}

			return err, req
		},
		ErrHandler: func(err error, w http.ResponseWriter, r *http.Request) {
			http.Error(w, err.Error(), 500)
		},
	}
}

// FromMiddleware grabs session object passed by middleware
func FromMiddleware(c context.Context) (sess *Session, found bool) {
	sess, found = c.Value(SessionObjectKey).(*Session)
	return
}
