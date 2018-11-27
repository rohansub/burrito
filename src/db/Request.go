package db

import (
	"github.com/rcsubra2/burrito/src/environment"
	re "regexp"
	"strings"
)

type Req struct {
	Method string
	GetReq GetReq
}

type Param struct {
	IsString 	bool
	Val 	string
}

func (p* Param) GetValue(envs []*environment.Env) (string, bool){
	if p.IsString {
		return p.Val, true
	}
	for _, e := range envs {
		entry := e.GetValue(p.Val)
		if entry != nil {
			v, ok := entry.(string)
			if ok {
				return v, true
			}
		}
	}
	return "", false
}


type GetReq struct {
	ArgNames []Param
}

// CreateDBGetReq - creates database get request given a list of string arguments
func CreateDBGetReq(argStrs []string) *GetReq {
	args := make([]Param, 0)
	str := re.MustCompile(`^'(\w*)'$`)

	for _, s := range argStrs {
		stripped := strings.Trim(s, " ")
		var arg Param
		// check if it is a string
		if str.MatchString(stripped) {
			matches := str.FindStringSubmatch(stripped)
			arg = Param{
				IsString: true,
				Val: matches[1],
			}
		} else { // otherwise it is a variable
			arg = Param{
				IsString: false,
				Val: stripped,
			}
		}
		args = append(args, arg)
	}
	return &GetReq{
		ArgNames: args,
	}
}

