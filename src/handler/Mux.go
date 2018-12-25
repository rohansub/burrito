package handler

import (
	"github.com/rohansub/burrito/src/environment"
	"net/http"
)

// Router - list of routes that are handled by the burrito server
type Router struct {
	routes []*route
	zester Zester
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

func (r* Router) CheckRoutes() error {
	return r.zester.CheckPaths()
}




// Source: https://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc

func (h *Router) Handler(pattern string, handler BurritoHandler) {
	pathObj := NewPathObject(pattern)
	h.zester.Add(pathObj)
	h.routes = append(h.routes, &route{pathObj, handler})
}

func (h *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request, *environment.Env)) {
	pathObj := NewPathObject(pattern)
	h.zester.Add(pathObj)
	h.routes = append(h.routes, &route{pathObj, BurritoHandlerFunc(handler)})
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
