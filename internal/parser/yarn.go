package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ParseYarnLock parses a yarn.lock file
func ParseYarnLock(lockfilePath string) (*Lockfile, error) {
	file, err := os.Open(lockfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open yarn.lock: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close file: %v\n", closeErr)
		}
	}()

	lockfile := &Lockfile{
		Name:               extractProjectNameFromPath(lockfilePath),
		Version:            "unknown",
		LockfileVersion:    1, // Yarn lockfile v1
		Packages:           make(map[string]*Package),
		DirectDependencies: make(map[string]string),
	}

	scanner := bufio.NewScanner(file)
	var currentPackage *Package
	var currentPath string
	var inDependencies bool

	// Regex patterns
	// Package line must not start with whitespace (to exclude "dependencies:" etc)
	packageLineRe := regexp.MustCompile(`^([^:\s][^:]*?):\s*$`)
	versionRe := regexp.MustCompile(`^\s+version\s+"([^"]+)"`)
	resolvedRe := regexp.MustCompile(`^\s+resolved\s+"([^"]+)"`)
	integrityRe := regexp.MustCompile(`^\s+integrity\s+(.+)`)
	dependenciesRe := regexp.MustCompile(`^\s+dependencies:\s*$`)
	depEntryRe := regexp.MustCompile(`^\s+([^\s]+)\s+"([^"]+)"`)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// New package entry
		if match := packageLineRe.FindStringSubmatch(line); match != nil {
			// Save previous package if exists
			if currentPackage != nil && currentPath != "" {
				lockfile.Packages[currentPath] = currentPackage
			}

			// Parse package name from "package@version:" format
			// Can be "package@^1.0.0:" or "package@npm:other@1.0.0:"
			packageSpec := strings.TrimSpace(match[1])
			packageName := extractPackageNameFromSpec(packageSpec)

			currentPath = "node_modules/" + packageName
			currentPackage = &Package{
				Name:         packageName,
				Dependencies: make(map[string]string),
			}
			inDependencies = false
			continue
		}

		if currentPackage == nil {
			continue
		}

		// Version
		if match := versionRe.FindStringSubmatch(line); match != nil {
			currentPackage.Version = match[1]
			continue
		}

		// Resolved URL
		if match := resolvedRe.FindStringSubmatch(line); match != nil {
			currentPackage.Resolved = match[1]
			continue
		}

		// Integrity
		if match := integrityRe.FindStringSubmatch(line); match != nil {
			currentPackage.Integrity = match[1]
			continue
		}

		// Dependencies section
		if dependenciesRe.MatchString(line) {
			inDependencies = true
			continue
		}

		// Dependency entry
		if inDependencies {
			if match := depEntryRe.FindStringSubmatch(line); match != nil {
				depName := match[1]
				depVersion := match[2]
				currentPackage.Dependencies[depName] = depVersion
			} else if !strings.HasPrefix(line, "  ") {
				// End of dependencies section
				inDependencies = false
			}
		}
	}

	// Save last package
	if currentPackage != nil && currentPath != "" {
		lockfile.Packages[currentPath] = currentPackage
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading yarn.lock: %w", err)
	}

	// Try to read package.json for project name and direct dependencies
	if err := enrichFromPackageJSON(lockfilePath, lockfile); err != nil {
		// Non-fatal, continue without enrichment
	}

	return lockfile, nil
}

// extractPackageNameFromSpec extracts package name from yarn spec
// Examples:
//
//	"lodash@^4.17.20" -> "lodash"
//	"@babel/core@^7.0.0" -> "@babel/core"
//	"package@npm:other@1.0.0" -> "package"
func extractPackageNameFromSpec(spec string) string {
	// Remove quotes
	spec = strings.Trim(spec, "\"'")

	// Handle npm: aliases
	if strings.Contains(spec, "@npm:") {
		parts := strings.SplitN(spec, "@npm:", 2)
		spec = strings.TrimPrefix(parts[0], "@")
		return spec
	}

	// For scoped packages (@org/package@version)
	if strings.HasPrefix(spec, "@") {
		parts := strings.SplitN(spec[1:], "@", 2)
		if len(parts) > 0 {
			return "@" + parts[0]
		}
	}

	// Regular packages (package@version)
	parts := strings.SplitN(spec, "@", 2)
	return parts[0]
}

// extractProjectNameFromPath extracts project name from lockfile path
func extractProjectNameFromPath(path string) string {
	// Get parent directory name using filepath for cross-platform compatibility
	dir := filepath.Dir(path)
	return filepath.Base(dir)
}

// enrichFromPackageJSON reads package.json to get project info and direct deps
func enrichFromPackageJSON(lockfilePath string, lockfile *Lockfile) error {
	// Find package.json in same directory
	dir := filepath.Dir(lockfilePath)
	pacakageJSONPath := filepath.Join(dir, "package.json")

	data, err := os.ReadFile(pacakageJSONPath)
	if err != nil {
		return err // Non-fatal
	}

	var pkgJSON struct {
		Name         string            `json:"name"`
		Version      string            `json:"version"`
		Dependencies map[string]string `json:"dependencies"`
	}

	if err := json.Unmarshal(data, &pkgJSON); err != nil {
		return err
	}

	// Enrich lockfile
	if pkgJSON.Name != "" {
		lockfile.Name = pkgJSON.Name
	}
	if pkgJSON.Version != "" {
		lockfile.Version = pkgJSON.Version
	}
	if len(pkgJSON.Dependencies) > 0 {
		lockfile.DirectDependencies = pkgJSON.Dependencies
	}

	return nil
}
