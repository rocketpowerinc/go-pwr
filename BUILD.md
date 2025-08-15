# Build Documentation

This document explains how to build go-pwr from source and create releases using the provided build scripts.

## 📋 Prerequisites

### Required Tools

- **Go 1.19+** - [Download](https://golang.org/dl/)
- **Git** - [Download](https://git-scm.com/downloads)
- **GitHub CLI (gh)** - Required for automated releases
  - Windows: `winget install GitHub.cli`
  - macOS: `brew install gh`
  - Linux: `sudo apt install gh` or `sudo dnf install gh`

### Authentication

Make sure you're authenticated with GitHub CLI:

```bash
gh auth login
```

## 🔧 Build Scripts Overview

### Available Scripts

| Script            | Purpose                                 | Platform    |
| ----------------- | --------------------------------------- | ----------- |
| `build.ps1`       | Single platform build (Windows-focused) | Windows     |
| `build.sh`        | Single platform build (Unix-focused)    | Linux/macOS |
| `build-all.ps1`   | Multi-platform build + release          | Windows     |
| `new-release.ps1` | Generate release notes template         | Windows     |

## 🚀 Build Methods

### Method 1: Quick Development Build

For testing and development:

```powershell
# Windows
.\build.ps1

# Linux/macOS
chmod +x build.sh
./build.sh
```

This creates a single binary for your current platform with git commit info.

### Method 2: Multi-Platform Build

Build for all supported platforms:

```powershell
# Windows (recommended for releases)
.\build-all.ps1

# Output: build/ directory with 5 binaries
```

### Method 3: Standard Go Build

Traditional Go build without extras:

```bash
# Simple build
go build -o go-pwr cmd/go-pwr/main.go

# With version info (manual)
go build -ldflags "-X main.gitCommit=$(git rev-parse --short HEAD) -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o go-pwr cmd/go-pwr/main.go
```

## 📦 Release Process

### Complete Release Workflow

#### Step 1: Prepare New Release

```powershell
# Generate release notes template
.\new-release.ps1 -Version "1.0.5"
```

This will:

- Create `RELEASE_NOTES_v1.0.5.md` with a template
- Auto-detect the previous version for changelog links
- Open the file in notepad for editing

#### Step 2: Update Version

Edit `cmd/go-pwr/main.go` and update:

```go
const version = "1.0.5"  // Update this line
```

#### Step 3: Edit Release Notes

Fill out the generated `RELEASE_NOTES_v1.0.5.md` with:

- New features
- Bug fixes
- Breaking changes
- Any other relevant information

#### Step 4: Build and Release

```powershell
# Build for all platforms and create GitHub release
.\build-all.ps1 -Version "1.0.5" -Release
```

This will:

- Build binaries for all platforms
- Create git tag `v1.0.5`
- Push the tag to GitHub
- Create GitHub release with all binaries
- Test the release with `go install`

### Alternative: Build-Only Release

If you prefer to create the GitHub release manually:

```powershell
# Just build the binaries
.\build-all.ps1 -Version "1.0.5"

# Then manually create release at:
# https://github.com/rocketpowerinc/go-pwr/releases/new
```

## 🔍 Script Details

### `build-all.ps1` Options

```powershell
# Show help
.\build-all.ps1 -Help

# Build only (auto-detect version from main.go)
.\build-all.ps1

# Build with specific version
.\build-all.ps1 -Version "1.0.5"

# Build and create GitHub release
.\build-all.ps1 -Release

# Build specific version and release
.\build-all.ps1 -Version "1.0.5" -Release
```

### `new-release.ps1` Options

```powershell
# Generate template (auto-detect previous version)
.\new-release.ps1 -Version "1.0.5"

# Generate template with specific previous version
.\new-release.ps1 -Version "1.0.5" -PreviousVersion "1.0.3"
```

## 📁 Build Outputs

### Directory Structure After Build

```
build/
├── go-pwr-windows-amd64.exe    # Windows 64-bit
├── go-pwr-darwin-amd64         # macOS Intel
├── go-pwr-darwin-arm64         # macOS Apple Silicon
├── go-pwr-linux-amd64          # Linux 64-bit
└── go-pwr-linux-arm64          # Linux ARM64
```

### Binary Information

Each binary includes:

- Version number from `main.go`
- Git commit hash (short)
- Build timestamp
- Go version used
- Target platform/architecture

Verify with:

```bash
./go-pwr-linux-amd64 -v
```

## 🐛 Troubleshooting

### Common Issues

#### "gh command not found"

```bash
# Install GitHub CLI
winget install GitHub.cli  # Windows
brew install gh            # macOS
sudo apt install gh        # Ubuntu/Debian
```

#### "git tag already exists"

```bash
# Delete local and remote tag
git tag -d v1.0.5
git push origin :refs/tags/v1.0.5
```

#### "release notes file not found"

Make sure you've created the release notes file:

```powershell
.\new-release.ps1 -Version "1.0.5"
```

#### Build fails with "GOOS/GOARCH not set"

The script should handle this automatically, but if it fails:

```powershell
# Reset environment
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
```

### Version Detection Issues

If the script can't detect the version from `main.go`:

```powershell
# Explicitly specify version
.\build-all.ps1 -Version "1.0.5"
```

Make sure the version line in `cmd/go-pwr/main.go` follows this exact format:

```go
const version = "1.0.5"
```

## 🔄 CI/CD Integration

### GitHub Actions (Future)

The build scripts can be integrated into GitHub Actions:

```yaml
# Example workflow
- name: Build All Platforms
  run: .\build-all.ps1 -Version ${{ github.ref_name }}

- name: Create Release
  run: .\build-all.ps1 -Version ${{ github.ref_name }} -Release
```

## 📚 Best Practices

### Version Management

1. **Semantic Versioning**: Use `MAJOR.MINOR.PATCH` format
2. **Update main.go**: Always update the version constant
3. **Git Tags**: Let the script create tags automatically
4. **Release Notes**: Always fill out comprehensive release notes

### Testing

1. **Test Local Build**: Run `.\build.ps1` first to test locally
2. **Test All Platforms**: Use `.\build-all.ps1` to ensure all platforms build
3. **Test Installation**: Script automatically tests `go install` after release

### Security

1. **Review Changes**: Always review what's being released
2. **Test Binaries**: Test at least the primary platform binary
3. **Release Notes**: Include any security-related changes

## 📖 Examples

### Complete Release Example

```powershell
# 1. Start new release
.\new-release.ps1 -Version "1.0.5"

# 2. Edit cmd/go-pwr/main.go (update version)
# 3. Edit RELEASE_NOTES_v1.0.5.md (add changes)

# 4. Build and release
.\build-all.ps1 -Version "1.0.5" -Release

# 5. Verify release
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.5
go-pwr -v
```

### Development Build Example

```powershell
# Quick development build
.\build.ps1

# Test the build
.\go-pwr.exe -v
```

This documentation should cover everything needed to build and release go-pwr efficiently!
