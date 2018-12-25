package handler

import (
	"github.com/rohansub/burrito/src/environment"
	"strconv"
	"strings"
)

// PathSegment - represents one "Segment" of a path (ie. the path
// "/zesty/breakfast/burrito" would be split into three segments "zesty", "breakfast"
// and "burrito"
type PathSegment struct {
	mustMatch bool
	segStr    string
	typeMatch string
	varName   string
}




// SegMatch - determine if a string matches the given segment
//            returns true, and an EnvEntry if there is a variable in the string
func (ps * PathSegment) SegMatchAndExtractVars(str string) (bool, *environment.EnvEntry) {
	if ps.mustMatch {

		return (ps.segStr == str), nil;

	} else if (ps.typeMatch == "int") {

		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return false, nil;
		}
		return true, environment.CreateIntEntry(ps.varName, i);

	} else if (ps.typeMatch == "float") {

		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return false, nil;
		}
		return true, environment.CreateFloatEntry(ps.varName, f);

	} else { // Otherwise, it must be a string

		return true, environment.CreateStringEntry(ps.varName, str);
	}
}

// NewPathSegment - creates a PathSegment from a string representation
func NewPathSegment(str string) *PathSegment {
	if len(str) == 0 {
		return &PathSegment{
			mustMatch: true,
			segStr: "",
			typeMatch: "",
		}
	} else if str[0] != ':' {
		return &PathSegment {
			mustMatch: true,
			segStr: str,
			typeMatch: "",
		}
	} else {
		ind := strings.Index(str[1:], ":")

		if ind != -1 {
			return &PathSegment {
				mustMatch: false,
				varName: str[1:ind+1],
				typeMatch: str[ind+2:],
			}
		} else {
			return &PathSegment{
				mustMatch: false,
				varName: str[1:],
			}
		}
	}

}

// PathObject - represents one "Segment" of a path (ie. the path
// "/zesty/breakfast/burrito" would be represented as one path object
type PathObject struct {
	parts []*PathSegment
	str string
}


func NewPathObject(url string) *PathObject{
	// Ignores the string before the first "/"
	pathStrLst := strings.Split(url, "/")
	pathObjLst := make([]*PathSegment, len(pathStrLst)-1)

	for i, p := range pathStrLst[1:] {
		pathObjLst[i] = NewPathSegment(p)
	}

	return &PathObject{
		parts: pathObjLst,
		str: url,
	}
}