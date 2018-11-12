package burrito

import (
	"errors"
	re "regexp"
)

type Resp struct {
	respType string // file/json/string
	body     string
}

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
