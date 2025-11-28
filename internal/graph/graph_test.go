package graph

import (
	"testing"

	"github.com/fullstack-spiderman/hulud-scan/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildGraph_Simple(t *testing.T) {
	// Arrange - Parse our test lockfile
	lockfilePath := "../../testdata/npm-project/package-lock.json"
	lockfile, err := parser.ParseLockfile(lockfilePath)
	require.NoError(t, err)

	// Act - Build the graph
	graph, err := BuildGraph(lockfile)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, graph)

	// Check root node
	require.NotNil(t, graph.Root)
	assert.Equal(t, "test-project", graph.Root.Package.Name)
	assert.Equal(t, 0, graph.Root.Depth, "Root should have depth 0")

	// Check that nodes were created for all packages
	// We have 3 packages in our test file
	assert.Len(t, graph.Nodes, 3, "Should have 3 nodes")

	// Check that direct dependencies exist and are marked correctly
	lodashNode := graph.Nodes["node_modules/lodash"]
	require.NotNil(t, lodashNode, "lodash node should exist")
	assert.Equal(t, "lodash", lodashNode.Package.Name)
	assert.True(t, lodashNode.IsDirect, "lodash should be marked as direct dependency")
	assert.Equal(t, 1, lodashNode.Depth, "Direct dependencies should have depth 1")

	expressNode := graph.Nodes["node_modules/express"]
	require.NotNil(t, expressNode, "express node should exist")
	assert.Equal(t, "express", expressNode.Package.Name)
	assert.True(t, expressNode.IsDirect, "express should be marked as direct dependency")
	assert.Equal(t, 1, expressNode.Depth)

	// Check transitive dependency
	bodyParserNode := graph.Nodes["node_modules/body-parser"]
	require.NotNil(t, bodyParserNode, "body-parser node should exist")
	assert.Equal(t, "body-parser", bodyParserNode.Package.Name)
	assert.False(t, bodyParserNode.IsDirect, "body-parser should NOT be direct dependency")
	assert.Equal(t, 2, bodyParserNode.Depth, "Transitive dependencies should have depth > 1")

	// Check that express has body-parser as a dependency
	assert.Len(t, expressNode.Dependencies, 1, "express should have 1 dependency")
	assert.Equal(t, bodyParserNode, expressNode.Dependencies[0], "express should depend on body-parser")

	// Check that body-parser knows express depends on it
	assert.Contains(t, bodyParserNode.Dependents, expressNode, "body-parser should list express as dependent")
}

func TestBuildGraph_EmptyLockfile(t *testing.T) {
	// Test edge case: lockfile with no packages
	lockfile := &parser.Lockfile{
		Name:     "empty-project",
		Version:  "1.0.0",
		Packages: make(map[string]*parser.Package),
	}

	graph, err := BuildGraph(lockfile)

	require.NoError(t, err)
	require.NotNil(t, graph)
	assert.NotNil(t, graph.Root)
	assert.Equal(t, "empty-project", graph.Root.Package.Name)
	assert.Empty(t, graph.Nodes, "Should have no dependency nodes")
}

func TestFindDependencyPath(t *testing.T) {
	// Arrange
	lockfilePath := "../../testdata/npm-project/package-lock.json"
	lockfile, err := parser.ParseLockfile(lockfilePath)
	require.NoError(t, err)

	graph, err := BuildGraph(lockfile)
	require.NoError(t, err)

	// Act - Find path from root to body-parser
	path := graph.FindPath("node_modules/body-parser")

	// Assert
	require.NotNil(t, path, "Should find a path to body-parser")

	// Path should be: root -> express -> body-parser
	assert.Len(t, path, 3, "Path should have 3 nodes")
	assert.Equal(t, "test-project", path[0], "First should be root project")
	assert.Equal(t, "express", path[1], "Second should be express")
	assert.Equal(t, "body-parser", path[2], "Last should be body-parser")
}

func TestFindDependencyPath_DirectDependency(t *testing.T) {
	// Test path to a direct dependency
	lockfilePath := "../../testdata/npm-project/package-lock.json"
	lockfile, err := parser.ParseLockfile(lockfilePath)
	require.NoError(t, err)

	graph, err := BuildGraph(lockfile)
	require.NoError(t, err)

	// Act - Find path to lodash (direct dependency)
	path := graph.FindPath("node_modules/lodash")

	// Assert
	require.NotNil(t, path)
	assert.Len(t, path, 2, "Direct dependency path should have 2 nodes")
	assert.Equal(t, "test-project", path[0])
	assert.Equal(t, "lodash", path[1])
}
