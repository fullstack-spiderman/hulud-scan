package parser

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ParseBunLock parses a bun.lockb file
// Note: bun.lockb is a binary format. We use `bun` CLI to read it.
func ParseBunLock(lockfilePath string) (*Lockfile, error) {
	// Check if bun is installed
	if !isBunInstalled() {
		return nil, fmt.Errorf("bun.lockb detected but 'bun' CLI is not installed or not in PATH\n" +
			"Please install Bun from https://bun.sh or use a different package manager")
	}

	// Get the project directory
	projectDir := filepath.Dir(lockfilePath)

	// Use `bun pm ls --all` to get package list
	// This outputs JSON with package information
	cmd := exec.Command("bun", "pm", "ls", "--all")
	cmd.Dir = projectDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run 'bun pm ls': %w\nOutput: %s", err, string(output))
	}

	// Parse the output
	// Bun outputs each package on a line
	lockfile := &Lockfile{
		Name:               extractProjectNameFromPath(lockfilePath),
		Version:            "unknown",
		LockfileVersion:    1,
		Packages:           make(map[string]*Package),
		DirectDependencies: make(map[string]string),
	}

	// Parse bun pm ls output
	// Format is typically: package@version
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "node_modules") {
			continue
		}

		// Try to parse package@version format
		if strings.Contains(line, "@") {
			name, version := parseBunPackageLine(line)
			if name != "" && version != "" {
				nodePath := "node_modules/" + name
				lockfile.Packages[nodePath] = &Package{
					Name:         name,
					Version:      version,
					Dependencies: make(map[string]string),
				}
			}
		}
	}

	// Enrich from package.json (error is non-fatal)
	_ = enrichFromPackageJSON(lockfilePath, lockfile)

	return lockfile, nil
}

// isBunInstalled checks if bun CLI is available
func isBunInstalled() bool {
	_, err := exec.LookPath("bun")
	return err == nil
}

// parseBunPackageLine parses a line from bun pm ls output
// Examples:
//
//	"lodash@4.17.21" -> "lodash", "4.17.21"
//	"@babel/core@7.20.0" -> "@babel/core", "7.20.0"
func parseBunPackageLine(line string) (name string, version string) {
	line = strings.TrimSpace(line)

	// Remove tree characters and whitespace
	line = strings.TrimLeft(line, " ├─└│")
	line = strings.TrimSpace(line)

	// Handle scoped packages
	if strings.HasPrefix(line, "@") {
		// Find the last @ which separates version
		lastAt := strings.LastIndex(line, "@")
		if lastAt > 0 {
			name = line[:lastAt]
			version = line[lastAt+1:]
			return
		}
	}

	// Regular packages
	parts := strings.SplitN(line, "@", 2)
	if len(parts) == 2 {
		name = parts[0]
		version = parts[1]
		return
	}

	return "", ""
}
