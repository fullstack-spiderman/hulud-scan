package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectLockfile(t *testing.T) {
	tests := []struct {
		name         string
		projectPath  string
		expectedType LockfileType
		shouldFail   bool
	}{
		{
			name:         "detect npm lockfile",
			projectPath:  "../../testdata/npm-project",
			expectedType: LockfileTypeNPM,
			shouldFail:   false,
		},
		{
			name:         "detect yarn lockfile",
			projectPath:  "../../testdata/yarn-project",
			expectedType: LockfileTypeYarn,
			shouldFail:   false,
		},
		{
			name:         "detect pnpm lockfile",
			projectPath:  "../../testdata/pnpm-project",
			expectedType: LockfileTypePNPM,
			shouldFail:   false,
		},
		{
			name:        "no lockfile found",
			projectPath: "../../testdata/nonexistent",
			shouldFail:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := DetectLockfile(tt.projectPath)

			if tt.shouldFail {
				assert.Error(t, err)
				assert.Nil(t, info)
			} else {
				require.NoError(t, err)
				require.NotNil(t, info)
				assert.Equal(t, tt.expectedType, info.Type)
				assert.NotEmpty(t, info.Path)
			}
		})
	}
}

func TestParseAuto(t *testing.T) {
	tests := []struct {
		name             string
		projectPath      string
		expectedName     string
		expectedPackages int
		shouldFail       bool
	}{
		{
			name:             "parse npm project",
			projectPath:      "../../testdata/npm-project",
			expectedName:     "test-project",
			expectedPackages: 3, // lodash, express, body-parser
			shouldFail:       false,
		},
		{
			name:             "parse yarn project",
			projectPath:      "../../testdata/yarn-project",
			expectedName:     "yarn-test-project",
			expectedPackages: 3, // lodash, express, body-parser
			shouldFail:       false,
		},
		{
			name:             "parse pnpm project",
			projectPath:      "../../testdata/pnpm-project",
			expectedName:     "pnpm-test-project",
			expectedPackages: 3, // lodash, express, body-parser
			shouldFail:       false,
		},
		{
			name:        "no lockfile found",
			projectPath: "../../testdata/nonexistent",
			shouldFail:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockfile, info, err := ParseAuto(tt.projectPath)

			if tt.shouldFail {
				assert.Error(t, err)
				assert.Nil(t, lockfile)
				assert.Nil(t, info)
			} else {
				require.NoError(t, err)
				require.NotNil(t, lockfile)
				require.NotNil(t, info)
				assert.Equal(t, tt.expectedName, lockfile.Name)
				assert.Len(t, lockfile.Packages, tt.expectedPackages)
				assert.NotEmpty(t, lockfile.DirectDependencies)
			}
		})
	}
}

func TestLockfileTypeString(t *testing.T) {
	tests := []struct {
		lockType LockfileType
		expected string
	}{
		{LockfileTypeNPM, "npm"},
		{LockfileTypeYarn, "yarn"},
		{LockfileTypePNPM, "pnpm"},
		{LockfileTypeBun, "bun"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.lockType))
		})
	}
}
