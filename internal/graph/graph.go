package graph

import (
	"github.com/arjunu/hulud-scan/internal/parser"
)

// BuildGraph constructs a dependency graph from a parsed lockfile
func BuildGraph(lockfile *parser.Lockfile) (*Graph, error) {
	// Create the graph
	graph := &Graph{
		Nodes: make(map[string]*Node),
	}

	// Create root node (the project itself)
	// The root doesn't have a Package entry in the lockfile's Packages map
	graph.Root = &Node{
		Package: &parser.Package{
			Name:    lockfile.Name,
			Version: lockfile.Version,
		},
		Dependencies: make([]*Node, 0),
		Dependents:   make([]*Node, 0),
		IsDirect:     false, // Root is not a dependency
		Depth:        0,     // Root is at depth 0
	}

	// Step 1: Create nodes for all packages
	for path, pkg := range lockfile.Packages {
		graph.Nodes[path] = &Node{
			Package:      pkg,
			Dependencies: make([]*Node, 0),
			Dependents:   make([]*Node, 0),
			IsDirect:     false, // Will be set later
			Depth:        -1,    // Will be calculated later
		}
	}

	// Step 2: Build dependency edges
	// For each package, link it to its dependencies
	for path, node := range graph.Nodes {
		pkg := lockfile.Packages[path]

		// For each dependency this package has
		for depName, depVersion := range pkg.Dependencies {
			// Find the dependency node
			// Dependencies are typically in node_modules/depName
			depPath := "node_modules/" + depName

			depNode, exists := graph.Nodes[depPath]
			if exists {
				// Add edge: node -> depNode
				node.Dependencies = append(node.Dependencies, depNode)

				// Add reverse edge: depNode knows node depends on it
				depNode.Dependents = append(depNode.Dependents, node)
			}
			// Note: If not exists, it might be a nested dependency
			// For now, we skip those (they're in the lockfile differently)
			// We can enhance this later
			_ = depVersion // Avoid unused variable error
		}
	}

	// Step 3: Calculate depth and mark direct dependencies
	// We'll use BFS (Breadth-First Search) from the root
	calculateDepth(graph, lockfile)

	return graph, nil
}

// calculateDepth performs BFS to calculate depth and mark direct dependencies
func calculateDepth(graph *Graph, lockfile *parser.Lockfile) {
	// Queue for BFS: [node, depth]
	type queueItem struct {
		node  *Node
		depth int
	}

	queue := []queueItem{}

	// Find all direct dependencies by looking at root project's dependencies
	// In package-lock.json v3, the root package has empty string key with dependencies
	// But we're building from our Root node, so we need to find which packages
	// are direct vs transitive

	// Direct dependencies are those at depth 1 from root
	// We'll start by finding packages not depended on by other packages
	// Actually, better approach: assume all packages are direct unless proven otherwise

	// Mark packages that are depended on by other packages as transitive
	transitivePackages := make(map[string]bool)
	for _, node := range graph.Nodes {
		for _, dep := range node.Dependencies {
			// This dep is depended on by node, so it's transitive from node's perspective
			// But we need to check if node itself is direct
			// This is getting complex - let's use a simpler approach
			transitivePackages[dep.Package.Name] = true
		}
	}

	// Actually, simpler approach: all packages at depth 1 are direct
	// Let's do BFS and mark as we go

	// Start with root's dependencies
	// We need to infer root's dependencies from the graph
	// Root depends on packages that aren't depended on by any package in node_modules
	// Or: packages at the top level of node_modules

	// Even simpler: just traverse from root, marking depth
	// Direct = depth 1

	// Initialize: all packages start at depth -1
	// Root is depth 0
	queue = append(queue, queueItem{node: graph.Root, depth: 0})
	visited := make(map[*Node]bool)

	// We need to connect root to its direct dependencies
	// In our lockfile, direct dependencies are those not nested
	// Let's assume all packages in graph.Nodes are potentially reachable
	// and find shortest path to each

	// Better approach: Use BFS from each package upward through Dependents
	// to find distance to root

	// Actually, let's use a working approach:
	// 1. Top-level packages (node_modules/X without further nesting in path) are direct
	// 2. Calculate depth from root

	// Mark direct dependencies using the lockfile's DirectDependencies map
	for depName := range lockfile.DirectDependencies {
		// Find the node for this direct dependency
		depPath := "node_modules/" + depName
		if node, exists := graph.Nodes[depPath]; exists {
			node.IsDirect = true
			node.Depth = 1

			// Connect root to this direct dependency
			graph.Root.Dependencies = append(graph.Root.Dependencies, node)
			node.Dependents = append(node.Dependents, graph.Root)
		}
	}

	// Now calculate depth for transitive dependencies using BFS
	queue = []queueItem{{node: graph.Root, depth: 0}}
	visited[graph.Root] = true

	for len(queue) > 0 {
		// Dequeue
		item := queue[0]
		queue = queue[1:]

		currentNode := item.node
		currentDepth := item.depth

		// Process all dependencies
		for _, dep := range currentNode.Dependencies {
			if !visited[dep] {
				visited[dep] = true

				// Set depth if not already set or if we found a shorter path
				if dep.Depth == -1 || dep.Depth > currentDepth+1 {
					dep.Depth = currentDepth + 1
				}

				// Enqueue for further processing
				queue = append(queue, queueItem{node: dep, depth: currentDepth + 1})
			}
		}
	}

	// Handle any unvisited nodes (shouldn't happen in a well-formed graph)
	for _, node := range graph.Nodes {
		if node.Depth == -1 {
			node.Depth = 999 // Mark as unreachable
		}
	}
}

// FindPath finds a dependency path from root to the specified package
// Returns the path as a list of package names
func (g *Graph) FindPath(targetPath string) DependencyPath {
	targetNode, exists := g.Nodes[targetPath]
	if !exists {
		return nil
	}

	// Use BFS to find shortest path from root to target
	type queueItem struct {
		node *Node
		path DependencyPath
	}

	queue := []queueItem{
		{node: g.Root, path: DependencyPath{g.Root.Package.Name}},
	}

	visited := make(map[*Node]bool)
	visited[g.Root] = true

	for len(queue) > 0 {
		// Dequeue
		item := queue[0]
		queue = queue[1:]

		currentNode := item.node
		currentPath := item.path

		// Check if we reached the target
		if currentNode == targetNode {
			return currentPath
		}

		// Explore dependencies
		for _, dep := range currentNode.Dependencies {
			if !visited[dep] {
				visited[dep] = true

				// Create new path including this dependency
				newPath := make(DependencyPath, len(currentPath)+1)
				copy(newPath, currentPath)
				newPath[len(currentPath)] = dep.Package.Name

				queue = append(queue, queueItem{node: dep, path: newPath})
			}
		}
	}

	// Target not reachable from root
	return nil
}
