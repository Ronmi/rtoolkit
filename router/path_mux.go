package router

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// ContextKey represents a key used in context
type ContextKey string

const PathVarKey = ContextKey("pathData")

// GetPathVariable extracts variables from context
func GetPathVariable(c context.Context) (data []string, ok bool) {
	v := c.Value(PathVarKey)
	if v == nil {
		return
	}

	data, ok = v.([]string)
	return
}

// FillPathVariable is simple halper to load path variables into custom variable
//
// It returns how many variables are filled. Zero if error.
//
// Say we have a pattern /user/*/*
//
//    var uidStr, action string
//    if FillPathVariable(request.Context(), &uidStr, &action) < 0 {
//            log.Fatalf("unknown user or action")
//    }
func FillPathVariable(c context.Context, vars ...*string) int {
	data, ok := GetPathVariable(c)
	if !ok {
		return 0
	}

	min := len(data)
	if l := len(vars); l < min {
		min = l
	}

	for i := 0; i < min; i++ {
		v := vars[i]
		*v = data[i]
	}

	return min
}

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
func (n *pathNode) match(r *http.Request) (h http.Handler, data []string, found bool) {
	arr := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	return n.doMatch(arr, make([]string, 0, len(arr)))
}

func (n *pathNode) doMatch(arr []string, oldData []string) (h http.Handler, data []string, found bool) {
	data = oldData
	if len(arr) < 1 {
		if n.h != nil {
			return n.h, oldData, true
		}
		return
	}

	cur := arr[0]
	if cur != "" {
		if next, ok := n.child[cur]; ok {
			if h, data, found = next.doMatch(arr[1:], oldData); found {
				return
			}
		}

		if n.catchAll != nil {
			tmpData := oldData
			l := len(tmpData)
			tmpData = append(tmpData, "")
			if h, data, found = n.catchAll.doMatch(arr[1:], tmpData); found {
				data[l] = cur
				return
			}
			data = oldData
		}
	}

	if next, ok := n.child[""]; ok {
		return next.h, data, true
	}

	return
}

func (n *pathNode) register(wild, pattern string, h http.Handler) {
	arr := strings.Split(strings.TrimLeft(pattern, "/"), "/")
	n.doRegister(wild, arr, h)
}

func (n *pathNode) doRegister(wild string, arr []string, h http.Handler) {
	if len(arr) < 1 {
		n.h = h
		return
	}

	cur := arr[0]
	if strings.Index(cur, wild) != -1 && cur != wild {
		panic("router: wildcard cannot use with others")
	}
	if len(arr) != 1 && cur == "" {
		panic("router: pattern cannot have empty string")
	}
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
	if h, data, found := m.mappings.match(r); found {
		req := r.WithContext(context.WithValue(r.Context(), PathVarKey, data))
		h.ServeHTTP(w, req)
		return
	}
	m.ErrHandler.ServeHTTP(w, r)
}
