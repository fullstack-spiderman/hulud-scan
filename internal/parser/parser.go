package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseLockfile reads and parses a package-lock.json file
func ParseLockfile(lockfilePath string) (*Lockfile, error) {
	// Read the file
	data, err := os.ReadFile(lockfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lockfile: %w", err)
	}

	// Parse JSON into a temporary structure that matches package-lock.json format
	var raw struct {
		Name            string `json:"name"`
		Version         string `json:"version"`
		LockfileVersion int    `json:"lockfileVersion"`
		Packages        map[string]struct {
			Version      string            `json:"version"`
			Resolved     string            `json:"resolved"`
			Integrity    string            `json:"integrity"`
			Dependencies map[string]string `json:"dependencies"`
		} `json:"packages"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse lockfile JSON: %w", err)
	}

	// Extract direct dependencies from root package (empty string key)
	directDeps := make(map[string]string)
	if rootPkg, exists := raw.Packages[""]; exists {
		directDeps = rootPkg.Dependencies
	}

	// Convert to our internal Lockfile structure
	lockfile := &Lockfile{
		Name:               raw.Name,
		Version:            raw.Version,
		LockfileVersion:    raw.LockfileVersion,
		Packages:           make(map[string]*Package),
		DirectDependencies: directDeps,
	}

	// Process each package
	for path, pkg := range raw.Packages {
		// Skip the root package (empty string key)
		if path == "" {
			continue
		}

		// Extract package name from path (e.g., "node_modules/lodash" -> "lodash")
		name := extractPackageName(path)

		lockfile.Packages[path] = &Package{
			Name:         name,
			Version:      pkg.Version,
			Resolved:     pkg.Resolved,
			Integrity:    pkg.Integrity,
			Dependencies: pkg.Dependencies,
		}
	}

	return lockfile, nil
}

// extractPackageName extracts the package name from a node_modules path
// e.g., "node_modules/lodash" -> "lodash"
// e.g., "node_modules/@babel/core" -> "@babel/core"
func extractPackageName(path string) string {
	// Remove "node_modules/" prefix
	name := strings.TrimPrefix(path, "node_modules/")

	// Handle scoped packages (e.g., "@babel/core")
	if strings.HasPrefix(name, "@") {
		// For scoped packages, keep the scope and package name
		parts := strings.Split(name, "/")
		if len(parts) >= 2 {
			return filepath.Join(parts[0], parts[1])
		}
	}

	// For regular packages, just return the first part
	parts := strings.Split(name, "/")
	return parts[0]
}
