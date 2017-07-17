package router

import (
	"errors"
	"net/http"
	"strings"
)

// node is an element in mapping tree
//
// A node can contains mappings to child node and a handler
type node struct {
	child    map[string]*node
	catchAll *node
	h        http.Handler
}

func createNode() *node {
	return &node{child: map[string]*node{}}
}

// Go idiom
func (n *node) match(arr []string) (h http.Handler, found bool) {
	if len(arr) < 1 {
		if n.h != nil {
			return n.h, true
		}
		return
	}

	cur := arr[0]
	if next, ok := n.child[cur]; ok {
		h, found = next.match(arr[1:])
	}

	if !found && n.catchAll != nil {
		h, found = n.catchAll.match(arr[1:])
	}

	return
}

func (n *node) register(wild string, arr []string, h http.Handler) {
	if len(arr) < 1 {
		n.h = h
		return
	}

	cur := arr[0]
	if cur == wild {
		if n.catchAll == nil {
			n.catchAll = createNode()
		}
		n.catchAll.register(wild, arr[1:], h)
		return
	}

	next, ok := n.child[cur]
	if !ok {
		next = createNode()
		n.child[cur] = next
	}

	next.register(wild, arr[1:], h)
}

// Mux is a http.ServerMux compitable mux implementation
type Mux struct {
	mappings   *node
	Wildcard   string
	ErrHandler http.Handler
}

func errHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	return
}

// New creates a new Mux with default wildcard and error handler.
//
// Wildcard defaults to "*". Error handler returns 404 NOT FOUND for every error.
func New() *Mux {
	return &Mux{
		mappings:   createNode(),
		Wildcard:   "*",
		ErrHandler: http.HandlerFunc(errHandler),
	}
}

// Handle registers a handler for specified pattern
//
// It panics if pattern is invalid.
//
// This method is not thread-safe.
func (m *Mux) Handle(pattern string, h http.Handler) {
	if !strings.HasPrefix(pattern, "/") {
		panic(errors.New("mux: pattern must begin with /"))
	}

	arr := strings.Split(strings.Trim(pattern, "/"), "/")
	m.mappings.register(m.Wildcard, arr, h)
}

// Handle registers a handler for specified pattern
//
// It panics if pattern is invalid.
//
// This method is not thread-safe.
func (m *Mux) HandleFunc(pattern string, h func(http.ResponseWriter, *http.Request)) {
	m.Handle(pattern, http.HandlerFunc(h))
}

// Dispatch finds correct handler, or return Mux.ErrHandler if error
func (m *Mux) Dispatch(uri string) http.Handler {
	arr := strings.Split(strings.Trim(uri, "/"), "/")
	if h, ok := m.mappings.match(arr); ok {
		return h
	}

	return m.ErrHandler
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Dispatch(r.URL.Path).ServeHTTP(w, r)
}
