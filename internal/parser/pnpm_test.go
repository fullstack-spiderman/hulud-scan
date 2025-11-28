package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePNPMLock(t *testing.T) {
	lockfilePath := "../../testdata/pnpm/clean/pnpm-lock.yaml"

	lockfile, err := ParsePNPMLock(lockfilePath)
	require.NoError(t, err)
	require.NotNil(t, lockfile)

	// Check project metadata
	assert.Equal(t, "test-pnpm-clean", lockfile.Name)
	assert.Equal(t, "1.0.0", lockfile.Version)
	assert.Equal(t, 6, lockfile.LockfileVersion)

	// Check packages
	assert.Len(t, lockfile.Packages, 4, "Should have 4 packages: lodash, axios, follow-redirects, form-data")

	// Check lodash
	lodashPkg := lockfile.Packages["node_modules/lodash"]
	require.NotNil(t, lodashPkg)
	assert.Equal(t, "lodash", lodashPkg.Name)
	assert.Equal(t, "4.17.21", lodashPkg.Version)
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

func TestExtractPNPMPackageInfo(t *testing.T) {
	tests := []struct {
		name            string
		pkgPath         string
		expectedName    string
		expectedVersion string
	}{
		{
			name:            "simple package",
			pkgPath:         "/lodash/4.17.21",
			expectedName:    "lodash",
			expectedVersion: "4.17.21",
		},
		{
			name:            "scoped package",
			pkgPath:         "/@babel/core/7.20.0",
			expectedName:    "@babel/core",
			expectedVersion: "7.20.0",
		},
		{
			name:            "package with peer deps suffix",
			pkgPath:         "/react-dom/18.2.0_react@18.2.0",
			expectedName:    "react-dom",
			expectedVersion: "18.2.0_react@18.2.0",
		},
		{
			name:            "scoped with peer deps",
			pkgPath:         "/@testing-library/react/13.4.0_react@18.2.0",
			expectedName:    "@testing-library/react",
			expectedVersion: "13.4.0_react@18.2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, version := extractPNPMPackageInfo(tt.pkgPath)
			assert.Equal(t, tt.expectedName, name)
			assert.Equal(t, tt.expectedVersion, version)
		})
	}
}

func TestParsePNPMLock_FileNotFound(t *testing.T) {
	_, err := ParsePNPMLock("../../testdata/nonexistent/pnpm-lock.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read pnpm-lock.yaml")
}

func TestParsePNPMLockfileVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{
			name:     "integer version",
			input:    6,
			expected: 6,
		},
		{
			name:     "float version",
			input:    5.4,
			expected: 5,
		},
		{
			name:     "string version",
			input:    "6.0",
			expected: 6,
		},
		{
			name:     "invalid type",
			input:    true,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePNPMLockfileVersion(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
