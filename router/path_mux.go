package router

import (
	"errors"
	"net/http"
	"strings"
)

// pathNode is an element in mapping tree, dispatching by path
//
// A pathNode can contains mappings to child pathNode and a handler
type pathNode struct {
	child    map[string]*pathNode
	catchAll *pathNode
	h        http.Handler
}

func createPathNode() *pathNode {
	return &pathNode{child: map[string]*pathNode{}}
}

// Go idiom
func (n *pathNode) match(r *http.Request) (h http.Handler, found bool) {
	arr := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	return n.doMatch(arr)
}

func (n *pathNode) doMatch(arr []string) (h http.Handler, found bool) {
	if len(arr) < 1 {
		if n.h != nil {
			return n.h, true
		}
		return
	}

	cur := arr[0]
	if next, ok := n.child[cur]; ok {
		h, found = next.doMatch(arr[1:])
	}

	if !found && n.catchAll != nil {
		h, found = n.catchAll.doMatch(arr[1:])
	}

	return
}

func (n *pathNode) register(wild, pattern string, h http.Handler) {
	arr := strings.Split(strings.Trim(pattern, "/"), "/")
	n.doRegister(wild, arr, h)
}

func (n *pathNode) doRegister(wild string, arr []string, h http.Handler) {
	if len(arr) < 1 {
		n.h = h
		return
	}

	cur := arr[0]
	if cur == wild {
		if n.catchAll == nil {
			n.catchAll = createPathNode()
		}
		n.catchAll.doRegister(wild, arr[1:], h)
		return
	}

	next, ok := n.child[cur]
	if !ok {
		next = createPathNode()
		n.child[cur] = next
	}

	next.doRegister(wild, arr[1:], h)
}

// PathMux is a http.ServerMux compitable mux implementation, dispatches by path
type PathMux struct {
	mappings   *pathNode
	Wildcard   string
	ErrHandler http.Handler
}

func errHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	return
}

// ByPath creates a new PathMux with default settings.
//
// Wildcard defaults to "*". Error handler returns 404 NOT FOUND for every error.
func ByPath() *PathMux {
	return &PathMux{
		mappings:   createPathNode(),
		Wildcard:   "*",
		ErrHandler: http.HandlerFunc(errHandler),
	}
}

// Handle registers a handler for specified pattern
//
// It panics if pattern is invalid.
//
// This method is not thread-safe.
func (m *PathMux) Handle(pattern string, h http.Handler) {
	if !strings.HasPrefix(pattern, "/") {
		panic(errors.New("mux: pattern must begin with /"))
	}

	m.mappings.register(m.Wildcard, pattern, h)
}

// HandleFunc registers a handler function for specified pattern
//
// It panics if pattern is invalid.
//
// This method is not thread-safe.
func (m *PathMux) HandleFunc(pattern string, h func(http.ResponseWriter, *http.Request)) {
	m.Handle(pattern, http.HandlerFunc(h))
}

// ServeHTTP finds correct handler and executes it, or use PathMux.ErrHandler if no match
func (m *PathMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, found := m.mappings.match(r); found {
		h.ServeHTTP(w, r)
		return
	}
	m.ErrHandler.ServeHTTP(w, r)
}
