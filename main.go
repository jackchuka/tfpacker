package main

import (
	"fmt"
	"os"

	"github.com/jackchuka/tfpacker/cmd/tfpacker"
)

func main() {
	if err := tfpacker.Eexecute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
