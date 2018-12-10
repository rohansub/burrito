package parser

import (
	"github.com/rcsubra2/burrito/src/db"
	re "regexp"
)



func createDBCall(respStr string, databases map[string]db.Database) *db.DatabaseAction{
	isDB := re.MustCompile(`^(\w+).(\w+)\((.*)\)$`)
	if !isDB.MatchString(respStr) {
		return nil
	}
	matches := isDB.FindStringSubmatch(respStr)

	name := matches[1]
	fname := matches[2]
	args := matches[3]

	data, ok := databases[name]

	if ok && data.IsCorrectSyntax(fname, args){
		return &db.DatabaseAction{
			Name: name,
			Fname: fname,
			Args: args,
		}
	}

	return nil
}


