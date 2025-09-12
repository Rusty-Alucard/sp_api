# golangci-lint v2 Upgrade Information

## Current Status
- **Previous version**: v1.61.0 (built with go1.23.1)
- **Current version**: v2.4.0 (built with go1.25.0)
- **Status**: Successfully upgraded and working

## Key Changes in v2
1. **Configuration version**: Must use `version: 2` in `.golangci.yml`
2. **Linter changes**: Some linters renamed or deprecated
3. **Performance improvements**: Better caching and parallel execution
4. **Go version compatibility**: Built with Go 1.25.0

## Installation Commands
```bash
# Install v2.4.0 (latest)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ./bin v2.4.0

# Verify installation
./bin/golangci-lint --version
```

## Configuration Updates Required
- `.golangci.yml` already has `version: 2` ✓
- All linters are compatible with v2 ✓
- No breaking changes detected ✓

## Benefits of v2
- Better performance
- Improved caching
- More accurate linting
- Better Go 1.25+ support
- Enhanced error reporting

## Troubleshooting
- If v1 was installed via goenv, it may conflict
- Use `./bin/golangci-lint` to ensure using local v2
- Check `which golangci-lint` to verify PATH
- Use installation script instead of `go install` for reliability