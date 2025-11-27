# pnpm Test Project

Clean pnpm project for testing pnpm-lock.yaml parsing.

## Dependencies
- lodash@^4.17.21
- express@^4.18.0
- body-parser@1.20.1 (transitive)

## Lockfile
- Format: pnpm-lock.yaml (lockfileVersion 6.0)
- Packages: 3 (2 direct, 1 transitive)

## Test Command
```bash
./hulud-scan scan testdata/pnpm-project
```

## Expected Result
✅ No compromised packages detected

## Regenerate Lockfile
```bash
cd testdata/pnpm-project
pnpm install
```
