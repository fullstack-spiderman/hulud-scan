package parser

import (
	"fmt"
	"os"
	"path/filepath"
)

// LockfileType represents the type of package manager lockfile
type LockfileType string

const (
	LockfileTypeNPM  LockfileType = "npm"
	LockfileTypeYarn LockfileType = "yarn"
	LockfileTypePNPM LockfileType = "pnpm"
	LockfileTypeBun  LockfileType = "bun"
)

// LockfileInfo contains detected lockfile information
type LockfileInfo struct {
	Type     LockfileType
	Path     string
	Filename string
}

// DetectLockfile detects which lockfile exists in the project directory
// Priority order: package-lock.json > yarn.lock > pnpm-lock.yaml > bun.lockb
func DetectLockfile(projectPath string) (*LockfileInfo, error) {
	// Check in priority order (npm is most common)
	lockfiles := []struct {
		filename string
		lockType LockfileType
	}{
		{"package-lock.json", LockfileTypeNPM},
		{"yarn.lock", LockfileTypeYarn},
		{"pnpm-lock.yaml", LockfileTypePNPM},
		{"bun.lockb", LockfileTypeBun},
	}

	for _, lf := range lockfiles {
		lockfilePath := filepath.Join(projectPath, lf.filename)
		if _, err := os.Stat(lockfilePath); err == nil {
			return &LockfileInfo{
				Type:     lf.lockType,
				Path:     lockfilePath,
				Filename: lf.filename,
			}, nil
		}
	}

	return nil, fmt.Errorf("no supported lockfile found in %s (looking for: package-lock.json, yarn.lock, pnpm-lock.yaml, bun.lockb)", projectPath)
}

// String returns a human-readable name for the lockfile type
func (t LockfileType) String() string {
	switch t {
	case LockfileTypeNPM:
		return "npm (package-lock.json)"
	case LockfileTypeYarn:
		return "Yarn (yarn.lock)"
	case LockfileTypePNPM:
		return "pnpm (pnpm-lock.yaml)"
	case LockfileTypeBun:
		return "Bun (bun.lockb)"
	default:
		return string(t)
	}
}

// ParseAuto automatically detects and parses the appropriate lockfile
func ParseAuto(projectPath string) (*Lockfile, *LockfileInfo, error) {
	// Detect which lockfile exists
	info, err := DetectLockfile(projectPath)
	if err != nil {
		return nil, nil, err
	}

	// Parse based on detected type
	var lockfile *Lockfile
	switch info.Type {
	case LockfileTypeNPM:
		lockfile, err = ParseLockfile(info.Path)
	case LockfileTypeYarn:
		lockfile, err = ParseYarnLock(info.Path)
	case LockfileTypePNPM:
		lockfile, err = ParsePNPMLock(info.Path)
	case LockfileTypeBun:
		lockfile, err = ParseBunLock(info.Path)
	default:
		return nil, nil, fmt.Errorf("unsupported lockfile type: %s", info.Type)
	}

	if err != nil {
		return nil, info, err
	}

	return lockfile, info, nil
}
