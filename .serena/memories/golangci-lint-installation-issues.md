# golangci-lint Installation Issues and Solutions

## Problem Description
```
no required module provides package github.com/golangci/golangci-lint/v2/cmd/golangci-lint; to add it:
        go get github.com/golangci/golangci-lint/v2/cmd/golangci-lint

Exited with code exit status 1

/bin/bash: line 1: ./bin/golangci-lint: No such file or directory

Exited with code exit status 127
```

## Root Causes
1. **Go toolchain version mismatch**: `go install` fails due to version conflicts between Go compiler and tools
2. **Missing module dependency**: The module is not properly added to go.mod
3. **Binary not found**: The build process fails, leaving no executable

## Solutions (in order of preference)

### 1. Use Installation Script (Recommended for CI)
```bash
# Install to GOPATH/bin
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0

# Install to current directory bin/
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ./bin v1.61.0
```

### 2. Use go tool with dedicated module (Recommended for local development)
```bash
# Create dedicated module file
go mod init -modfile=golangci-lint.mod github.com/your-org/your-repo/golangci-lint

# Add golangci-lint as a tool
go get -tool -modfile=golangci-lint.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v1.61.0

# Run golangci-lint as a tool
go tool -modfile=golangci-lint.mod golangci-lint run
```

### 3. Use Homebrew (macOS only)
```bash
brew install golangci-lint
# or for newer versions
brew tap golangci/tap
brew install golangci/tap/golangci-lint
```

### 4. Use Docker (CI environments)
```bash
docker run --rm -t -v $(pwd):/app -w /app \
--user $(id -u):$(id -g) \
-v $(go env GOCACHE):/.cache/go-build -e GOCACHE=/.cache/go-build \
-v $(go env GOMODCACHE):/.cache/mod -e GOMODCACHE=/.cache/mod \
-v ~/.cache/golangci-lint:/.cache/golangci-lint -e GOLANGCI_LINT_CACHE=/.cache/golangci-lint \
golangci/golangci-lint:v1.61.0 golangci-lint run
```

## Updated Makefile Solution
```makefile
install-tools:
	mkdir -p bin
	# Use installation script instead of go build
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ./bin v1.61.0

lint:
	./bin/golangci-lint run

lint-fix:
	./bin/golangci-lint run --fix
```

## CircleCI Configuration Update
```yaml
- run:
    name: Install golangci-lint
    command: |
      curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ./bin v1.61.0
- run:
    name: Run Lint
    command: ./bin/golangci-lint run
```

## Key Points
- **Avoid `go install`**: Prone to toolchain version conflicts
- **Use specific versions**: Always pin to specific versions for reproducibility
- **Prefer installation script**: Most reliable method across different environments
- **Cache binaries**: In CI, cache the binary to speed up builds
- **Use dedicated modules**: For local development with `go tool`

## Troubleshooting
- If script fails: Check network connectivity and GitHub access
- If binary not found: Verify PATH includes the installation directory
- If version conflicts: Use Docker or installation script instead of go install