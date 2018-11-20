package mux

import (
	"net/http"
	"strings"
)


type Router struct {
	routes []*route
}

func NewRouter() (*Router) {
	return &Router {
		routes: make([]*route, 0),
	}
}

func extractList(lst []string) []*PathSegment{
	lstPaths := make([]*PathSegment, len(lst))
	for i, p := range lst {
		lstPaths[i] = NewPathSegment(p)
	}
	return lstPaths
}



// Source: https://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc

func (h *Router) Handler(pattern string, handler http.Handler) {
	lst := strings.Split(pattern, "/")
	h.routes = append(h.routes, &route{extractList(lst), handler})
}

func (h *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	lst := strings.Split(pattern, "/")
	h.routes = append(h.routes, &route{extractList(lst), http.HandlerFunc(handler)})
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		match, _ := route.Match(r.URL.Path)
		// TODO: ADD to env is necessary
		if match {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}
