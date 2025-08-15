# Release Notes - go-pwr v1.0.4

## 🆕 New Features

### Custom Repository Support

- Users can now configure go-pwr to use their own script repositories
- RocketPowerInc's scriptbin remains the default for new users
- Easy switching between custom and default repositories
- Repository URL validation with support for GitHub, GitLab, and other Git hosting services

### Enhanced Command Line Interface

- **New Flags:**
  - `-h`, `-help` - Show comprehensive help with all available flags and examples
  - `-v`, `-version` - Show detailed version information including git commit and build date
  - `-show-repo` - Display current and default repository URLs
  - `-set-repo <url>` - Set a custom repository URL with validation
  - `-reset-repo` - Reset to default RocketPowerInc scriptbin

### Improved Version Information

- Version output now includes git commit hash for better debugging
- Build date information for tracking release builds
- Go version and platform information
- Repository link for easy access to source code

## 🔧 Technical Improvements

### Build System Enhancements

- Added `build.ps1` PowerShell script for Windows cross-platform building
- Added `build.sh` bash script for Unix-like systems
- Updated Makefile with git commit injection support
- Cross-platform build script (`build-all.ps1`) for release preparation

### Repository Management

- Smart repository path handling (custom repos use separate directories)
- Fresh clone on each startup ensures up-to-date scripts
- Configuration persistence across application restarts
- Backward compatibility with existing installations

## 📚 Documentation Updates

- Added comprehensive Repository Setup Guide (`REPOSITORY_SETUP.md`)
- Enhanced README with repository setup instructions
- Command line flag documentation
- Usage examples and troubleshooting guide

## 🔄 Repository Setup Quick Start

```bash
# View current repository
go-pwr -show-repo

# Set a custom repository
go-pwr -set-repo https://github.com/yourusername/your-scripts.git

# Reset to default repository
go-pwr -reset-repo

# Show version and build info
go-pwr -v

# Show help
go-pwr -h
```

## 📦 Available Downloads

- **Windows**: `go-pwr-windows-amd64.exe`
- **macOS Intel**: `go-pwr-darwin-amd64`
- **macOS Apple Silicon**: `go-pwr-darwin-arm64`
- **Linux AMD64**: `go-pwr-linux-amd64`
- **Linux ARM64**: `go-pwr-linux-arm64`

## 🔧 Installation

1. Download the appropriate binary for your platform
2. Make it executable (Unix-like systems): `chmod +x go-pwr-*`
3. Move to a directory in your PATH or run directly
4. Start with `go-pwr` to launch the interactive TUI

## ⚠️ Breaking Changes

None - This release is fully backward compatible with previous versions.

## 🐛 Bug Fixes

- Improved error handling for repository operations
- Better validation for repository URLs
- Enhanced cross-platform compatibility for build system

---

**Full Changelog**: https://github.com/rocketpowerinc/go-pwr/compare/v1.0.3...v1.0.4
