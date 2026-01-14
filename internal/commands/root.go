package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information (set at build time)
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "aicof",
	Short: "AICoF - Artificial Intelligence Coding Framework CLI",
	Long: `AICoF CLI manages the Artificial Intelligence Coding Framework.

It helps you initialize projects with AI coding guardrails, update framework
versions, and manage language/framework guides without cloning the repository.

Examples:
  aicof init my-project          # Initialize a new project
  aicof init .                   # Initialize in current directory
  aicof update                   # Update to latest framework version
  aicof add language rust        # Add Rust language guide
  aicof list --available         # List all available components
  aicof doctor                   # Check installation health`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")
}

// exitWithError prints an error message and exits
func exitWithError(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s: %v\n", msg, err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	}
	os.Exit(1)
}
