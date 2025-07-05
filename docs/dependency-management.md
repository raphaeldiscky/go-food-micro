# Safe Dependency Management in Go Microservices

This guide provides comprehensive strategies for safely updating dependencies in your Go microservices project.

## Overview

Your project consists of:

- **Shared Package**: `internal/pkg` - Common utilities and shared code
- **Services**:
  - `catalogwriteservice` - Product catalog write operations
  - `catalogreadservice` - Product catalog read operations
  - `orderservice` - Order management
- **Go Version**: 1.23.6 (latest)

## Available Tools

### 1. Basic Update Script

```bash
# Update all services
task update-dependencies

# Update specific service
./scripts/update-dependencies.sh pkg
./scripts/update-dependencies.sh catalogwriteservice
```

### 2. Enhanced Safe Update Script (Recommended)

```bash
# Safe update with backup, testing, and rollback capability
task safe-update-dependencies

# Check for outdated dependencies
task check-dependencies

# Direct script usage
./scripts/safe-update-dependencies.sh update
./scripts/safe-update-dependencies.sh check
./scripts/safe-update-dependencies.sh rollback
```

## Safe Update Process

### Step 1: Pre-Update Assessment

```bash
# Check current Go version
go version

# Check for outdated dependencies
task check-dependencies

# Review current dependency state
cd internal/pkg && go list -m all
```

### Step 2: Create Backup

The safe update script automatically creates backups:

- `go.mod` and `go.sum` files are backed up
- Timestamped backup directory: `backup-YYYYMMDD-HHMMSS/`

### Step 3: Update Dependencies

```bash
# Recommended: Use safe update
task safe-update-dependencies

# Alternative: Manual update
cd internal/pkg
go get -u -t -d -v ./...
go mod tidy
go mod verify
```

### Step 4: Verification

The safe update script automatically:

- ✅ Builds each service
- ✅ Runs tests
- ✅ Verifies module integrity
- ✅ Provides rollback capability if issues occur

## Manual Update Commands

### Check Outdated Dependencies

```bash
# Check all modules for updates
go list -u -m all

# Check specific module
go list -u -m github.com/gin-gonic/gin

# See what would be updated
go get -u -d ./...
```

### Update Specific Dependencies

```bash
# Update to latest version
go get -u github.com/gin-gonic/gin

# Update to specific version
go get github.com/gin-gonic/gin@v1.9.1

# Update to latest patch version
go get -u=patch github.com/gin-gonic/gin
```

### Module Maintenance

```bash
# Clean up unused dependencies
go mod tidy

# Verify module integrity
go mod verify

# Download dependencies
go mod download

# Vendor dependencies (if using vendoring)
go mod vendor
```

## Best Practices

### 1. **Incremental Updates**

- Update one service at a time
- Test thoroughly before moving to the next
- Use semantic versioning constraints in `go.mod`

### 2. **Version Pinning**

```go
// In go.mod - pin to specific versions for stability
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/go-sql-driver/mysql v1.7.1
)
```

### 3. **Dependency Analysis**

```bash
# Analyze dependency graph
go mod graph

# Check for security vulnerabilities
go list -json -deps ./... | jq 'select(.Vulnerabilities)'

# Use govulncheck (if available)
govulncheck ./...
```

### 4. **Testing Strategy**

```bash
# Run all tests
make unit-test
make integration-test
make e2e-test

# Test specific service
cd internal/services/catalogwriteservice
go test ./... -v
```

## Troubleshooting

### Common Issues

#### 1. **Version Conflicts**

```bash
# Check for version conflicts
go mod why -m github.com/conflicting/module

# Resolve conflicts
go mod tidy
go mod download
```

#### 2. **Build Failures**

```bash
# Clean module cache
go clean -modcache

# Rebuild
go mod download
go mod tidy
go build ./...
```

#### 3. **Test Failures**

- Check for breaking changes in major version updates
- Review changelogs for updated dependencies
- Consider using `go get -u=patch` for safer updates

### Rollback Process

```bash
# Automatic rollback (if using safe update script)
./scripts/safe-update-dependencies.sh rollback

# Manual rollback
cp backup-YYYYMMDD-HHMMSS/service/go.mod internal/services/service/
cp backup-YYYYMMDD-HHMMSS/service/go.sum internal/services/service/
cd internal/services/service
go mod tidy
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Dependency Update
on:
  schedule:
    - cron: "0 2 * * 1" # Weekly on Monday at 2 AM

jobs:
  update-deps:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: "1.23"

      - name: Check for updates
        run: make check-dependencies

      - name: Safe update
        run: make safe-update-dependencies

      - name: Run tests
        run: make unit-test

      - name: Create PR
        if: success()
        run: |
          git config user.name "Dependency Bot"
          git config user.email "bot@example.com"
          git add .
          git commit -m "chore: update dependencies"
          git push origin HEAD:dependency-update
```

## Monitoring and Alerts

### Dependency Health Checks

```bash
# Check for deprecated packages
go list -f '{{if .Deprecated}}{{.Path}}{{end}}' all

# Check for security issues
go list -json -deps ./... | jq 'select(.Vulnerabilities)'

# Monitor dependency size
go list -f '{{.Path}} {{.Size}}' all
```

### Automated Monitoring

- Set up weekly dependency update checks
- Monitor for security vulnerabilities
- Track dependency update frequency and success rate

## Advanced Techniques

### 1. **Dependency Graph Analysis**

```bash
# Generate dependency graph
go mod graph | dot -Tpng -o deps.png

# Analyze import paths
go list -f '{{.ImportPath}} -> {{join .Imports " "}}' ./...
```

### 2. **Custom Update Strategies**

```bash
# Update only direct dependencies
go list -f '{{if not .Indirect}}{{.Path}}{{end}}' -m all | xargs -I {} go get -u {}

# Update with specific constraints
go get -u -t -d -v ./... && go mod edit -go=1.23
```

### 3. **Vendor Management**

```bash
# Create vendor directory
go mod vendor

# Use vendored dependencies
go build -mod=vendor ./...

# Update vendor
go mod vendor
```

## Conclusion

The safe update script provides a comprehensive approach to dependency management with:

- ✅ Automatic backups
- ✅ Incremental updates
- ✅ Built-in testing
- ✅ Rollback capability
- ✅ Detailed logging

Use `make safe-update-dependencies` for the most reliable update process, and always test thoroughly in a staging environment before deploying to production.
