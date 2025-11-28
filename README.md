# hulud-scan 🔍

> Supply-chain security scanner for JavaScript/TypeScript projects

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE) ![Tests](https://img.shields.io/badge/tests-40%20passing-brightgreen) ![Coverage](https://img.shields.io/badge/coverage-80%25-green)

**hulud-scan** is an open-source CLI tool that detects compromised packages and supply-chain attacks in JavaScript/TypeScript projects by scanning lockfiles and comparing against known blocklists.

Named after the Shai-Hulud attacks (2024-2025), this tool helps protect your projects from malicious dependencies.

---

## 🚀 Features

- ✅ **Multi-Package Manager Support**
  - npm (`package-lock.json`)
  - Yarn Classic (`yarn.lock`)
  - pnpm (`pnpm-lock.yaml`)
  - Bun (`bun.lockb`)

- ✅ **Automatic Lockfile Detection**
  - No configuration needed
  - Auto-detects and uses the right parser

- ✅ **Comprehensive Dependency Analysis**
  - Scans direct AND transitive dependencies
  - Shows full dependency paths
  - Identifies direct vs indirect dependencies

- ✅ **Multiple Blocklist Sources**
  - Default: Wiz Shai-Hulud blocklist (795 packages)
  - Custom blocklists supported (CSV format)
  - Remote URLs with caching (1-hour TTL)

- ✅ **CI/CD Ready**
  - Exit codes for automation
  - JSON output format
  - Fast scanning (< 10s for typical projects)

- ✅ **Cross-Platform**
  - macOS, Linux, Windows
  - No dependencies required

---

## 📦 Installation

### Option 1: Download Prebuilt Binary (Recommended - No Go Required!)

**Latest Release:** [v1.0.2](https://github.com/fullstack-spiderman/hulud-scan/releases/latest)

**macOS:**

```bash
# Intel Macs
curl -LO https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_darwin_amd64.tar.gz
tar -xzf hulud-scan_1.0.2_darwin_amd64.tar.gz
sudo mv hulud-scan /usr/local/bin/

# Apple Silicon (M1/M2/M3/M4)
curl -LO https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_darwin_arm64.tar.gz
tar -xzf hulud-scan_1.0.2_darwin_arm64.tar.gz
sudo mv hulud-scan /usr/local/bin/
```

**Linux:**

```bash
# x86_64
curl -LO https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_linux_amd64.tar.gz
tar -xzf hulud-scan_1.0.2_linux_amd64.tar.gz
sudo mv hulud-scan /usr/local/bin/

# ARM64
curl -LO https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_linux_arm64.tar.gz
tar -xzf hulud-scan_1.0.2_linux_arm64.tar.gz
sudo mv hulud-scan /usr/local/bin/
```

**Windows:**

```powershell
# Download the latest release
# Visit: https://github.com/fullstack-spiderman/hulud-scan/releases/latest

# For Windows x64:
# 1. Download hulud-scan_1.0.2_windows_amd64.zip
# 2. Extract the archive
# 3. Add to PATH or run from extracted directory

# Or using PowerShell (x64):
Invoke-WebRequest -Uri "https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_windows_amd64.zip" -OutFile "hulud-scan.zip"
Expand-Archive -Path hulud-scan.zip -DestinationPath .
```

**Verify Installation:**

```bash
hulud-scan --version
```

### Option 2: Install with Go

**Prerequisites:** [Go 1.21+](https://go.dev/dl/)

```bash
# Install with Go
go install github.com/fullstack-spiderman/hulud-scan@latest

# Or build from source
git clone https://github.com/fullstack-spiderman/hulud-scan.git
cd hulud-scan
go build -o hulud-scan
```

### Option 3: Package Managers (Coming Soon)

🚧 Homebrew (macOS) and Scoop (Windows) support planned for future releases.

---

📖 **See [INSTALLATION.md](INSTALLATION.md) for detailed platform-specific instructions** (macOS, Linux, Windows)

---

## 🎯 Quick Start

```bash
# Scan current directory
hulud-scan scan .

# Scan specific project
hulud-scan scan /path/to/your/project

# Use custom blocklist
hulud-scan scan . --blocklist ./my-blocklist.csv

# Disable caching (always download fresh)
hulud-scan scan . --no-cache

# JSON output
hulud-scan scan . --format json
```

### Example Output

```text
🔍 Scanning project at: ./my-app
🔎 Detecting lockfile in: ./my-app
📄 Detected: npm (package-lock.json)
✅ Found 1247 packages
Project: my-app@1.0.0

📊 Building dependency graph...
📋 Loading blocklist from: https://github.com/wiz-sec-public/...
✅ Loaded 795 blocklist entries

🔍 Scanning for compromised packages...

============================================================
SCAN RESULTS
============================================================

Total packages scanned: 1247
Issues found: 1

⚠️  SECURITY ISSUES DETECTED:

1. malicious-package@1.0.0 [CRITICAL]
   Type: transitive dependency
   Path: my-app → express → body-parser → malicious-package
   Reason: Compromised package (Shai-Hulud attack)

❌ Critical security issues detected!
```

---

## 📚 Documentation

- **[INSTALLATION.md](INSTALLATION.md)** - Platform-specific installation guide (macOS, Linux, Windows)
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Developer guide for contributors
- **[STATUS.md](STATUS.md)** - Project status, roadmap, and completed features
- **[testdata/README.md](testdata/README.md)** - Test data documentation

---

## 🛠️ Usage

### Basic Commands

```bash
# Scan with default settings
hulud-scan scan .

# Scan specific directory
hulud-scan scan /path/to/project

# View help
hulud-scan --help
hulud-scan scan --help
```

### Advanced Options

```bash
# Custom blocklist
hulud-scan scan . --blocklist https://example.com/blocklist.csv
hulud-scan scan . --blocklist ./local-blocklist.csv

# Custom cache directory
hulud-scan scan . --cache-dir ~/.my-cache

# Disable caching
hulud-scan scan . --no-cache

# JSON output (for CI/CD)
hulud-scan scan . --format json > results.json
```

### Exit Codes

- `0` - No issues found ✅
- `1` - Critical issues detected or scan error ❌

Perfect for CI/CD pipelines!

---

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/parser/...
go test ./internal/scanner/...
go test ./internal/graph/...

# Verbose output
go test ./... -v
```

### Test Coverage

- **internal/graph**: 95.6%
- **internal/parser**: 72.9%
- **internal/scanner**: 82.3%

### 40 automated tests, all passing ✅

---

## 🏗️ Architecture

```text
hulud-scan/
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Root command
│   └── scan.go            # Scan command
├── internal/
│   ├── parser/            # Lockfile parsers
│   │   ├── parser.go      # npm (package-lock.json)
│   │   ├── yarn.go        # Yarn (yarn.lock)
│   │   ├── pnpm.go        # pnpm (pnpm-lock.yaml)
│   │   ├── bun.go         # Bun (bun.lockb)
│   │   └── detector.go    # Auto-detection
│   ├── graph/             # Dependency graph
│   │   └── graph.go       # Graph builder & traversal
│   └── scanner/           # Security scanner
│       ├── scanner.go     # Blocklist matching
│       ├── download.go    # Remote blocklist fetch
│       └── cache.go       # Caching layer
├── testdata/              # Test fixtures
└── main.go                # Entry point
```

---

## 🤝 Contributing

We welcome contributions! Whether it's:

- 🐛 Bug reports
- 💡 Feature requests
- 📖 Documentation improvements
- 🔧 Code contributions

📖 **See [CONTRIBUTING.md](CONTRIBUTING.md) for:**

- Development setup
- How to compile and test
- Code structure
- PR guidelines

---

## 🗺️ Roadmap

### ✅ Phase 1 - CLI (Complete!)

- [x] Multi-package manager support (npm, Yarn, pnpm, Bun)
- [x] Auto-detection
- [x] Blocklist scanning
- [x] Dependency graph analysis
- [x] CI/CD integration

### 🔄 Phase 2 - Enhanced CLI (Planned)

- [ ] Lifecycle script detection
- [ ] Config file support
- [ ] Whitelist/ignore mechanism
- [ ] Multiple output formats (HTML, SARIF)
- [ ] Yarn Berry (v2+) support

### 🔮 Phase 3 - TUI (Future)

- [ ] Interactive terminal UI
- [ ] Visual dependency tree
- [ ] Real-time filtering

### 🌐 Phase 4 - Web Dashboard (Future)

- [ ] Multi-repo monitoring
- [ ] Historical tracking
- [ ] Alerts & notifications
- [ ] Team collaboration

---

## 📊 Project Stats

- **Language**: Go
- **Lines of Code**: ~2,500 (excluding tests)
- **Test Lines**: ~1,500
- **Dependencies**: 3 (cobra, testify, yaml.v3)
- **Test Coverage**: 80%+
- **Tests**: 40 automated tests
- **Supported Formats**: 4 lockfile types

---

## 🔒 Security

### What hulud-scan Detects

✅ **Known compromised packages** in blocklists
✅ **Direct and transitive dependencies**
✅ **Full dependency chain** for each issue

### Limitations

⚠️ **Zero-day attacks** - Not yet in blocklists
⚠️ **Obfuscated malware** - Advanced hiding techniques
⚠️ **Lifecycle scripts** - Detection planned for Phase 2

**hulud-scan is a defense layer, not a silver bullet.** Use alongside other security practices.

---

## 📄 Blocklist Format

hulud-scan supports two CSV formats:

### Wiz Format (Simple)

```csv
Package,Version
malicious-pkg,=1.0.0
bad-package,=2.1.0
```

### Full Format (Detailed)

```csv
package_name,version,severity,reason,cve
lodash,4.17.20,critical,Prototype pollution,CVE-2020-8203
express,4.17.1,high,DoS vulnerability,CVE-2022-24999
```

**Severity levels**: `critical`, `high`, `medium`, `low`, `info`

---

## 🌟 Use Cases

### For Developers

```bash
# Scan before committing
hulud-scan scan .
git commit -m "feat: new feature"
```

### For CI/CD

```yaml
# GitHub Actions
- name: Scan dependencies
  run: hulud-scan scan .
```

### For Organizations

```bash
# Scan all projects
for dir in ~/projects/*; do
  hulud-scan scan "$dir"
done
```

### For Security Teams

```bash
# Custom enterprise blocklist
hulud-scan scan . --blocklist https://internal.corp/blocklist.csv
```

---

## 🙏 Credits

- Inspired by the [Wiz Security Shai-Hulud research](https://github.com/wiz-sec-public/wiz-research-iocs)
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Tested with [Testify](https://github.com/stretchr/testify)

---

## 📜 License

MIT License - see [LICENSE](LICENSE) file for details

---

## 🔗 Links

- **GitHub**: <https://github.com/fullstack-spiderman/hulud-scan>
- **Issues**: <https://github.com/fullstack-spiderman/hulud-scan/issues>
- **Wiz Blocklist**: <https://github.com/wiz-sec-public/wiz-research-iocs>

---

## 💬 Support

- 📖 Documentation in this repo
- 🐛 Report bugs via GitHub Issues
- 💡 Feature requests welcome
- 🤝 PRs encouraged (see [CONTRIBUTING.md](CONTRIBUTING.md))

---

**Stay safe! Scan your dependencies regularly.** 🛡️

---

*Built with ❤️ by [@arjun]([https://github.com/arjunu](https://github.com/fullstack-spiderman))*
