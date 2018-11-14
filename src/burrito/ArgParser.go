package burrito

import (
	"errors"
	"fmt"
	re "regexp"
)

const (
	// Protocol constants
	GET    string = "GET"
	PUT    string = "PUT"
	POST   string = "POST"
	DELETE string = "DELETE"

	// REGEX constants
	PATHRE string = `\s*'((/([.a-zA-Z0-9-]*)|(<(int|str):[a-z]*>))+)'\s*`
	TYPERE string = `\s*'((GET)|(POST)|(PUT)|(DELETE))'\s*`
)

// Arg - a path and url representing an argument
type Arg struct {
	reqType string
	path    string
}

/**
*
* Parse an Arg object given an argStr, return error when there is a syntax error
 */
func createArg(argStr string) (Arg, error) {
	// Check to see if both path and request are given
	expTwoParams := re.MustCompile(`\(` + PATHRE + `,` + TYPERE + `\)`)
	if expTwoParams.MatchString(argStr) {
		matches := expTwoParams.FindStringSubmatch(argStr)
		return Arg{
			path:    matches[1], // first match is the path
			reqType: matches[6], // sixth match ist the method type
		}, nil
	}
	// Check to see if only path is specified
	expOneParam := re.MustCompile(`\(` + PATHRE + `\)`)
	if expOneParam.MatchString(argStr) {
		matches := expOneParam.FindStringSubmatch(argStr)
		return Arg{
			reqType: GET,
			path:    matches[1],
		}, nil
	}
	// Check to see if no parameters are specified
	expNoParams := re.MustCompile(`\(\)`)
	if expNoParams.MatchString(argStr) {
		return Arg{
			reqType: GET,
			path:    "/",
		}, nil
	}

	// Return Error if compilation failed
	m := fmt.Sprintf("Failed! Argument - %s - had an error", argStr)
	return Arg{}, errors.New(m)
}
