# GitHub Actions Workflows

This document describes all GitHub Actions workflows configured for hulud-scan.

## 📋 Workflows Overview

| Workflow | Trigger | Purpose | Artifacts |
|----------|---------|---------|-----------|
| **test.yml** | PR, Push to main/develop | Run tests on multiple platforms | Coverage reports |
| **build.yml** | Push to main | Build multi-platform binaries | Binary artifacts |
| **release.yml** | Git tags (v*.*.*) | Create releases with binaries | GitHub Release + Assets |
| **docker.yml** | Push to main, Tags | Build and publish Docker images | Docker images |
| **codeql.yml** | PR, Push, Weekly | Security scanning | Security alerts |

---

## 🧪 test.yml - Continuous Integration

**Triggers:**
- Push to `main` or `develop`
- Pull requests to `main` or `develop`

**Jobs:**

### 1. Test Matrix
- **Platforms**: Ubuntu, macOS, Windows
- **Go Versions**: 1.21, 1.22, 1.23
- **Actions**:
  - ✅ Verify dependencies
  - ✅ Build project
  - ✅ Run tests with race detection
  - ✅ Upload coverage to Codecov (Ubuntu + Go 1.23 only)

### 2. Lint
- **Platform**: Ubuntu
- **Actions**:
  - ✅ Run golangci-lint

### 3. Security
- **Platform**: Ubuntu
- **Actions**:
  - ✅ Run Gosec security scanner
  - ✅ Run govulncheck for vulnerabilities

**Example:**
```yaml
# Runs on every PR
on:
  pull_request:
    branches: [ main, develop ]
```

---

## 🔨 build.yml - Build Binaries

**Triggers:**
- Push to `main` branch
- Manual workflow dispatch

**Build Matrix:**
- macOS Intel (darwin-amd64)
- macOS Apple Silicon (darwin-arm64)
- Linux AMD64 (linux-amd64)
- Linux ARM64 (linux-arm64)
- Windows AMD64 (windows-amd64.exe)

**Jobs:**

### 1. Build Binaries
- Cross-compile for all platforms
- Include version info via ldflags
- Upload as artifacts (30-day retention)

### 2. Test Binaries
- Download and test binaries on native platforms
- Run version, help, and scan commands
- Validate functionality before release

**Version Information:**
```bash
go build -ldflags="-s -w \
  -X main.Version=$(git describe --tags --always --dirty)"
```

---

## 🚀 release.yml - Release Workflow

**Triggers:**
- Git tags matching `v*.*.*` (e.g., v1.0.0)
- Manual workflow dispatch with version input

**Jobs:**

### 1. Create Release
- Extract version from tag
- Generate changelog from commits
- Create GitHub Release with installation instructions

### 2. Build and Upload Assets
- Build for all platforms (5 binaries)
- Generate SHA256 checksums
- Upload binaries and checksums to release

### 3. Publish Docker
- Build multi-arch Docker images (amd64, arm64)
- Push to GitHub Container Registry
- Tag with version, major.minor, major, and latest

**Release Assets:**
```
hulud-scan-darwin-amd64
hulud-scan-darwin-amd64.sha256
hulud-scan-darwin-arm64
hulud-scan-darwin-arm64.sha256
hulud-scan-linux-amd64
hulud-scan-linux-amd64.sha256
hulud-scan-linux-arm64
hulud-scan-linux-arm64.sha256
hulud-scan-windows-amd64.exe
hulud-scan-windows-amd64.exe.sha256
```

**Create a Release:**
```bash
# Create and push a tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Workflow runs automatically
```

---

## 🐳 docker.yml - Docker Images

**Triggers:**
- Push to `main` branch
- Git tags (v*.*.*)
- Pull requests (build only, no push)

**Features:**
- Multi-architecture support (amd64, arm64)
- Layer caching for faster builds
- Published to GitHub Container Registry (ghcr.io)

**Image Tags:**
- `latest` - Latest main branch
- `v1.0.0` - Specific version
- `v1.0` - Major.minor version
- `v1` - Major version
- `main-<sha>` - Main branch commits
- `pr-<number>` - Pull requests

**Pull the Image:**
```bash
docker pull ghcr.io/arjunu/hulud-scan:latest
```

**Run:**
```bash
docker run --rm -v $(pwd):/app ghcr.io/arjunu/hulud-scan:latest scan /app
```

---

## 🔒 codeql.yml - Security Scanning

**Triggers:**
- Push to `main` or `develop`
- Pull requests to `main`
- Weekly on Sunday (scheduled)

**Actions:**
- Initialize CodeQL for Go
- Analyze code for security vulnerabilities
- Report to GitHub Security tab

