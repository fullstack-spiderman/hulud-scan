package graph

import "github.com/fullstack-spiderman/hulud-scan/internal/parser"

// Node represents a package in the dependency graph
type Node struct {
	Package      *parser.Package // The actual package data
	Dependencies []*Node         // Packages this one depends on
	Dependents   []*Node         // Packages that depend on this one
	IsDirect     bool            // Is this a direct dependency of the root project?
	Depth        int             // How far from root (0 = direct, 1+ = transitive)
}

// Graph represents the complete dependency graph
type Graph struct {
	Root  *Node            // The root project
	Nodes map[string]*Node // All nodes indexed by package path
}

// DependencyPath represents a chain showing how a package is reached
// Example: ["my-app", "express", "body-parser", "lodash"]
type DependencyPath []string
