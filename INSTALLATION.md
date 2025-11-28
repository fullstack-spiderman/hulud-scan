# Installation & Usage Guide

Complete guide for installing and using hulud-scan on macOS, Windows, and Linux.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
  - [macOS](#macos)
  - [Linux](#linux)
  - [Windows](#windows)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### All Platforms

You need **Go 1.21+** installed to build hulud-scan.

#### Check if Go is installed

```bash
go version
```

If you see `go version go1.21.x` or higher, you're good to go! ✅

#### Install Go if needed

**macOS:**

```bash
# Using Homebrew
brew install go

# Or download from <https://go.dev/dl/>
```

**Linux (Ubuntu/Debian):**

```bash
# Using apt
sudo apt update
sudo apt install golang-go

# Or download from <https://go.dev/dl/>
```

**Windows:**

- Download installer from: <https://go.dev/dl/>
- Run the `.msi` installer
- Restart your terminal after installation

---

## Installation

### Package Managers

#### macOS - Homebrew

```bash
# Add the tap
brew tap fullstack-spiderman/tap

# Install hulud-scan
brew install hulud-scan

# Verify installation
hulud-scan --version
```

**Update:**

```bash
brew upgrade hulud-scan
```

**Uninstall:**

```bash
brew uninstall hulud-scan
brew untap fullstack-spiderman/tap
```

#### Windows - Scoop

```powershell
# Add the bucket
scoop bucket add hulud-scan https://github.com/fullstack-spiderman/scoop-bucket

# Install hulud-scan
scoop install hulud-scan

# Verify installation
hulud-scan --version
```

**Update:**

```powershell
scoop update hulud-scan
```

**Uninstall:**

```powershell
scoop uninstall hulud-scan
scoop bucket rm hulud-scan
```

---

### macOS

#### Option 1: Prebuilt Binary (Recommended - No Go Required!)

**Latest Release:** [v1.0.2](https://github.com/fullstack-spiderman/hulud-scan/releases/latest)

```bash
# Intel Macs (x86_64)
curl -LO https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_darwin_amd64.tar.gz
tar -xzf hulud-scan_1.0.2_darwin_amd64.tar.gz
sudo mv hulud-scan /usr/local/bin/

# Apple Silicon (M1/M2/M3/M4 - ARM64)
curl -LO https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_darwin_arm64.tar.gz
tar -xzf hulud-scan_1.0.2_darwin_arm64.tar.gz
sudo mv hulud-scan /usr/local/bin/

# Verify installation
hulud-scan --version
```

#### Option 2: Build from Source

**Prerequisites:** [Go 1.21+](https://go.dev/dl/)

```bash
# 1. Clone the repository
git clone https://github.com/fullstack-spiderman/hulud-scan.git
cd hulud-scan

# 2. Build the binary
go build -o hulud-scan

# 3. (Optional) Move to PATH for global access
sudo mv hulud-scan /usr/local/bin/

# 4. Verify installation
hulud-scan --version
```

#### Option 3: Direct Build Without Clone

**Prerequisites:** [Go 1.21+](https://go.dev/dl/)

```bash
# Install directly with Go
go install github.com/fullstack-spiderman/hulud-scan@latest

# The binary will be in ~/go/bin/hulud-scan
# Add ~/go/bin to your PATH in ~/.zshrc or ~/.bash_profile:
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

#### Using the Tool on macOS

```bash
# Scan current directory
./hulud-scan scan

# Scan specific project
./hulud-scan scan /path/to/your/project

# If installed globally:
hulud-scan scan /path/to/your/project
```

---

### Linux

#### Option 1: Prebuilt Binary (No Go Required!)

**Latest Release:** [v1.0.2](https://github.com/fullstack-spiderman/hulud-scan/releases/latest)

```bash
# x86_64 (Intel/AMD 64-bit)
curl -LO https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_linux_amd64.tar.gz
tar -xzf hulud-scan_1.0.2_linux_amd64.tar.gz
sudo mv hulud-scan /usr/local/bin/

# ARM64 (Raspberry Pi, ARM servers)
curl -LO https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_linux_arm64.tar.gz
tar -xzf hulud-scan_1.0.2_linux_arm64.tar.gz
sudo mv hulud-scan /usr/local/bin/

# Verify installation
hulud-scan --version
```

#### Option 2: Build from Source (Ubuntu/Debian/Fedora/Arch)

**Prerequisites:** [Go 1.21+](https://go.dev/dl/)

```bash
# 1. Install Git (if not installed)
sudo apt install git  # Ubuntu/Debian
# sudo dnf install git  # Fedora
# sudo pacman -S git    # Arch

# 2. Clone the repository
git clone https://github.com/fullstack-spiderman/hulud-scan.git
cd hulud-scan

# 3. Build the binary
go build -o hulud-scan

# 4. Make it executable
chmod +x hulud-scan

# 5. (Optional) Move to PATH for global access
sudo mv hulud-scan /usr/local/bin/

# 6. Verify installation
hulud-scan --version
```

#### Option 3: Direct Install with Go

**Prerequisites:** [Go 1.21+](https://go.dev/dl/)

```bash
# Install directly
go install github.com/fullstack-spiderman/hulud-scan@latest

# Add Go bin to PATH (add to ~/.bashrc or ~/.zshrc)
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

#### Using the Tool on Linux

```bash
# Scan current directory
./hulud-scan scan

# Scan specific project
./hulud-scan scan /home/user/my-project

# If installed globally:
hulud-scan scan /home/user/my-project
```

---

### Windows

#### Option 1: Prebuilt Binary (Easiest Option!)

**Latest Release:** [v1.0.2](https://github.com/fullstack-spiderman/hulud-scan/releases/latest)

**Using Browser:**

1. Visit: [GitHub Releases](https://github.com/fullstack-spiderman/hulud-scan/releases/latest)
2. Download:
   - `hulud-scan_1.0.2_windows_amd64.zip` (for x64 systems)
   - `hulud-scan_1.0.2_windows_arm64.zip` (for ARM64 systems)
3. Extract the `.zip` file
4. Add to your PATH or run from the extracted directory

**Using PowerShell (x64):**

```powershell
# Download and extract
Invoke-WebRequest -Uri "https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_windows_amd64.zip" -OutFile "hulud-scan.zip"
Expand-Archive -Path hulud-scan.zip -DestinationPath .

# Verify installation
.\hulud-scan.exe --version
```

**Using PowerShell (ARM64):**

```powershell
# Download and extract
Invoke-WebRequest -Uri "https://github.com/fullstack-spiderman/hulud-scan/releases/latest/download/hulud-scan_1.0.2_windows_arm64.zip" -OutFile "hulud-scan.zip"
Expand-Archive -Path hulud-scan.zip -DestinationPath .

# Verify installation
.\hulud-scan.exe --version
```

#### Option 2: Build from Source (PowerShell/CMD)

**Prerequisites:** [Go 1.21+](https://go.dev/dl/)

```powershell
# 1. Install Git (if not installed)
# Download from: https://git-scm.com/download/win

# 2. Clone the repository
git clone https://github.com/fullstack-spiderman/hulud-scan.git
cd hulud-scan

# 3. Build the binary
go build -o hulud-scan.exe

# 4. The executable is now ready: hulud-scan.exe
```

#### Option 3: Install with Go

**Prerequisites:** [Go 1.21+](https://go.dev/dl/)

```powershell
# Install directly
go install github.com/fullstack-spiderman/hulud-scan@latest

# The binary will be in: %USERPROFILE%\go\bin\hulud-scan.exe
# Add to PATH if needed (System Properties > Environment Variables)
```

#### Using the Tool on Windows

**PowerShell:**

```powershell
# Scan current directory
.\hulud-scan.exe scan

# Scan specific project
.\hulud-scan.exe scan C:\Users\YourName\my-project

# If added to PATH:
hulud-scan scan C:\Users\YourName\my-project
```

**Command Prompt (cmd):**

```cmd
# Scan current directory
hulud-scan.exe scan

# Scan specific project
hulud-scan.exe scan C:\Users\YourName\my-project
```

**Git Bash (recommended for Windows):**

```bash
# Scan current directory
./hulud-scan.exe scan

# Scan specific project
./hulud-scan.exe scan /c/Users/YourName/my-project
```

---

## Quick Start

### Basic Usage (All Platforms)

```bash
# 1. Navigate to your JavaScript/TypeScript project
cd /path/to/your/project

# 2. Run the scan
hulud-scan scan .

# Or specify the project path:
hulud-scan scan /path/to/your/project
```

### What Gets Scanned?

hulud-scan automatically detects and scans:

- ✅ `package-lock.json` (npm)
- ✅ `yarn.lock` (Yarn)
- ✅ `pnpm-lock.yaml` (pnpm)
- ✅ `bun.lockb` (Bun)

---

## Usage Examples

### Example 1: Scan a Project with Default Settings

**macOS/Linux:**

```bash
./hulud-scan scan ~/projects/my-app
```

**Windows (PowerShell):**

```powershell
.\hulud-scan.exe scan C:\projects\my-app
```

**Output:**

```text
🔍 Scanning project at: ~/projects/my-app
🔎 Detecting lockfile in: ~/projects/my-app
📄 Detected: npm (package-lock.json)
✅ Found 1247 packages
Project: my-app@1.0.0

📊 Building dependency graph...
📋 Loading blocklist from: <https://github.com/wiz-sec-public/...>
✅ Loaded 795 blocklist entries

🔍 Scanning for compromised packages...

============================================================
SCAN RESULTS
============================================================

Total packages scanned: 1247
Issues found: 0

✅ No compromised packages detected!
```

### Example 2: Use Custom Blocklist

**macOS/Linux:**

```bash
./hulud-scan scan . --blocklist ./my-blocklist.csv
```

**Windows:**

```powershell
.\hulud-scan.exe scan . --blocklist .\my-blocklist.csv
```

### Example 3: Disable Caching (Always Fresh Download)

**macOS/Linux:**

```bash
./hulud-scan scan . --no-cache
```

**Windows:**

```powershell
.\hulud-scan.exe scan . --no-cache
```

### Example 4: Scan Multiple Projects

**macOS/Linux (Bash):**

```bash
#!/bin/bash
for project in ~/projects/*; do
    echo "Scanning $project"
    ./hulud-scan scan "$project"
done
```

**Windows (PowerShell):**

```powershell
Get-ChildItem C:\projects | ForEach-Object {
    Write-Host "Scanning $($_.FullName)"
    .\hulud-scan.exe scan $_.FullName
}
```

### Example 5: Use in CI/CD

**GitHub Actions (All Platforms):**

```yaml
name: Security Scan
on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install hulud-scan
        run: go install github.com/fullstack-spiderman/hulud-scan@latest

      - name: Scan dependencies
        run: hulud-scan scan .
```

**GitLab CI (All Platforms):**

```yaml
security_scan:
  image: golang:1.21
  script:
    - go install github.com/fullstack-spiderman/hulud-scan@latest
    - hulud-scan scan .
```

---

## Platform-Specific Notes

### macOS Notes

- **Binary Location**: After building, `hulud-scan` is in the current directory
- **Global Install**: Move to `/usr/local/bin/` for system-wide access
- **Permissions**: No special permissions needed (unless moving to `/usr/local/bin/`)
- **Shell**: Works in Terminal, iTerm2, and any bash/zsh shell

### Linux Notes

- **Binary Location**: After building, `hulud-scan` is in the current directory
- **Global Install**: Move to `/usr/local/bin/` or `~/.local/bin/`
- **Permissions**: May need `sudo` to move to `/usr/local/bin/`
- **Shell**: Works in any shell (bash, zsh, fish, etc.)

### Windows Notes

- **Binary Name**: Use `hulud-scan.exe` (not just `hulud-scan`)
- **Path Prefix**: Use `.\` for current directory in PowerShell/CMD
- **Path Format**:
  - PowerShell/CMD: `C:\Users\Name\project`
  - Git Bash: `/c/Users/Name/project`
- **Recommended Shell**: Git Bash for Unix-like experience
- **Admin Rights**: Not required for scanning

---

## Troubleshooting

### Issue: "command not found" or "not recognized"

**macOS/Linux:**

```bash
# Make sure the binary is executable
chmod +x hulud-scan

# Run with ./ prefix if not in PATH
./hulud-scan scan

# Or add to PATH
export PATH="$PATH:$(pwd)"
```

**Windows:**

```powershell
# Use full filename with .exe
.\hulud-scan.exe scan

# Or add current directory to PATH temporarily
$env:PATH += ";$(Get-Location)"
```

### Issue: "go: command not found"

Go is not installed or not in PATH.

**macOS:**

```bash
brew install go
# Or download from <https://go.dev/dl/>
```

**Linux:**

```bash
sudo apt install golang-go  # Ubuntu/Debian
sudo dnf install golang     # Fedora
```

**Windows:**

- Download and install from: <https://go.dev/dl/>
- Restart terminal after installation

### Issue: "no supported lockfile found"

Your project doesn't have a lockfile. Generate one:

**npm:**

```bash
npm install
```

**Yarn:**

```bash
yarn install
```

**pnpm:**

```bash
pnpm install
```

**Bun:**

```bash
bun install
```

### Issue: Bun lockfile not detected

For Bun, you need the Bun CLI installed:

**macOS/Linux:**

```bash
curl -fsSL <https://bun.sh/install> | bash
```

**Windows:**

```powershell
powershell -c "irm bun.sh/install.ps1 | iex"
```

### Issue: Permission denied (Linux/macOS)

```bash
# Make binary executable
chmod +x hulud-scan

# If moving to /usr/local/bin
sudo mv hulud-scan /usr/local/bin/
```

### Issue: Slow download on first run

The first run downloads the blocklist (795 entries). Subsequent runs use cache.

To force fresh download:

```bash
./hulud-scan scan . --no-cache
```

---

## Advanced Configuration

### Custom Cache Directory

**macOS/Linux:**

```bash
./hulud-scan scan . --cache-dir ~/.my-custom-cache
```

**Windows:**

```powershell
.\hulud-scan.exe scan . --cache-dir C:\Users\YourName\.cache\hulud
```

### All Available Flags

```bash
hulud-scan scan [path] [flags]

Flags:
  -f, --format string       Output format (table or json) (default "table")
  -c, --config string       Path to config file
      --blocklist string    Blocklist URL or local file path
      --cache-dir string    Cache directory for downloaded blocklists
      --no-cache           Disable caching (always download fresh)
  -h, --help               Help for scan
```

### Example: JSON Output

**macOS/Linux:**

```bash
./hulud-scan scan . --format json > results.json
```

**Windows:**

```powershell
.\hulud-scan.exe scan . --format json > results.json
```

---

## Uninstallation

### macOS/Linux

```bash
# If installed globally
sudo rm /usr/local/bin/hulud-scan

# Remove cache
rm -rf ~/.hulud-scan
```

### Windows Uninstall

```powershell
# Delete the binary
Remove-Item .\hulud-scan.exe

# Remove cache
Remove-Item -Recurse -Force $env:USERPROFILE\.hulud-scan
```

---

## Getting Help

```bash
# View all commands
hulud-scan --help

# View scan command help
hulud-scan scan --help

# Check version
hulud-scan --version
```

---

## Support

- **Issues**: <https://github.com/fullstack-spiderman/hulud-scan/issues>
- **Documentation**: See README.md and STATUS.md
- **Examples**: See testdata/ directory

---
