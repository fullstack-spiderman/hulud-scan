# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure
- Multi-package manager support (npm, Yarn, pnpm, Bun)
- Automatic lockfile detection
- Dependency graph builder with BFS traversal
- Blocklist scanner with local and remote support
- HTTP download with caching (1-hour TTL)
- GitHub URL auto-conversion (blob → raw)
- CLI with scan command
- Version command
- JSON output format
- CI/CD with GitHub Actions
- Docker support
- Multi-platform builds (macOS, Linux, Windows)
- Security scanning (CodeQL, Gosec, govulncheck)
- Comprehensive test suite (40 tests, 80%+ coverage)
- Complete documentation (README, INSTALLATION, CONTRIBUTING, STATUS)

### Changed
- N/A (initial release)

### Deprecated
- N/A (initial release)

### Removed
- N/A (initial release)

### Fixed
- N/A (initial release)

### Security
- N/A (initial release)

## [1.0.0] - YYYY-MM-DD (Template for first release)

### Added
- Initial stable release
- Full support for npm (package-lock.json)
- Full support for Yarn Classic (yarn.lock)
- Full support for pnpm (pnpm-lock.yaml)
- Full support for Bun (bun.lockb)
- Automatic lockfile type detection
- Complete dependency graph analysis
- Direct vs transitive dependency tracking
- Blocklist scanning with severity levels
- Default Wiz Shai-Hulud blocklist (795 packages)
- Custom blocklist support (CSV format)
- Remote blocklist download with HTTPS
- Caching mechanism for downloaded blocklists
- Cross-platform support (macOS, Linux, Windows)
- Multi-architecture binaries (amd64, arm64)
- Docker images (linux/amd64, linux/arm64)
- Exit codes for CI/CD integration
- Configuration flags (format, cache-dir, no-cache)
- Comprehensive documentation
- Test coverage: 80%+
- GitHub Actions workflows (test, build, release, docker, security)

---

## Release Notes Template

Copy this template for future releases:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security fixes
```

---

## Version History

- **Unreleased** - In development
- **1.0.0** - First stable release (pending)

---

## Links

- [Releases](https://github.com/arjunu/hulud-scan/releases)
- [Issues](https://github.com/arjunu/hulud-scan/issues)
- [Pull Requests](https://github.com/arjunu/hulud-scan/pulls)
