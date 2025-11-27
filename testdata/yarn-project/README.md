# Yarn Test Project

Clean Yarn Classic (v1) project for testing yarn.lock parsing.

## Dependencies
- lodash@^4.17.21
- express@^4.18.0
- body-parser@1.20.1 (transitive)

## Lockfile
- Format: yarn.lock (Yarn Classic v1 format)
- Packages: 3 (2 direct, 1 transitive)

## Test Command
```bash
./hulud-scan scan testdata/yarn-project
```

## Expected Result
✅ No compromised packages detected

## Regenerate Lockfile
```bash
cd testdata/yarn-project
yarn install
```
