package jsonapi

import "net/http"

// Middleware is a wrapper for handler
type Middleware func(Handler) Handler

// Registerer is set of tools to enable use of middleware
//
//     With(
//         mySessionMiddleware()
//     ).With(
//         apilog.Use(apilog.JSON(log.New(os.Stdout, "myapp", log.LstdFlags)))
//     ).RegisterAll(mux, "/api", myHandler)
//
// Request processing flow will be:
//
//     1. mux.ServeHTTP
//     2. mySessionMiddleWare
//     3. JSON logging middleware
//     4. myHandler
type Registerer interface {
	Register(apis []API, mux *http.ServeMux)
	RegisterAll(mux *http.ServeMux, prefix string, handlers interface{})
	With(m Middleware) Registerer
}

func With(m Middleware) Registerer {
	return &registerer{
		m: m,
	}
}

type registerer struct {
	m      Middleware
	parent Registerer
}

func (r *registerer) Register(apis []API, mux *http.ServeMux) {
	for x, a := range apis {
		apis[x].Handler = r.m(a.Handler)
	}

	if r.parent != nil {
		r.parent.Register(apis, mux)
		return
	}

	Register(apis, mux)
}

func (r *registerer) RegisterAll(mux *http.ServeMux, prefix string, handlers interface{}) {
	r.Register(findMatchedMethods(prefix, handlers), mux)
}

func (r *registerer) With(m Middleware) Registerer {
	return &registerer{
		m:      m,
		parent: r,
	}
}
