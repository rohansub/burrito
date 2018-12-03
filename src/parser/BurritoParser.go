package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	str "strings"
)

// ParsedRoutes contains a map of Arg structs to the Resp object associated with it
type ParsedRoutes struct {
	Routes map[string]map[string][]Resp // Key - URL, Value - Map (Key - Method type, Value - Response)
}

// AddRules - Add a rule to an existing ParsedRoutes Object
func (rts *ParsedRoutes) AddRules(ar Arg, bodies []Resp) error {
	_, pathExists := rts.Routes[ar.path]
	if !pathExists {
		rts.Routes[ar.path] = make(map[string][]Resp)
	}
	_, methodExists := rts.Routes[ar.path][ar.reqType]
	if methodExists {

		m := fmt.Sprintf("Multiple Arguments with path %s and route %s were found",
			ar.path, ar.reqType)
		return errors.New(m)
	}
	rts.Routes[ar.path][ar.reqType] = bodies
	return nil

}

// ParseBurritoFile - parse a server file into a ParsedRoutes object
func ParseBurritoFile(filepath string) (ParsedRoutes, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return ParsedRoutes{}, err
	}
	source := fmt.Sprintf("%s", content)
	routes := str.Split(source, ";")

	pbData := ParsedRoutes{
		Routes: make(map[string]map[string][]Resp),
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
		arg, err := createArg(str.Trim(parts[0], " \n\t"))
		if err != nil {
			m := "Failed to compile this line: " + line
			return ParsedRoutes{}, errors.New(m)
		}

		resp := make([]Resp, len(parts)-1)
		for i, partStr := range parts[1:] {
			resp[i], err = createResp(str.Trim(partStr, " \n\t"))
			if err != nil {
				m := "Failed to compile this response part: " + partStr
				return ParsedRoutes{}, errors.New(m)
			}
		}
		err = pbData.AddRules(arg, resp)

		if err != nil {
			return ParsedRoutes{}, err
		}
	}
	return pbData, nil
}
