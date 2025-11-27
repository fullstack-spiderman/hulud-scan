# Contributing to hulud-scan

Thank you for your interest in contributing! This guide will help you get started with development.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Building & Compiling](#building--compiling)
- [Running Tests](#running-tests)
- [Code Structure](#code-structure)
- [Adding New Features](#adding-new-features)
- [Submitting Changes](#submitting-changes)
- [Style Guide](#style-guide)

---

## Getting Started

### Prerequisites

- **Go 1.21+** - [Download here](https://go.dev/dl/)
- **Git** - [Download here](https://git-scm.com/)
- Basic knowledge of:
  - Go programming
  - Package managers (npm, Yarn, pnpm, Bun)
  - Dependency management concepts

### Fork and Clone

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/hulud-scan.git
cd hulud-scan

# 3. Add upstream remote
git remote add upstream https://github.com/arjunu/hulud-scan.git

# 4. Create a branch for your changes
git checkout -b feature/my-new-feature
```

---

## Development Setup

### Install Dependencies

```bash
# Download Go dependencies
go mod download

# Verify installation
go mod verify
```

### Dependencies Used

```go
require (
    github.com/spf13/cobra v1.10.1       // CLI framework
    github.com/stretchr/testify v1.11.1  // Testing
    gopkg.in/yaml.v3 v3.0.1              // YAML parsing (pnpm)
)
```

---

## Building & Compiling

### Quick Build

```bash
# Build for your current platform
go build -o hulud-scan

# Test the binary
./hulud-scan --version
./hulud-scan scan testdata/npm-project
```

### Build for Specific Platforms

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o hulud-scan-darwin-amd64

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o hulud-scan-darwin-arm64

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o hulud-scan-linux-amd64

# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o hulud-scan-windows-amd64.exe
```

### Build with Optimizations

```bash
# Production build with size optimization
go build -ldflags="-s -w" -o hulud-scan

# With version information
VERSION=$(git describe --tags --always)
go build -ldflags="-X main.Version=$VERSION" -o hulud-scan
```

### Clean Build

```bash
# Remove build artifacts
go clean

# Remove all cached data
go clean -cache -testcache -modcache
```

---

## Running Tests

### Run All Tests

```bash
# Run all tests
go test ./...

# Verbose output
go test ./... -v

# With coverage
go test ./... -cover

# Coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Specific Tests

```bash
# Test a specific package
go test ./internal/parser/...
go test ./internal/graph/...
go test ./internal/scanner/...

# Test a specific function
go test ./internal/parser -run TestParseYarnLock

# Test with race detection
go test ./... -race
```

### Test-Driven Development (TDD) Workflow

hulud-scan was built using TDD. Follow this pattern:

```bash
# 1. Write a failing test (RED)
# Edit: internal/parser/newfeature_test.go

# 2. Run the test - it should fail
go test ./internal/parser -run TestNewFeature

# 3. Write minimal code to pass (GREEN)
# Edit: internal/parser/newfeature.go

# 4. Run the test - it should pass
go test ./internal/parser -run TestNewFeature

# 5. Refactor if needed (REFACTOR)
# Improve code without breaking tests

# 6. Run all tests to ensure nothing broke
go test ./...
```

### Watch Mode for Tests

```bash
# Install nodemon (if using)
npm install -g nodemon

# Watch and re-run tests on file changes
nodemon --exec go test ./... --signal SIGTERM
```

---

## Code Structure

### Project Layout

```
hulud-scan/
├── main.go                 # Entry point
├── cmd/                    # CLI commands
│   ├── root.go            # Root command setup
│   └── scan.go            # Scan command logic
├── internal/              # Internal packages
│   ├── parser/            # Lockfile parsing
│   │   ├── types.go       # Data structures
│   │   ├── parser.go      # npm parser
│   │   ├── yarn.go        # Yarn parser
│   │   ├── pnpm.go        # pnpm parser
│   │   ├── bun.go         # Bun parser
│   │   ├── detector.go    # Auto-detection
│   │   └── *_test.go      # Tests
│   ├── graph/             # Dependency graph
│   │   ├── types.go       # Graph structures
│   │   ├── graph.go       # Graph builder
│   │   └── graph_test.go  # Tests
│   └── scanner/           # Security scanning
│       ├── types.go       # Scanner types
│       ├── scanner.go     # Blocklist matching
│       ├── download.go    # HTTP download
│       ├── cache.go       # Caching logic
│       └── *_test.go      # Tests
├── testdata/              # Test fixtures
│   ├── npm-project/
│   ├── yarn-project/
│   ├── pnpm-project/
│   └── bun-project/
├── go.mod                 # Go module definition
└── go.sum                 # Dependency checksums
```

### Module Responsibilities

#### `cmd/` - Command Line Interface
- Cobra command definitions
- Flag parsing
- User-facing output formatting
- Orchestrates calls to internal packages

#### `internal/parser/` - Lockfile Parsing
- Parse different lockfile formats
- Convert to unified `Lockfile` structure
- Auto-detect lockfile type
- Extract package metadata

#### `internal/graph/` - Dependency Graph
- Build dependency graph from lockfile
- Calculate dependency depth (BFS)
- Mark direct vs transitive dependencies
- Find dependency paths

#### `internal/scanner/` - Security Scanning
- Load blocklists (local/remote)
- Match packages against blocklist
- Cache downloaded blocklists
- Report findings

---

## Adding New Features

### Example 1: Add Support for a New Lockfile Format

Let's say you want to add support for "NewPM" package manager:

**1. Create the parser (`internal/parser/newpm.go`):**

```go
package parser

import (
    "fmt"
    "os"
)

// ParseNewPMLock parses a newpm.lock file
func ParseNewPMLock(lockfilePath string) (*Lockfile, error) {
    data, err := os.ReadFile(lockfilePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read newpm.lock: %w", err)
    }

    // Parse the format and convert to Lockfile structure
    lockfile := &Lockfile{
        Name:               extractProjectNameFromPath(lockfilePath),
        Version:            "unknown",
        LockfileVersion:    1,
        Packages:           make(map[string]*Package),
        DirectDependencies: make(map[string]string),
    }

    // TODO: Add parsing logic here

    return lockfile, nil
}
```

**2. Write tests first (`internal/parser/newpm_test.go`):**

```go
package parser

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestParseNewPMLock(t *testing.T) {
    lockfilePath := "../../testdata/newpm-project/newpm.lock"

    lockfile, err := ParseNewPMLock(lockfilePath)
    require.NoError(t, err)
    require.NotNil(t, lockfile)

    assert.Equal(t, "newpm-test-project", lockfile.Name)
    assert.Len(t, lockfile.Packages, 3)
}
```

**3. Create test fixtures (`testdata/newpm-project/`):**

```bash
mkdir -p testdata/newpm-project
# Add newpm.lock and package.json
```

**4. Update detector (`internal/parser/detector.go`):**

```go
const (
    LockfileTypeNPM   LockfileType = "npm"
    LockfileTypeYarn  LockfileType = "yarn"
    LockfileTypePNPM  LockfileType = "pnpm"
    LockfileTypeBun   LockfileType = "bun"
    LockfileTypeNewPM LockfileType = "newpm"  // Add this
)

func DetectLockfile(projectPath string) (*LockfileInfo, error) {
    lockfiles := []struct {
        filename string
        lockType LockfileType
    }{
        {"package-lock.json", LockfileTypeNPM},
        {"yarn.lock", LockfileTypeYarn},
        {"pnpm-lock.yaml", LockfileTypePNPM},
        {"bun.lockb", LockfileTypeBun},
        {"newpm.lock", LockfileTypeNewPM},  // Add this
    }
    // ...
}

func ParseAuto(projectPath string) (*Lockfile, *LockfileInfo, error) {
    // ...
    switch info.Type {
    case LockfileTypeNPM:
        lockfile, err = ParseLockfile(info.Path)
    // ... other cases
    case LockfileTypeNewPM:  // Add this
        lockfile, err = ParseNewPMLock(info.Path)
    }
    // ...
}
```

**5. Run tests:**

```bash
go test ./internal/parser/... -v
```

**6. Test manually:**

```bash
go build -o hulud-scan
./hulud-scan scan testdata/newpm-project
```

### Example 2: Add a New CLI Flag

**1. Add flag to command (`cmd/scan.go`):**

```go
func init() {
    rootCmd.AddCommand(scanCmd)

    // Existing flags...

    // Add new flag
    scanCmd.Flags().Bool("verbose", false, "Enable verbose output")
}
```

**2. Use flag in command logic:**

```go
func runScan(projectPath string, cmd *cobra.Command) error {
    verbose, _ := cmd.Flags().GetBool("verbose")

    if verbose {
        fmt.Println("Verbose mode enabled")
        // Add more detailed logging
    }
    // ...
}
```

**3. Test manually:**

```bash
go build -o hulud-scan
./hulud-scan scan . --verbose
```

---

## Submitting Changes

### Before Submitting

**1. Run all tests:**
```bash
go test ./...
```

**2. Run linter (if available):**
```bash
golangci-lint run
# Or
go vet ./...
```

**3. Format code:**
```bash
go fmt ./...
```

**4. Update documentation:**
- Update README.md if needed
- Update INSTALLATION.md if adding platform-specific steps
- Add comments to exported functions

### Commit Guidelines

**Commit message format:**
```
<type>: <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding/updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

**Examples:**
```bash
git commit -m "feat: add support for NewPM lockfile format"
git commit -m "fix: correct Yarn scoped package parsing"
git commit -m "docs: update installation guide for Windows"
git commit -m "test: add tests for pnpm parser edge cases"
```

### Pull Request Process

**1. Update your fork:**
```bash
git fetch upstream
git rebase upstream/main
```

**2. Push your branch:**
```bash
git push origin feature/my-new-feature
```

**3. Create Pull Request on GitHub:**
- Go to your fork on GitHub
- Click "New Pull Request"
- Fill in PR template:
  - **Title**: Clear, descriptive title
  - **Description**: What changed and why
  - **Testing**: How you tested the changes
  - **Related Issues**: Link to issues (if any)

**4. Address review feedback:**
```bash
# Make changes based on feedback
git add .
git commit -m "fix: address review comments"
git push origin feature/my-new-feature
```

---

## Style Guide

### Go Code Style

Follow standard Go conventions:

**1. Use `gofmt`:**
```bash
go fmt ./...
```

**2. Naming conventions:**
```go
// Exported functions: PascalCase
func ParseLockfile() {}

// Unexported functions: camelCase
func extractPackageName() {}

// Constants: PascalCase or ALL_CAPS
const MaxRetries = 3
const DEFAULT_TIMEOUT = 30
```

**3. Error handling:**
```go
// Always handle errors
data, err := os.ReadFile(path)
if err != nil {
    return nil, fmt.Errorf("failed to read file: %w", err)
}

// Use %w for error wrapping
```

**4. Comments:**
```go
// Exported functions must have comments
// ParseLockfile parses an npm package-lock.json file
// and returns a Lockfile structure containing all dependencies.
func ParseLockfile(path string) (*Lockfile, error) {
    // ...
}
```

**5. Testing:**
```go
// Use table-driven tests
func TestParsePackageName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "simple package",
            input:    "node_modules/lodash",
            expected: "lodash",
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := parsePackageName(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

---

## Common Development Tasks

### Add a New Test

```bash
# 1. Create test file (if not exists)
touch internal/parser/myfeature_test.go

# 2. Write test
# (See examples above)

# 3. Run test
go test ./internal/parser -run TestMyFeature -v
```

### Debug a Test

```bash
# Add debug output
func TestMyFeature(t *testing.T) {
    t.Logf("Debug: value=%v", myValue)
    // ...
}

# Run with verbose
go test ./internal/parser -run TestMyFeature -v
```

### Profile Performance

```bash
# CPU profiling
go test ./internal/parser -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Memory profiling
go test ./internal/parser -memprofile=mem.prof
go tool pprof mem.prof
```

### Update Dependencies

```bash
# Update all dependencies
go get -u ./...

# Update specific dependency
go get -u github.com/spf13/cobra

# Tidy up
go mod tidy
```

---

## Getting Help

- **Questions**: Open a GitHub Discussion
- **Bugs**: Open a GitHub Issue
- **Chat**: (Add Discord/Slack if available)

---

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help others learn and grow

---

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Credited in commit history

---

**Thank you for contributing to hulud-scan!** 🙏

Every contribution, big or small, helps make the JavaScript ecosystem safer. 🛡️
