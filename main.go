package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		if len(os.Args) > 2 {
			fmt.Fprintf(os.Stderr, "Too much arguments: needed 1, given %v\n", len(os.Args)-1)
		} else {
			fmt.Fprintf(os.Stderr, "Too few arguments: needed 1, given %v\n", len(os.Args)-1)
		}
		os.Exit(1)
	}

	// This section is just for tests
	fmt.Fprintf(os.Stdout, "%v\n", os.Args[1])
}
