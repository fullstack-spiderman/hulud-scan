package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLockfile_Simple(t *testing.T) {
	// Arrange - Setup test data
	lockfilePath := "../../testdata/npm/clean/package-lock.json"

	// Act - Execute the function we're testing
	lockfile, err := ParseLockfile(lockfilePath)

	// Assert - Verify the results
	require.NoError(t, err, "ParseLockfile should not return an error")
	require.NotNil(t, lockfile, "Lockfile should not be nil")

	// Check basic metadata
	assert.Equal(t, "test-clean", lockfile.Name)
	assert.Equal(t, "1.0.0", lockfile.Version)
	assert.Equal(t, 3, lockfile.LockfileVersion)

	// Check that packages were parsed
	assert.NotEmpty(t, lockfile.Packages, "Packages map should not be empty")
	assert.Len(t, lockfile.Packages, 4, "Should have 4 packages: lodash, axios, follow-redirects, form-data")

	// Check specific packages exist
	lodash := lockfile.Packages["node_modules/lodash"]
	require.NotNil(t, lodash, "lodash package should exist")
	assert.Equal(t, "lodash", lodash.Name)
	assert.Equal(t, "4.17.21", lodash.Version)

	axios := lockfile.Packages["node_modules/axios"]
	require.NotNil(t, axios, "axios package should exist")
	assert.Equal(t, "axios", axios.Name)
	assert.Equal(t, "1.6.0", axios.Version)
	assert.Contains(t, axios.Dependencies, "follow-redirects")
	assert.Contains(t, axios.Dependencies, "form-data")
}

func TestParseLockfile_FileNotFound(t *testing.T) {
	// Test error handling when file doesn't exist
	lockfilePath := "../../testdata/nonexistent.json"

	lockfile, err := ParseLockfile(lockfilePath)

	assert.Error(t, err, "Should return error for nonexistent file")
	assert.Nil(t, lockfile, "Lockfile should be nil on error")
	assert.Contains(t, err.Error(), "failed to read lockfile")
}

func TestExtractPackageName(t *testing.T) {
	// Table-driven test - a Go testing pattern!
	tests := []struct {
		name     string // Test case name
		path     string // Input
		expected string // Expected output
	}{
		{
			name:     "simple package",
			path:     "node_modules/lodash",
			expected: "lodash",
		},
		{
			name:     "scoped package",
			path:     "node_modules/@babel/core",
			expected: "@babel/core",
		},
		{
			name:     "nested dependency",
			path:     "node_modules/express/node_modules/body-parser",
			expected: "express",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Run creates a subtest
			result := extractPackageName(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}
