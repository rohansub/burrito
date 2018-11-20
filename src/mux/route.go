package mux

import (
	"net/http"
	"strings"
)

type route struct {
	pattern []*PathSegment
	handler http.Handler
}



func (r * route) Match(path string) (bool, *Env) {
	parts := strings.Split(path, "/")

	env := CreateEnv()

	// return false i the length of the url doesn't match this route
	if len(parts) != len(r.pattern) {
		return false, env
	}

	for i, st := range parts {
		p := r.pattern[i]
		match, _ := p.SegMatchAndExtractVars(st)
		if match {
			// append to env
		} else {
			return false, nil
		}

	}
	return true, env


}