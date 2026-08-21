package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.String("username", "momin-chezgi", "A valid user-name in GitHub")

	flag.Parse()

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Too much arguments: needed 1, given %v\n", len(os.Args)-1)
		os.Exit(1)
	}

	// This section is just for tests
	fmt.Fprintf(os.Stdout, "%v\n", os.Args[1])
}
