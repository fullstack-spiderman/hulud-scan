package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseYarnLock(t *testing.T) {
	lockfilePath := "../../testdata/yarn/clean/yarn.lock"

	lockfile, err := ParseYarnLock(lockfilePath)
	require.NoError(t, err)
	require.NotNil(t, lockfile)

	// Check project metadata
	assert.Equal(t, "test-yarn-clean", lockfile.Name)
	assert.Equal(t, "1.0.0", lockfile.Version)
	assert.Equal(t, 1, lockfile.LockfileVersion)

	// Check packages
	assert.Len(t, lockfile.Packages, 4, "Should have 4 packages: lodash, axios, follow-redirects, form-data")

	// Check lodash
	lodashPkg := lockfile.Packages["node_modules/lodash"]
	require.NotNil(t, lodashPkg)
	assert.Equal(t, "lodash", lodashPkg.Name)
	assert.Equal(t, "4.17.21", lodashPkg.Version)
	assert.Contains(t, lodashPkg.Resolved, "lodash")
	assert.NotEmpty(t, lodashPkg.Integrity)

	// Check axios
	axiosPkg := lockfile.Packages["node_modules/axios"]
	require.NotNil(t, axiosPkg)
	assert.Equal(t, "axios", axiosPkg.Name)
	assert.Equal(t, "1.6.0", axiosPkg.Version)
	assert.NotEmpty(t, axiosPkg.Dependencies)
	assert.Contains(t, axiosPkg.Dependencies, "follow-redirects")
	assert.Contains(t, axiosPkg.Dependencies, "form-data")

	// Check direct dependencies
	assert.Len(t, lockfile.DirectDependencies, 2)
	assert.Contains(t, lockfile.DirectDependencies, "lodash")
	assert.Contains(t, lockfile.DirectDependencies, "axios")
}

func TestExtractPackageNameFromSpec(t *testing.T) {
	tests := []struct {
		name     string
		spec     string
		expected string
	}{
		{
			name:     "simple package",
			spec:     "lodash@^4.17.20",
			expected: "lodash",
		},
		{
			name:     "scoped package",
			spec:     "@babel/core@^7.0.0",
			expected: "@babel/core",
		},
		{
			name:     "npm alias",
			spec:     "package@npm:other@1.0.0",
			expected: "package",
		},
		{
			name:     "scoped with quotes",
			spec:     "\"@types/node@^18.0.0\"",
			expected: "@types/node",
		},
		{
			name:     "exact version",
			spec:     "react@17.0.2",
			expected: "react",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPackageNameFromSpec(tt.spec)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseYarnLock_FileNotFound(t *testing.T) {
	_, err := ParseYarnLock("../../testdata/nonexistent/yarn.lock")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open yarn.lock")
}

func TestExtractProjectNameFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "yarn.lock path",
			path:     "/home/user/projects/my-app/yarn.lock",
			expected: "my-app",
		},
		{
			name:     "package-lock.json path",
			path:     "/var/www/web-app/package-lock.json",
			expected: "web-app",
		},
		{
			name:     "pnpm-lock.yaml path",
			path:     "/projects/backend/pnpm-lock.yaml",
			expected: "backend",
		},
		{
			name:     "bun.lockb path",
			path:     "/code/frontend/bun.lockb",
			expected: "frontend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractProjectNameFromPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}
