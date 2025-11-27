package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arjunu/hulud-scan/internal/graph"
	"github.com/arjunu/hulud-scan/internal/parser"
	"github.com/arjunu/hulud-scan/internal/scanner"
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

	// --blocklist flag for blocklist URL or local path
	scanCmd.Flags().String("blocklist",
		"https://github.com/wiz-sec-public/wiz-research-iocs/blob/main/reports/shai-hulud-2-packages.csv",
		"Blocklist URL or local file path")

	// --cache-dir flag for cache directory
	homeDir, _ := os.UserHomeDir()
	defaultCacheDir := filepath.Join(homeDir, ".hulud-scan", "cache")
	scanCmd.Flags().String("cache-dir", defaultCacheDir, "Cache directory for downloaded blocklists")

	// --no-cache flag to disable caching
	scanCmd.Flags().Bool("no-cache", false, "Disable caching (always download fresh)")
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
	fmt.Printf("Project: %s@%s\n", lockfile.Name, lockfile.Version)

	// Step 2: Build dependency graph
	fmt.Println("\n📊 Building dependency graph...")
	dependencyGraph, err := graph.BuildGraph(lockfile)
	if err != nil {
		return fmt.Errorf("failed to build graph: %w", err)
	}

	// Step 3: Load or download blocklist
	blocklistPath, _ := cmd.Flags().GetString("blocklist")
	cacheDir, _ := cmd.Flags().GetString("cache-dir")
	noCache, _ := cmd.Flags().GetBool("no-cache")

	if noCache {
		cacheDir = "" // Disable caching
	}

	fmt.Printf("📋 Loading blocklist from: %s\n", blocklistPath)
	blocklist, err := scanner.LoadOrDownloadBlocklist(blocklistPath, cacheDir)
	if err != nil {
		return fmt.Errorf("failed to load blocklist: %w", err)
	}
	fmt.Printf("✅ Loaded %d blocklist entries\n", len(blocklist.Entries))

	// Step 4: Scan for compromised packages
	fmt.Println("\n🔍 Scanning for compromised packages...")
	result := scanner.ScanGraph(dependencyGraph, blocklist)

	// Step 5: Display results
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("SCAN RESULTS")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	fmt.Printf("Total packages scanned: %d\n", result.TotalPackages)
	fmt.Printf("Issues found: %d\n\n", result.IssuesFound)

	if result.IssuesFound == 0 {
		fmt.Println("✅ No compromised packages detected!")
		return nil
	}

	// Display findings
	fmt.Printf("⚠️  SECURITY ISSUES DETECTED:\n\n")

	for i, finding := range result.Findings {
		fmt.Printf("%d. %s@%s [%s]\n", i+1, finding.PackageName, finding.Version, strings.ToUpper(string(finding.Severity)))

		// Show dependency path
		pathStr := strings.Join(finding.Path, " → ")
		dependencyType := "transitive"
		if finding.IsDirect {
			dependencyType = "direct"
		}
		fmt.Printf("   Type: %s dependency\n", dependencyType)
		fmt.Printf("   Path: %s\n", pathStr)
		fmt.Printf("   Reason: %s\n", finding.Reason)

		if finding.CVE != "" {
			fmt.Printf("   CVE: %s\n", finding.CVE)
		}

		fmt.Println()
	}

	// Exit with error code if critical issues found
	hasCritical := false
	for _, finding := range result.Findings {
		if finding.Severity == scanner.SeverityCritical {
			hasCritical = true
			break
		}
	}

	if hasCritical {
		fmt.Println("❌ Critical security issues detected!")
		return fmt.Errorf("scan failed: critical issues found")
	}

	return nil
}
