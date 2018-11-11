package main

import (
	"fmt"
	"os"

	"./parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main <burrito_file>")
		return
	}
	fmt.Println(parser.ParseBurritoFile(os.Args[1]))

}
