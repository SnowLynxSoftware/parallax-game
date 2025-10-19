# GitHub Actions

## PR Test Validation Workflow

**File:** `.github/workflows/pr-test-validation.yml`

### Overview

This workflow automatically validates Pull Requests targeting the `develop` branch by running build checks, code quality validations, and unit tests.

### Trigger

- **Event:** Pull Request
- **Target Branch:** `develop`
- **PR Types:** opened, synchronize, reopened

### Jobs

#### 1. Build & Quality Checks

Runs first to ensure code quality and successful compilation.

**Steps:**

1. Checkout code
2. Setup Go 1.25
3. Download dependencies (`make deps`)
4. Check code formatting (`make fmt`)
   - Fails if code is not properly formatted
5. Run go vet (`make vet`)
   - Static analysis to catch common issues
6. Build application (`make build`)
   - Compiles for Linux/amd64

**Duration:** ~2-3 minutes

#### 2. Run Unit Tests

Runs only if Build & Quality Checks succeed.

**Steps:**

1. Checkout code
2. Setup Go 1.25
3. Download dependencies (`make deps`)
4. Run all tests (`make test`)
   - Executes full test suite with verbose output

**Duration:** ~2-3 seconds

### Total Workflow Time

**~2-4 minutes** (depending on cache hits and test duration)

### Performance Optimizations

- Go modules are cached between runs
- Sequential execution prevents unnecessary test runs if build fails
- Dependencies are shared via cache

### Usage

This workflow runs automatically on PRs to `develop`. To ensure your PR passes:

```bash
# Before creating a PR, run these locally:
make deps      # Install dependencies
make fmt       # Format code
make vet       # Run static analysis
make build     # Verify build succeeds
make test      # Run all tests
```

Or use the comprehensive check:

```bash
make check     # Runs fmt, vet, and test
```

### Troubleshooting

**If "Check code formatting" fails:**

```bash
make fmt
git add .
git commit -m "Fix code formatting"
git push
```

**If "Run go vet" fails:**
Fix the issues reported by vet, then:

```bash
make vet       # Verify fixes
git add .
git commit -m "Fix vet issues"
git push
```

**If "Build application" fails:**
Check the error logs in the GitHub Actions output and fix compilation errors.

**If "Run Unit Tests" fails:**

```bash
make test      # Run tests locally to see failures
# Fix failing tests
git add .
git commit -m "Fix failing tests"
git push
```
