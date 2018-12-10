package db

import (
	"github.com/rcsubra2/burrito/src/environment"
)

type DatabaseFunction func(group environment.EnvironmentGroup) (map[string]interface{}, error)

type DatabaseAction struct {
	Name string
	Fname string
	Args string
}


type Database interface {
	IsCorrectSyntax(fname string, args string) bool
	Run(fname string, args string, group environment.EnvironmentGroup) (map[string]interface{}, error)
}

