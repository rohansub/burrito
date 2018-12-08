package main

import (
	"fmt"
	"os"

	"github.com/rcsubra2/burrito/src/server"
	"github.com/rcsubra2/burrito/src/parser"
	)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main <burrito_file>")
		return
	}
	routes, err := parser.ParseBurritoFile(os.Args[1])
	if err != nil {
		fmt.Println("Your Burrito has failed to compile")
		fmt.Println("ERROR", err)
		return
	}

	serv, err := server.NewBurritoServer(&routes, nil)
	if err != nil {
		fmt.Println("The zester has caught a route conflict!, server cannot run")
		fmt.Println(err)
		return
	}
	serv.Run()
}
