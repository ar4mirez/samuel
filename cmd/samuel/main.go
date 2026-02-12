package main

import (
	"fmt"
	"os"

	"github.com/ar4mirez/samuel/internal/commands"
	"github.com/fatih/color"
)

func main() {
	if err := commands.Execute(); err != nil {
		// Print error in red
		red := color.New(color.FgRed).SprintFunc()
		fmt.Fprintf(os.Stderr, "%s %s\n", red("Error:"), err.Error())
		os.Exit(1)
	}
}
