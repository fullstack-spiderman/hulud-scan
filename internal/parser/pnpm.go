package parser

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParsePNPMLock parses a pnpm-lock.yaml file
func ParsePNPMLock(lockfilePath string) (*Lockfile, error) {
	data, err := os.ReadFile(lockfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pnpm-lock.yaml: %w", err)
	}

	// pnpm lockfile structure
	var pnpmLock struct {
		LockfileVersion interface{} `yaml:"lockfileVersion"`
		Packages        map[string]struct {
			Resolution struct {
				Integrity string `yaml:"integrity"`
			} `yaml:"resolution"`
			Dependencies    map[string]string `yaml:"dependencies"`
			DevDependencies map[string]string `yaml:"devDependencies"`
		} `yaml:"packages"`
	}

	if err := yaml.Unmarshal(data, &pnpmLock); err != nil {
		return nil, fmt.Errorf("failed to parse pnpm-lock.yaml: %w", err)
	}

	lockfile := &Lockfile{
		Name:               extractProjectNameFromPath(lockfilePath),
		Version:            "unknown",
		LockfileVersion:    parsePNPMLockfileVersion(pnpmLock.LockfileVersion),
		Packages:           make(map[string]*Package),
		DirectDependencies: make(map[string]string),
	}

	// Parse packages
	// pnpm format: "/lodash/4.17.21" -> Package
	for pkgPath, pkgData := range pnpmLock.Packages {
		// Extract name and version from path
		// Format: "/package-name/version" or "/@scope/package-name/version"
		name, version := extractPNPMPackageInfo(pkgPath)

		// Skip if we couldn't parse
		if name == "" {
			continue
		}

		// Create standard path format
		nodePath := "node_modules/" + name

		pkg := &Package{
			Name:         name,
			Version:      version,
			Integrity:    pkgData.Resolution.Integrity,
			Dependencies: make(map[string]string),
		}

		// Merge dependencies and devDependencies
		for depName, depVersion := range pkgData.Dependencies {
			pkg.Dependencies[depName] = depVersion
		}
		for depName, depVersion := range pkgData.DevDependencies {
			pkg.Dependencies[depName] = depVersion
		}

		lockfile.Packages[nodePath] = pkg
	}

	// Enrich from package.json
	if err := enrichFromPackageJSON(lockfilePath, lockfile); err != nil {
		// Non-fatal
	}

	return lockfile, nil
}

// extractPNPMPackageInfo extracts package name and version from pnpm path
// Examples:
//   "/lodash/4.17.21" -> "lodash", "4.17.21"
//   "/@babel/core/7.20.0" -> "@babel/core", "7.20.0"
func extractPNPMPackageInfo(pkgPath string) (name string, version string) {
	// Remove leading slash
	pkgPath = strings.TrimPrefix(pkgPath, "/")

	// Handle scoped packages
	if strings.HasPrefix(pkgPath, "@") {
		// Format: @scope/package/version
		parts := strings.Split(pkgPath, "/")
		if len(parts) >= 3 {
			name = parts[0] + "/" + parts[1] // @scope/package
			version = parts[2]
			return
		}
	}

	// Regular packages
	// Format: package/version
	parts := strings.Split(pkgPath, "/")
	if len(parts) >= 2 {
		name = parts[0]
		version = parts[1]
		return
	}

	return "", ""
}

// parsePNPMLockfileVersion converts pnpm lockfile version to int
func parsePNPMLockfileVersion(version interface{}) int {
	switch v := version.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		// Try to parse as int
		var result int
		if _, err := fmt.Sscanf(v, "%d", &result); err != nil {
			// If parsing fails, return 0 as default
			return 0
		}
		return result
	default:
		return 0
	}
}
