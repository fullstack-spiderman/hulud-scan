package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hulud-scan",
	Short: "A supply-chain security scanner for JavaScript/TypeScript projects",
	Long: `hulud-scan analyzes your project's dependencies (from lockfiles)
and detects known compromised packages and suspicious lifecycle scripts.

Examples:
  hulud-scan scan ./my-project
  hulud-scan scan --format json ./my-project`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// init() runs automatically when the package is imported
	// We'll add global flags here later
}
