package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:\n$ isso problem.json")
		os.Exit(0)
	}

	file := os.Args[1]
	_ = file
}
