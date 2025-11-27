# Bun Test Project

Bun project setup for testing bun.lockb parsing.

## Dependencies
- lodash@^4.17.21
- express@^4.18.0

## Lockfile
- Format: bun.lockb (binary format - not included in repo)
- **Note**: bun.lockb is a binary file and requires Bun to generate

## Generate Lockfile
```bash
# Install Bun first: https://bun.sh
curl -fsSL https://bun.sh/install | bash

# Generate lockfile
cd testdata/bun-project
bun install
```

## Test Command
```bash
./hulud-scan scan testdata/bun-project
```

## Expected Result
✅ No compromised packages detected

## How Bun Parsing Works
Since bun.lockb is binary, hulud-scan uses the `bun pm ls` command to read the lockfile. The Bun CLI must be installed and in PATH.

If Bun is not installed, you'll see:
```
Error: bun.lockb detected but 'bun' CLI is not installed or not in PATH
Please install Bun from https://bun.sh or use a different package manager
```
