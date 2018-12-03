package db

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rcsubra2/burrito/src/environment"
	"github.com/rcsubra2/burrito/src/utils"
	re "regexp"
	"strings"
)

type Req interface {
	Run(client Database, envs []*environment.Env) (map[string]string, error)
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
func (r * GetReq) Run(client Database, envs []*environment.Env) (map[string]string, error) {
	keys := make([]string, len(r.ArgNames))
	for i, ar := range r.ArgNames {
		val, ok := ar.GetValue(envs)
		if ok {
			keys[i] = val
		}
	}
	return client.Get(keys), nil
}


type SetReq struct {
	ArgNames []utils.Pair
}

// CreateDBSetReq - creates database get request given a list of string arguments
func CreateDBSetReq(argStrs []string) *SetReq {
	args := make([]utils.Pair, 0)
	for _, s := range argStrs {
		stripped := strings.Trim(s, " ")
		rePair := re.MustCompile(`\(\s*(.*)\s*,\s*(.*)\s*\)`)
		if rePair.MatchString(stripped){
			matches := rePair.FindStringSubmatch(stripped)
			p := utils.Pair {
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

func (req * SetReq) Run(client Database, envs []*environment.Env) (map[string]string, error) {
	pairs := make([]utils.Pair, len(req.ArgNames))
	for i, ar := range req.ArgNames {
		kParam := ar.Fst.(Param)
		vParam := ar.Snd.(Param)

		k, okKey := kParam.GetValue(envs)
		v, okVal := vParam.GetValue(envs)

		if okKey && okVal {
			pairs[i] = utils.Pair{
				Fst: k,
				Snd: v,
			}
		} else {
			m := fmt.Sprintf("No such variable %s in environment", ar.Snd)
			return map[string]string{}, errors.New(m)
		}
	}
	v := client.Set(pairs)
	if v {
		return map[string]string{}, nil
	} else {
		return map[string]string{}, errors.New("DB call failed")
	}

}


func extractParam(strParam string) Param {
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