**View Results:**
- Navigate to: Repository → Security → Code scanning alerts

---

## 📊 Workflow Status Badges

Add these to your README.md:

```markdown
![Test](https://github.com/fullstack-spiderman/hulud-scan/actions/workflows/test.yml/badge.svg)
![Build](https://github.com/fullstack-spiderman/hulud-scan/actions/workflows/build.yml/badge.svg)
![Release](https://github.com/fullstack-spiderman/hulud-scan/actions/workflows/release.yml/badge.svg)
![Docker](https://github.com/fullstack-spiderman/hulud-scan/actions/workflows/docker.yml/badge.svg)
![CodeQL](https://github.com/fullstack-spiderman/hulud-scan/actions/workflows/codeql.yml/badge.svg)
```

---

## 🔑 Required Secrets

### For Docker Publishing (Optional)
- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub password/token

**Note:** GitHub Container Registry uses `GITHUB_TOKEN` automatically (no setup needed).

### For Codecov (Optional)
- `CODECOV_TOKEN` - Codecov upload token

---

## 🚀 How to Create a Release

### Option 1: Git Tag (Recommended)
```bash
# 1. Make sure you're on main and up to date
git checkout main
git pull

# 2. Create an annotated tag
git tag -a v1.0.0 -m "Release v1.0.0

- Feature 1
- Feature 2
- Bug fix 3"

# 3. Push the tag
git push origin v1.0.0

# 4. Workflow runs automatically!
# Check: https://github.com/fullstack-spiderman/hulud-scan/actions
```

### Option 2: Manual Workflow Dispatch
```bash
# 1. Go to GitHub Actions
# 2. Select "Release" workflow
# 3. Click "Run workflow"
# 4. Enter version (e.g., v1.0.0)
# 5. Click "Run workflow" button
```

---

## 📦 Artifacts Retention

| Workflow | Artifact | Retention |
|----------|----------|-----------|
| test.yml | Coverage reports | 7 days |
| build.yml | Binaries | 30 days |
| release.yml | Release assets | Permanent |
| docker.yml | Docker images | Per registry policy |

---

## 🔧 Workflow Configuration

### Enable/Disable Workflows

Edit workflow files in `.github/workflows/` and commit changes.

### Customize Build Matrix

Edit `build.yml` or `release.yml`:

```yaml
strategy:
  matrix:
    include:
      - goos: darwin
        goarch: amd64
      # Add more platforms here
```

### Change Triggers

Edit the `on:` section:

```yaml
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  # Add more triggers
```

---

## 📈 Monitoring Workflows

### View Workflow Runs
1. Go to repository on GitHub
2. Click "Actions" tab
3. Select workflow from left sidebar
4. View runs, logs, and artifacts

### View Build Status
- Check status badges in README
- Review pull request checks
- Monitor email notifications

### Download Artifacts
1. Go to workflow run
2. Scroll to "Artifacts" section
3. Click to download

---

## 🐛 Troubleshooting

### Workflow Failed?

**Common Issues:**

1. **Tests failing**
   - Check test logs in Actions tab
   - Run `go test ./...` locally

2. **Build failing**
   - Verify Go version compatibility
   - Check for missing dependencies

3. **Release failing**
   - Ensure tag follows `v*.*.*` format
   - Check if release already exists

4. **Docker build failing**
   - Verify Dockerfile syntax
   - Check for missing files in .dockerignore

### Re-run Failed Workflows
1. Go to failed workflow run
2. Click "Re-run failed jobs"
3. Monitor progress

---

## 🎯 Best Practices

1. **Always test locally first**
   ```bash
   go test ./...
   go build
   ```

2. **Use semantic versioning**
   - v1.0.0 (major.minor.patch)
   - v1.0.0-beta.1 (pre-release)

3. **Write good commit messages**
   - They become your changelog!

4. **Review workflow logs**
   - Catch issues early

5. **Keep workflows updated**
   - Use latest action versions
   - Update Go versions

---

## 📝 Workflow Maintenance

### Update Actions
```yaml
# From:
uses: actions/checkout@v3

# To:
uses: actions/checkout@v4
```

### Update Go Versions
```yaml
strategy:
  matrix:
    go-version: ['1.21', '1.22', '1.23']  # Add new versions
```

### Add New Platforms
```yaml
matrix:
  include:
    - goos: freebsd  # Add FreeBSD support
      goarch: amd64
```

---

## 🔗 Useful Links

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Action](https://github.com/actions/setup-go)
- [Docker Buildx Action](https://github.com/docker/build-push-action)
- [CodeQL Action](https://github.com/github/codeql-action)

---

*Last updated: 2025-11-28*
