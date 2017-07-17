package router

import "net/http"

// Node represents an element in mapping tree
type Node interface {
	Register(wild, pattern string, h http.Handler)
	Serve(fallback http.Handler, w http.ResponseWriter, r *http.Request)
}
