package parser

import (
	"encoding/json"
	"errors"
	"github.com/rohansub/burrito/src/db"
	re "regexp"
)

// Resp struct representing Response - contains a Body and a RespType
type Resp struct {
	RespType string // FILE or JSON or STR or DB
	Body     interface{}
	DBReq    db.DatabaseAction
}

// createResp - parse a Resp struct from the string input,
//              return error if string can't be parsed
func createResp(respStr string, databases map[string]db.Database) (Resp, error) {
	isString := re.MustCompile(`^s'(.*)'$`)
	if isString.MatchString(respStr) {
		matches := isString.FindStringSubmatch(respStr)
		return Resp{
			RespType: "STR",
			Body:     matches[1],
		}, nil
	}

	isFile := re.MustCompile(`^('(.*)')$`)
	if isFile.MatchString(respStr) {
		matches := isFile.FindStringSubmatch(respStr)
		return Resp{
			RespType: "FILE",
			Body:     matches[2],
		}, nil
	}

	isDB := createDBCall(respStr, databases)
	if isDB != nil {
		return Resp{
			RespType: "DB",
			DBReq: 	  *isDB,
		}, nil
	}


	var data interface{}
	err := json.Unmarshal([]byte(respStr), &data)
	if err == nil {
		return Resp{
			RespType: "JSON",
			Body: data,
		}, nil
	}


	return Resp{}, errors.New("Response Body - " + respStr + " - not recognized")
}
