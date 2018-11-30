package db

import (
	"github.com/rcsubra2/burrito/src/environment"
	re "regexp"
	"strings"
)

type Req interface {
	Run(client Database, envs []*environment.Env) map[string]string
}

type Param struct {
	IsString 	bool
	Val 	string
}

// GetValue - given a list of environments find the value of the parameter
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

//GetReq - a struct that is an implements Req, and used as structure for get request
type GetReq struct {
	ArgNames []Param
}

// CreateDBGetReq - creates database get request given a list of string arguments
func CreateDBGetReq(argStrs []string) *GetReq {
	args := make([]Param, 0)
	for _, s := range argStrs {
		stripped := strings.Trim(s, " ")
		var arg Param = extractParam(stripped)
		args = append(args, arg)
	}
	return &GetReq{
		ArgNames: args,
	}
}

// Run - perform the request on given database.
func (r * GetReq) Run(client Database, envs []*environment.Env) map[string]string {
	keys := make([]string, len(r.ArgNames))
	for i, ar := range r.ArgNames {
		val, ok := ar.GetValue(envs)
		if ok {
			keys[i] = val
		}
	}

	return client.Get(keys)
}

type Pair struct {
	Fst Param
	Snd Param
}


type SetReq struct {
	ArgNames []Pair
}

// CreateDBGetReq - creates database get request given a list of string arguments
func CreateDBSetReq(argStrs []string) *SetReq {
	args := make([]Pair, 0)
	for _, s := range argStrs {
		stripped := strings.Trim(s, " ")
		rePair := re.MustCompile(`\(\s*(.*)\s*,\s*(.*)\s*\)`)
		if rePair.MatchString(stripped){
			matches := rePair.FindStringSubmatch(stripped)
			p := Pair {
				Fst: extractParam(matches[1]),
				Snd: extractParam(matches[2]),
			}
			args = append(args, p)
		} else {
			panic("CreateDBSetReq not called correctly!, " +
				"list of strings cannot be parsed into SetReq")
		}

		extractParam(stripped)
	}
	return &SetReq{
		ArgNames: args,
	}
}

func (req * SetReq) Run(client Database, envs []*environment.Env) map[string]string {
	return nil
}


func extractParam(strParam string) Param{
	strRE := re.MustCompile(`^'(\w*)'$`)
	if strRE.MatchString(strParam) {
		matches := strRE.FindStringSubmatch(strParam)
		return Param{
			IsString: true,
			Val:      matches[1],
		}
	} else { // otherwise it is a variable
		return Param{
			IsString: false,
			Val:      strParam,
		}
	}
}