# npm Test Project

Clean npm project for testing package-lock.json parsing.

## Dependencies
- lodash@^4.17.21
- express@^4.18.0
- body-parser@1.20.1 (transitive)

## Lockfile
- Format: package-lock.json (npm v3 lockfile format)
- Packages: 3 (2 direct, 1 transitive)

## Test Command
```bash
./hulud-scan scan testdata/npm-project
```

## Expected Result
✅ No compromised packages detected
