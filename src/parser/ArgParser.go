package parser

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
	PATH_NO_VAR string = `(\w*)`
	PATH_VAR string = `((:(\w*))(:((int)|(str)|(flt)))?)`
	PATHRE string = `\s*'((/(` + PATH_NO_VAR + `|`+ PATH_VAR +`))+)'\s*`
	TYPERE string = `\s*'((GET)|(POST)|(PUT)|(DELETE))'\s*`
)
//| (/:[\w](:(int|str|flt))?) )+
// Arg - a path and url representing an argument
type Arg struct {
	reqType string
	path    string
}


func getPath(path string) (string, error){
	expTwoParams := re.MustCompile(PATHRE)
	if expTwoParams.MatchString(path) {
		matches := expTwoParams.FindStringSubmatch(path)
		return matches[1], nil
	}
	return "", errors.New("Match not found")
}

func getMethod(method string) (string, error){
	expTwoParams := re.MustCompile(TYPERE)
	if expTwoParams.MatchString(method) {
		matches := expTwoParams.FindStringSubmatch(method)
		return matches[1], nil
	}
	return "", errors.New("Match not found")
}

/**
*
* Parse an Arg object given an argStr, return error when there is a syntax error
 */
func createArg(argStr string) (Arg, error) {
	// Check to see if both path and request are given
	expTwoParams := re.MustCompile(`\(([^,]*),([^,]*)\)`)
	if expTwoParams.MatchString(argStr) {
		matches := expTwoParams.FindStringSubmatch(argStr)
		pth, err1 := getPath(matches[1])
		meth, err2 := getMethod(matches[2])

		if err1 == nil && err2 == nil {
			return Arg{
				path:    pth, // first match is the path
				reqType: meth, // sixth match ist the method type
			}, nil
		}


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
