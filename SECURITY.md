# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of hulud-scan seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Please DO NOT:

- Open a public GitHub issue for security vulnerabilities
- Disclose the vulnerability publicly before it has been addressed

### Please DO:

**Report security vulnerabilities via GitHub Security Advisories:**

1. Go to the [Security tab](https://github.com/arjunu/hulud-scan/security)
2. Click "Report a vulnerability"
3. Fill out the form with details

**Or send an email to:** support@symbiotesai.com

### What to Include:

Please include the following information:
- Type of vulnerability (e.g., RCE, injection, authentication bypass)
- Full paths of affected source files
- Location of the affected source code (tag/branch/commit)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the vulnerability
- Any suggested fixes (if you have them)

## Response Timeline

- **Initial Response**: Within 48 hours of receiving your report
- **Assessment**: Within 7 days of initial response
- **Fix Development**: Varies based on severity
- **Public Disclosure**: After a fix is available and deployed

## Vulnerability Disclosure Process

1. **Report received** - We'll acknowledge receipt within 48 hours
2. **Assessment** - We'll assess the vulnerability and its impact
3. **Fix development** - We'll develop a fix for the vulnerability
4. **Testing** - We'll test the fix to ensure it resolves the issue
5. **Release** - We'll release a new version with the fix
6. **Advisory** - We'll publish a security advisory
7. **Credit** - We'll credit you (if desired) in the advisory

## Security Best Practices for Users

### When Using hulud-scan:

1. **Always use the latest version**

   ```bash
   # Check your version
   hulud-scan --version

   # Update to latest
   go install github.com/arjunu/hulud-scan@latest
   ```

2. **Verify checksums when downloading binaries**

   ```bash
   # Download and verify
   curl -L https://github.com/arjunu/hulud-scan/releases/download/v1.0.0/hulud-scan-linux-amd64 -o hulud-scan
   curl -L https://github.com/arjunu/hulud-scan/releases/download/v1.0.0/hulud-scan-linux-amd64.sha256 -o hulud-scan.sha256
   sha256sum -c hulud-scan.sha256
   ```

3. **Run with minimal permissions**
   - Don't run as root/admin unless necessary
   - Use in sandboxed environments when scanning untrusted projects

4. **Be cautious with custom blocklists**
   - Only use blocklists from trusted sources
   - Verify blocklist URLs are HTTPS
   - Review blocklist contents before using

5. **Keep Go updated**
   - hulud-scan is written in Go
   - Keep your Go runtime updated for security patches

### Secure Configuration

```bash
# Use HTTPS for blocklists
hulud-scan scan . --blocklist https://trusted-source.com/blocklist.csv

# Use caching to reduce exposure to network attacks
hulud-scan scan . --cache-dir ~/.hulud-scan/cache

# Specify custom cache directory in a secure location
hulud-scan scan . --cache-dir /secure/path/cache
```

## Known Security Considerations

### By Design

1. **Network Access** - hulud-scan downloads blocklists over HTTPS
   - Mitigation: Uses HTTPS by default, validates certificates

2. **File System Access** - Reads project lockfiles and cache
   - Mitigation: Only reads files, doesn't write to project directories

3. **External Dependencies** - Relies on Go packages
   - Mitigation: Minimal dependencies, regular updates via Dependabot

4. **Bun CLI Execution** - Executes `bun pm ls` for Bun lockfiles
   - Mitigation: Only when bun.lockb detected, command hardcoded

### Not Covered

- **Zero-day vulnerabilities** in packages (not in blocklists yet)
- **Malicious code execution** from scanned packages (we only read, don't install)
- **Network interception** (use VPN/secure network when downloading blocklists)

## Security Features

✅ **What hulud-scan does to stay secure:**

- Minimal dependencies (only 3 external packages)
- HTTPS-only for remote blocklist downloads
- SHA256 checksums for release binaries
- No code execution from scanned packages
- Read-only access to project files
- Security scanning via CodeQL in CI
- Regular dependency updates via Dependabot
- Vulnerability scanning with govulncheck

## Security Scanning Results

Our CI/CD pipeline runs:

- **CodeQL** - Static analysis for code vulnerabilities
- **Gosec** - Go security checker
- **govulncheck** - Go vulnerability database checker
- **Dependabot** - Automated dependency updates

View our [Security Advisories](https://github.com/arjunu/hulud-scan/security/advisories) for any past issues.

## Bug Bounty Program

We currently do not have a bug bounty program. However, we greatly appreciate security researchers who report vulnerabilities responsibly.

## Questions?

If you have questions about this security policy, please open an issue or contact the maintainers.

---

**Last updated:** 2025-11-28
