package mux

import (
	"strconv"
	"strings"
)

type PathSegment struct{
	mustMatch bool
	segStr    string
	typeMatch string
	varName   string
}




// SegMatch - determine if a string matches the given segment
//            returns true, and an EnvEntry if there is a variable in the string
func (ps * PathSegment) SegMatchAndExtractVars(str string) (bool, *EnvEntry) {
	if ps.mustMatch {

		return (ps.segStr == str), nil;

	} else if (ps.typeMatch == "int") {

		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return false, nil;
		}
		return true, CreateIntEntry(ps.varName, i);

	} else if (ps.typeMatch == "float") {

		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return false, nil;
		}
		return true, CreateFloatEntry(ps.varName, f);

	} else { // Otherwise, it must be a string

		return true, CreateStringEntry(ps.varName, str);
	}
}





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
