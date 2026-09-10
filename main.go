package main

import (
	"fmt"
	"os"
)

const neededArgs = 1

func main() {

	if len(os.Args) != neededArgs+1 {
		if len(os.Args) > neededArgs+1 {
			fmt.Fprintf(os.Stderr, "Too much arguments: needed %v, given %v\n",neededArgs, len(os.Args)-1)
		} else {
			fmt.Fprintf(os.Stderr, "Too few arguments: needed %v, given %v\n",neededArgs, len(os.Args)-1)
		}
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%v\n", os.Args[neededArgs])
}
