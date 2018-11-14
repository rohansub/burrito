package burrito

import (
	"errors"
	re "regexp"
)

// Resp struct representing Response - contains a body and a respType
type Resp struct {
	respType string // FILE or JSON or STR
	body     string
}

// createResp - parse a Resp struct from the string input,
//              return error if string can't be parsed
func createResp(respStr string) (Resp, error) {
	isString := re.MustCompile(`s'(.*)'`)
	if isString.MatchString(respStr) {
		matches := isString.FindStringSubmatch(respStr)
		return Resp{
			respType: "STR",
			body:     matches[1],
		}, nil
	}

	isFile := re.MustCompile(`'(.*)'`)
	if isFile.MatchString(respStr) {
		matches := isFile.FindStringSubmatch(respStr)
		return Resp{
			respType: "FILE",
			body:     matches[1],
		}, nil
	}

	return Resp{}, errors.New("Response Body - " + respStr + " - not recognized")
}
