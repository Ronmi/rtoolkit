// Package middleware provides few tools to create simple middleware
package middleware

import "net/http"

// Middleware defines a simple middleware
//
//
type Middleware struct {
	// Handler is your custom function holding the logic
	Handler func(http.ResponseWriter, *http.Request) (error, *http.Request)

	// Next is executed if Handler returns nil, normally you should
	// set it to a http.ServerMux or another middleware
	Next http.Handler

	// ErrHandler is executed if Handler returns an error. If it is
	// nil, Next is executed even if Handler returns error.
	ErrHandler func(error, http.ResponseWriter, *http.Request)
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, req := m.Handler(w, r)
	if req == nil {
		req = r
	}
	if err != nil && m.ErrHandler != nil {
		m.ErrHandler(err, w, req)
		return
	}

	m.Next.ServeHTTP(w, req)
}

// WrapErrHandler sets ErrHandler to a http.HandlerFunc
//
// It is mostly to use with predefined handlers like
//
//     m.WrapErrHandler(http.NotFound)
func (m *Middleware) WrapErrHandler(h func(http.ResponseWriter, *http.Request)) *Middleware {
	m.ErrHandler = func(e error, w http.ResponseWriter, r *http.Request) {
		h(w, r)
	}

	return m
}

// Clone creates an identical Middleware, and replace the Next handler
func (m *Middleware) Clone(next http.Handler) *Middleware {
	return &Middleware{
		Next:       next,
		Handler:    m.Handler,
		ErrHandler: m.ErrHandler,
	}
}

// CloneFunc is identical to Clone, but accept http.HandlerFunc
func (m *Middleware) CloneFunc(next func(http.ResponseWriter, *http.Request)) *Middleware {
	return &Middleware{
		Next:       http.HandlerFunc(next),
		Handler:    m.Handler,
		ErrHandler: m.ErrHandler,
	}
}
