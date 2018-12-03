package parser

import (
	"github.com/rcsubra2/burrito/src/db"
	re "regexp"
	"strings"
)



func createRespForDB(respStr string) db.Req {
	isGet := re.MustCompile(`^DB.GET\(((?:(?:\w+|(?:'\w*'))\s*,\s*)*)\)$`)
	if isGet.MatchString(respStr) {
		matches := isGet.FindStringSubmatch(respStr)
		argStrs := strings.Split(matches[1], ",")
		gr := db.CreateDBGetReq(argStrs[0:len(argStrs)-1])
		return gr
	}

	pair := `\((?:(?:\w+|(?:'\w*'))\s*),\s*(?:(?:\w+|(?:'\w*'))\s*)\)`
	pairList := `(?:` + pair + `,\s*)*`
	isSet := re.MustCompile(`^DB.SET\(`+ pairList + `\)$`)

	if isSet.MatchString(respStr) {
		pairRE := re.MustCompile(pair)
		matches := pairRE.FindAllString(respStr, -1)
		gr := db.CreateDBSetReq(matches)
		return gr
	}

	isDel := re.MustCompile(`^DB.DEL\(((?:(?:\w+|(?:'\w*'))\s*,\s*)*)\)$`)
	if isDel.MatchString(respStr) {
		matches := isDel.FindStringSubmatch(respStr)
		argStrs := strings.Split(matches[1], ",")
		gr := db.CreateDBDelReq(argStrs[0:len(argStrs)-1])
		return gr
	}

	return nil
}


