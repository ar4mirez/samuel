package main

import (
	"os"

	"github.com/ar4mirez/aicof/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
