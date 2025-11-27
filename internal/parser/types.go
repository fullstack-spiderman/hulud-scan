package parser

// Package represents a single package in the dependency tree
type Package struct {
	Name         string            // Package name (e.g., "lodash")
	Version      string            // Exact version (e.g., "4.17.21")
	Resolved     string            // URL where package was downloaded from
	Integrity    string            // Hash for verification
	Dependencies map[string]string // Direct dependencies (name -> version range)
}

// Lockfile represents the parsed package-lock.json structure
type Lockfile struct {
	Name               string              // Project name
	Version            string              // Project version
	LockfileVersion    int                 // npm lockfile format version
	Packages           map[string]*Package // Map of package path -> Package info
	DirectDependencies map[string]string   // Direct dependencies from root (name -> version range)
}
