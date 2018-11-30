package parser

import (
	"github.com/rcsubra2/burrito/src/db"
	re "regexp"
	"strings"
)



func createRespForDB(respStr string) db.Req {
	isGet := re.MustCompile(`^DB.GET\(((?:(?:\w+|'\w*')\s*,\s*)*)\)$`)
	if isGet.MatchString(respStr) {
		matches := isGet.FindStringSubmatch(respStr)
		argStrs := strings.Split(matches[1], ",")
		gr := db.CreateDBGetReq(argStrs[0:len(argStrs)-1])
		return gr
	}

	return nil
}


