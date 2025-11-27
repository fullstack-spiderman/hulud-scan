package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arjunu/hulud-scan/internal/parser"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan a project for compromised dependencies",
	Long: `Scan analyzes the package-lock.json file in the specified directory
and checks for known compromised packages and suspicious lifecycle scripts.`,
	Args: cobra.MaximumNArgs(1), // Accept 0 or 1 arguments
	Run: func(cmd *cobra.Command, args []string) {
		// This function runs when the command is executed
		path := "."  // Default to current directory
		if len(args) > 0 {
			path = args[0]
		}

		fmt.Printf("🔍 Scanning project at: %s\n", path)

		// Run the scan
		if err := runScan(path, cmd); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Register scanCmd as a subcommand of rootCmd
	rootCmd.AddCommand(scanCmd)

	// Add flags specific to the scan command
	// --format flag for output format (table or json)
	scanCmd.Flags().StringP("format", "f", "table", "Output format (table or json)")

	// --config flag for custom config file
	scanCmd.Flags().StringP("config", "c", "", "Path to config file")
}

// runScan performs the actual scanning logic
func runScan(projectPath string, cmd *cobra.Command) error {
	// Find package-lock.json in the project directory
	lockfilePath := filepath.Join(projectPath, "package-lock.json")

	// Check if file exists
	if _, err := os.Stat(lockfilePath); os.IsNotExist(err) {
		return fmt.Errorf("package-lock.json not found in %s", projectPath)
	}

	fmt.Printf("📄 Parsing lockfile: %s\n", lockfilePath)

	// Parse the lockfile
	lockfile, err := parser.ParseLockfile(lockfilePath)
	if err != nil {
		return fmt.Errorf("failed to parse lockfile: %w", err)
	}

	fmt.Printf("✅ Found %d packages\n", len(lockfile.Packages))
	fmt.Printf("\nProject: %s@%s\n", lockfile.Name, lockfile.Version)
	fmt.Printf("Lockfile version: %d\n\n", lockfile.LockfileVersion)

	// Print first 5 packages as a sample
	fmt.Println("Sample packages:")
	count := 0
	for path, pkg := range lockfile.Packages {
		if count >= 5 {
			break
		}
		fmt.Printf("  - %s@%s (path: %s)\n", pkg.Name, pkg.Version, path)
		count++
	}

	fmt.Println("\n✨ TODO: Implement blocklist checking and script detection")

	return nil
}
