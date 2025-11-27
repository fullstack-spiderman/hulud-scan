# Test Data Directory

This directory contains sample projects for testing hulud-scan with different package managers.

## Sample Projects

### 1. npm-project/
- **Package Manager**: npm
- **Lockfile**: `package-lock.json`
- **Purpose**: Clean project for testing npm lockfile parsing
- **Dependencies**: lodash@4.17.21, express@4.18.0, body-parser@1.20.1

### 2. yarn-project/
- **Package Manager**: Yarn Classic (v1)
- **Lockfile**: `yarn.lock`
- **Purpose**: Clean project for testing Yarn lockfile parsing
- **Dependencies**: lodash@4.17.21, express@4.18.0, body-parser@1.20.1

### 3. pnpm-project/
- **Package Manager**: pnpm
- **Lockfile**: `pnpm-lock.yaml`
- **Purpose**: Clean project for testing pnpm lockfile parsing
- **Dependencies**: lodash@4.17.21, express@4.18.0, body-parser@1.20.1

### 4. bun-project/
- **Package Manager**: Bun
- **Lockfile**: `bun.lockb` (not included - binary format)
- **Purpose**: Project setup for Bun testing
- **Dependencies**: lodash@4.17.21, express@4.18.0
- **Note**: To generate `bun.lockb`, run `bun install` in this directory with Bun installed

### 5. compromised-project/
- **Package Manager**: npm
- **Lockfile**: `package-lock.json`
- **Purpose**: Project with known vulnerable packages for testing detection
- **Dependencies**:
  - lodash@4.17.20 (CVE-2020-8203 - Prototype pollution)
  - express@4.17.1 (CVE-2022-24999 - DoS vulnerability)
- **Test With**: `--blocklist testdata/sample-blocklist.csv`

### 6. wiz-test-project/
- **Package Manager**: npm
- **Lockfile**: `package-lock.json`
- **Purpose**: Project with Shai-Hulud compromised package for testing Wiz blocklist
- **Dependencies**: 02-echo@0.0.7 (compromised in Shai-Hulud attack)
- **Test With**: Default Wiz blocklist (no flags needed)

## Testing Commands

```bash
# Test npm project
./hulud-scan scan testdata/npm-project --no-cache

# Test yarn project
./hulud-scan scan testdata/yarn-project --no-cache

# Test pnpm project
./hulud-scan scan testdata/pnpm-project --no-cache

# Test bun project (requires Bun installed)
cd testdata/bun-project && bun install && cd ../..
./hulud-scan scan testdata/bun-project --no-cache

# Test compromised packages detection
./hulud-scan scan testdata/compromised-project --blocklist testdata/sample-blocklist.csv

# Test Wiz Shai-Hulud blocklist
./hulud-scan scan testdata/wiz-test-project --no-cache
```

## Blocklist Files

### sample-blocklist.csv
Custom blocklist with known CVEs for testing:
- lodash@4.17.20 (CVE-2020-8203)
- lodash@4.17.19 (CVE-2019-10744)
- express@4.17.1 (CVE-2022-24999)
- minimist@1.2.5 (CVE-2021-44906)

Format: `package_name,version,severity,reason,cve`

## Adding New Test Projects

To add a new test project:

1. Create a directory: `testdata/my-test-project/`
2. Add `package.json` with dependencies
3. Generate lockfile using the appropriate package manager
4. Update this README with project details
5. Add test case in `internal/parser/*_test.go`
