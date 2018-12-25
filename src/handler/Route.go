package handler

import (
	"github.com/rohansub/burrito/src/environment"
	"strings"
)

type route struct {
	pattern *PathObject
	handler BurritoHandler
}

// Match - return whether a path matches with the route specified
func (r * route) Match(path string) (bool, *environment.Env) {
	// Extract the parts of the url by splitting by "/", ignore
	// the part of the path before the first "/"
	parts := strings.Split(path, "/")[1:]


	// return false i the length of the url doesn't match this route
	if len(parts) != len(r.pattern.parts) {
		return false, nil
	}

	env := environment.CreateEnv()

	for i, st := range parts {
		p := r.pattern.parts[i]
		match, entry := p.SegMatchAndExtractVars(st)
		if !match {
			return false, nil
		}
		if entry != nil {
			env.Add(*entry)
		}
	}
	return true, env


}