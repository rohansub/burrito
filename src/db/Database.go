package db

import (
	"github.com/rohansub/burrito/src/environment"
)

// DatabaseFunction - function performed on a database given an environment
type DatabaseFunction func(group environment.EnvironmentGroup) (map[string]interface{}, error)

// DatabaseAction - Name of database, function and arguments. Specify the
// type database call to be made
type DatabaseAction struct {
	Name string
	Fname string
	Args string
}

// Database - Interface implemented by any database integration that is implemented
type Database interface {
	IsCorrectSyntax(fname string, args string) bool
	Run(fname string, args string, group environment.EnvironmentGroup) (map[string]interface{}, error)
}

// DatabaseConstructor - A function that Creates a Database connection. Any constructor of
// a Database object should have a constructor of this type
type DatabaseConstructor func(isMock bool,  url string, pswd string) Database

