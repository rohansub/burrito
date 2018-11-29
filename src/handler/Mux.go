package handler

import (
	"github.com/rcsubra2/burrito/src/environment"
	"net/http"
	"strings"
)

// Router - list of routes that are handled by the burrito server
type Router struct {
	routes []*route
}

// BurritoHandlerFunc - type is an adapter to allow the use of
// ordinary functions as HTTP handlers, similar to http.HandlerFunc, but
// with a burrito environment variable included
type BurritoHandlerFunc func(http.ResponseWriter, *http.Request, *environment.Env);

// ServeHTTP calls f(w, r, e).
func (f BurritoHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, e *environment.Env) {
	f(w, r, e)
}

// BurritoHandler - responds to an HTTP request similar to http.Handler, but with a burrito Env included as
// a parameter
type BurritoHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, *environment.Env)
}


func NewRouter() (*Router) {
	return &Router {
		routes: make([]*route, 0),
	}
}

func extractList(lst []string) []*PathSegment{
	// Ignores the string before the first "/"
	lstPaths := make([]*PathSegment, len(lst)-1)
	for i, p := range lst[1:] {
		lstPaths[i] = NewPathSegment(p)
	}
	return lstPaths
}



// Source: https://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc

func (h *Router) Handler(pattern string, handler BurritoHandler) {
	lst := strings.Split(pattern, "/")
	h.routes = append(h.routes, &route{extractList(lst), handler})
}

func (h *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request, *environment.Env)) {
	lst := strings.Split(pattern, "/")
	h.routes = append(h.routes, &route{extractList(lst), BurritoHandlerFunc(handler)})
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		match, env := route.Match(r.URL.Path)
		if match {
			route.handler.ServeHTTP(w, r, env)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}
