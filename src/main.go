package main

import (
	"fmt"
	"os"

	"./burrito"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main <burrito_file>")
		return
	}
	routes, err := burrito.ParseBurritoFile(os.Args[1])
	if err != nil {
		fmt.Println("Your Burrito has failed to compile")
		fmt.Println("ERROR", err)
	}

	serv := burrito.BurritoServer{
		Routes: &routes,
	}
	serv.Run()
}
