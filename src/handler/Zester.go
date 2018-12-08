package handler

import (
	"errors"
	"fmt"
)

// Zester - object used to check whether any paths conflict with one another.
type Zester struct {
	paths []*PathObject
}

// Create a Zester object to check for path conflicts
func NewZester() *Zester {
	return &Zester{
		paths: make([]*PathObject, 0),
	}
}

// Add - add a PathSegment to the Zester object
func (z * Zester) Add(path *PathObject) {
	z.paths = append(z.paths, path)
}

// CheckPaths - See if any paths conflict. If so, return an error
// specifying which paths conflict, otherwise return nil
func (z * Zester) CheckPaths() error {
	for i, path1 := range z.paths {
		for _, path2 := range z.paths[i+1:] {
			if isConflicting(path1, path2) {
				m := fmt.Sprintf("Two paths conflict: %s and %s", path1.str, path2.str)
				return errors.New(m)
			}
		}
	}
	return nil
}

func isConflicting(one *PathObject, two *PathObject) bool {
	if len(one.parts) != len(two.parts) {
		return false
	}
	for i, segOne := range one.parts {
		segTwo := two.parts[i]
		// are the two segments not matching?
		if !(getSegmentMatch(segOne, segTwo) ||
			getSegmentMatch(segTwo, segOne) ||
			!segOne.mustMatch && !segTwo.mustMatch) {
			return false
		}
	}
	return true
}

// check true if two path segments could produce conflicting string url matches
// needs to be called twice, inverting the two variables
func getSegmentMatch(one *PathSegment, two *PathSegment) bool{
	return one.mustMatch &&
			(two.mustMatch &&
			two.segStr == one.segStr) ||
			!two.mustMatch && two.typeMatch == "str" ||
			!two.mustMatch && two.typeMatch == ""
}


