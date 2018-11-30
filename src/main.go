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

	serv := server.NewBurritoServer(&routes, nil)

	serv.Run()
}
