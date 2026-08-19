package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	userName := flag.String("username", "momin-chezgi", "A valid user-name in GitHub")

	flag.Parse()

	if len(os.Args) > 3 {
		fmt.Fprintf(os.Stderr, "Too much arguments: needed 1, given %v\n", len(os.Args)-2)
		os.Exit(1)
	}

	// This section is just for tests
	fmt.Fprintf(os.Stdout, "%v\n", *userName)
}
