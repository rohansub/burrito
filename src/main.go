package main

import (
	"fmt"
	"github.com/rcsubra2/burrito/src/config"
	"os"

	"github.com/rcsubra2/burrito/src/server"
	"github.com/rcsubra2/burrito/src/parser"
	)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main <burrito_file>")
		return
	}

	cfg, err := config.NewConfigFromFile("burritoconfig.json")
	if err != nil{
		fmt.Println("Error in Config File!")
		return
	}

	clients := cfg.CreateDatabaseClients()

	routes, err := parser.ParseBurritoFile(os.Args[1], clients)
	if err != nil {
		fmt.Println("Your Burrito has failed to compile")
		fmt.Println("ERROR", err)
		return
	}

	serv, err := server.NewBurritoServer(&routes, clients)
	if err != nil {
		fmt.Println("The zester has caught a route conflict!, server cannot run")
		fmt.Println(err)
		return
	}
	serv.Run()
}
