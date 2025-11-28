package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information (set by main package)
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
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
	Version: Version,
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hulud-scan %s\n", Version)
		fmt.Printf("Commit:  %s\n", Commit)
		fmt.Printf("Built:   %s\n", Date)
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// init() runs automatically when the package is imported
	rootCmd.AddCommand(versionCmd)
	rootCmd.SetVersionTemplate("hulud-scan version {{.Version}}\n")
}
