package jsonapi

// Middleware is a wrapper for handler
type Middleware func(Handler) Handler

// Registerer represents a chain of middleware
//
//     With(
//         myMiddleware
//     ).With(
//         apitool.LogIn(apitool.JSONFormat(
//             log.New(os.Stdout, "myapp", log.LstdFlags),
//         )),
//     ).RegisterAll(mux, "/api", myHandler)
//
// Request processing flow will be:
//
//     1. mux.ServeHTTP
//     2. myMiddleWare
//     3. Logging middleware
//     4. myHandler
type Registerer interface {
	Register(mux HTTPMux, apis []API)
	RegisterAll(mux HTTPMux, prefix string, handlers interface{})
	With(m Middleware) Registerer
}

// With creates a new Registerer
func With(m Middleware) Registerer {
	return &registerer{
		m: m,
	}
}

type registerer struct {
	m      Middleware
	parent Registerer
}

// Register is identical to jsonapi.Register(), but wraps api in middleware chain first
func (r *registerer) Register(mux HTTPMux, apis []API) {
	reg := Register
	if r.parent != nil {
		reg = r.parent.Register
	}

	for x, a := range apis {
		apis[x].Handler = r.m(a.Handler)
	}

	reg(mux, apis)
}

// RegisterAll is identical to jsonapi.RegisterAll(), but wraps api in middleware chain first
func (r *registerer) RegisterAll(mux HTTPMux, prefix string, handlers interface{}) {
	r.Register(mux, findMatchedMethods(prefix, handlers))
}

// With creaates a new Registerer and chains after current Registerer
func (r *registerer) With(m Middleware) Registerer {
	return &registerer{
		m:      m,
		parent: r,
	}
}
