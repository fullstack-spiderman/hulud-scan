package scanner

import (
	"testing"

	"github.com/fullstack-spiderman/hulud-scan/internal/graph"
	"github.com/fullstack-spiderman/hulud-scan/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBlocklist(t *testing.T) {
	// Arrange
	blocklistPath := "../../testdata/sample-blocklist.csv"

	// Act
	blocklist, err := LoadBlocklist(blocklistPath)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, blocklist)

	// Check that entries were loaded
	assert.NotEmpty(t, blocklist.Entries, "Should have loaded entries")

	// Check specific entry
	assert.GreaterOrEqual(t, len(blocklist.Entries), 1, "Should have at least 1 entry")

	// Find lodash entry
	var lodashEntry *BlocklistEntry
	for i := range blocklist.Entries {
		if blocklist.Entries[i].PackageName == "lodash" {
			lodashEntry = &blocklist.Entries[i]
			break
		}
	}

	require.NotNil(t, lodashEntry, "Should have lodash in blocklist")
	assert.Equal(t, "4.17.20", lodashEntry.Version)
	assert.Equal(t, SeverityCritical, lodashEntry.Severity)
	assert.Contains(t, lodashEntry.Reason, "Prototype pollution")
}

func TestScanGraph_WithCompromisedPackage(t *testing.T) {
	// Arrange - Create a compromised package in our test data
	// First, create a test lockfile with lodash 4.17.20 (blocklisted version)
	lockfile := &parser.Lockfile{
		Name:    "test-app",
		Version: "1.0.0",
		DirectDependencies: map[string]string{
			"lodash": "4.17.20",
		},
		Packages: map[string]*parser.Package{
			"node_modules/lodash": {
				Name:    "lodash",
				Version: "4.17.20", // This version is in blocklist!
			},
		},
	}

	// Build graph
	g, err := graph.BuildGraph(lockfile)
	require.NoError(t, err)

	// Load blocklist
	blocklist, err := LoadBlocklist("../../testdata/sample-blocklist.csv")
	require.NoError(t, err)

	// Act - Scan the graph
	result := ScanGraph(g, blocklist)

	// Assert
	require.NotNil(t, result)
	assert.Equal(t, 1, result.TotalPackages, "Should have scanned 1 package")
	assert.Equal(t, 1, result.IssuesFound, "Should have found 1 issue")
	assert.Len(t, result.Findings, 1, "Should have 1 finding")

	// Check the finding
	finding := result.Findings[0]
	assert.Equal(t, "lodash", finding.PackageName)
	assert.Equal(t, "4.17.20", finding.Version)
	assert.Equal(t, SeverityCritical, finding.Severity)
	assert.True(t, finding.IsDirect, "lodash is a direct dependency")
	assert.Contains(t, finding.Reason, "Prototype pollution")
}

func TestScanGraph_CleanPackages(t *testing.T) {
	// Test with packages not in blocklist
	lockfile := &parser.Lockfile{
		Name:    "test-app",
		Version: "1.0.0",
		DirectDependencies: map[string]string{
			"lodash": "4.17.21", // Different version (not blocklisted)
		},
		Packages: map[string]*parser.Package{
			"node_modules/lodash": {
				Name:    "lodash",
				Version: "4.17.21", // Clean version
			},
		},
	}

	g, err := graph.BuildGraph(lockfile)
	require.NoError(t, err)

	blocklist, err := LoadBlocklist("../../testdata/sample-blocklist.csv")
	require.NoError(t, err)

	// Act
	result := ScanGraph(g, blocklist)

	// Assert
	assert.Equal(t, 1, result.TotalPackages)
	assert.Equal(t, 0, result.IssuesFound, "Should have no issues")
	assert.Empty(t, result.Findings, "Should have no findings")
}

func TestBlocklist_IsBlocked(t *testing.T) {
	// Test the helper method for checking if a package is blocked
	blocklist, err := LoadBlocklist("../../testdata/sample-blocklist.csv")
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name        string
		packageName string
		version     string
		shouldBlock bool
	}{
		{
			name:        "blocked version",
			packageName: "lodash",
			version:     "4.17.20",
			shouldBlock: true,
		},
		{
			name:        "clean version",
			packageName: "lodash",
			version:     "4.17.21",
			shouldBlock: false,
		},
		{
			name:        "package not in blocklist",
			packageName: "some-safe-package",
			version:     "1.0.0",
			shouldBlock: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := blocklist.IsBlocked(tt.packageName, tt.version)
			if tt.shouldBlock {
				assert.NotNil(t, entry, "Should be blocked")
			} else {
				assert.Nil(t, entry, "Should not be blocked")
			}
		})
	}
}
