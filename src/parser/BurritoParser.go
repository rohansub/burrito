package parser

import (
	"errors"
	. "fmt"
	"io/ioutil"
	"log"
	str "strings"
)

// ParsedRoutes contains a map of Arg structs to the Resp object associated with it
type ParsedRoutes struct {
	routes map[Arg][]Resp
}

// AddRules - Add a rule to an existing ParsedRoutes Object
func (rts *ParsedRoutes) AddRules(ar Arg, bodies []Resp) {
	_, ok := rts.routes[ar]
	if ok {
		log.Fatalf("Multiple Arguments with path %s and route %s were found",
			ar.path, ar.reqType)
	}
	rts.routes[ar] = bodies
}

// ParseBurritoFile - parse a burrito file into a ParsedRoutes object
func ParseBurritoFile(filepath string) (ParsedRoutes, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return ParsedRoutes{}, err
	}
	source := Sprintf("%s", content)
	routes := str.Split(source, ";")

	pbData := ParsedRoutes{
		routes: make(map[Arg][]Resp),
	}

	for _, r := range routes {
		line := str.Trim(r, "\n ")
		if len(line) == 0 {
			continue
		}
		// handle commented out lines
		if line[0] == '#' {
			lns := str.Split(line, "\n")
			cont := true
			for _, l := range lns {
				if l[0] != '#' {
					cont = false
					line = l
				}
			}
			if cont {
				continue
			}
		}
		parts := str.Split(line, "=>")
		if len(parts) <= 1 {
			m := "Failed to compile this line: " + line
			return ParsedRoutes{}, errors.New(m)
		}
		arg, err := createArg(parts[0])
		if err != nil {
			m := "Failed to compile this line: " + line
			return ParsedRoutes{}, errors.New(m)
		}

		resp := make([]Resp, len(parts)-1)
		for i, partStr := range parts[1:] {
			resp[i], err = createResp(str.Trim(partStr, " "))
			if err != nil {
				m := "Failed to compile this response part: " + partStr
				return ParsedRoutes{}, errors.New(m)
			}
		}
		pbData.AddRules(arg, resp)
	}
	return pbData, nil
}
