package main

import (
	"fmt"
	"os"

	"github.com/Delnegend/common-facade/internal/generator"
)

func main() {
	if err := generator.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
