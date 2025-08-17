# Release Notes - go-pwr v1.0.8

## 🆕 New Features

### Bootstrap Scripts and Documentation
- Added platform-specific bootstrap scripts for automated installation
- New comprehensive bootstrap documentation for Windows, macOS, and Linux
- Streamlined setup process for new users across all platforms

### Enhanced Repository Management
- Improved repository selection and management interface
- Better handling of custom repository configurations
- Enhanced user experience for repository switching

## 🔧 Improvements

### UI/UX Enhancements
- Fixed ASCII art centering on large terminals for better visual presentation
- Improved terminal rendering and display consistency
- Enhanced overall user interface responsiveness

### Platform Compatibility
- Fixed Windows installation issues and improved compatibility
- Better handling of different terminal sizes and environments
- Improved cross-platform consistency

### Documentation
- Extensive documentation updates and improvements
- Better installation guides for all platforms
- Enhanced README with clearer instructions and examples

## 🐛 Bug Fixes

### Terminal Display
- Fixed ASCII art rendering issues on various terminal sizes
- Improved text centering and alignment
- Resolved display inconsistencies across different environments

### Installation Fixes
- Resolved Windows installation and compatibility issues
- Fixed linting issues in codebase
- Improved build and installation reliability

## 📚 Documentation Updates

- Added comprehensive bootstrap documentation for Windows, macOS, and Linux
- Updated installation guides with platform-specific instructions
- Enhanced README with better examples and usage instructions
- Improved build documentation and development guides

## 🔧 Technical Details

- **Git Commit**: 6f92ea3
- **Build Date**: 2025-08-17
- **Go Version**: go version go1.25.0 windows/amd64

## 📦 Available Downloads

- **Windows**: `go-pwr-windows-amd64.exe`
- **macOS Intel**: `go-pwr-darwin-amd64`
- **macOS Apple Silicon**: `go-pwr-darwin-arm64`
- **Linux AMD64**: `go-pwr-linux-amd64`
- **Linux ARM64**: `go-pwr-linux-arm64`

## 🔧 Installation

### Go Install (Recommended)
```bash
# Install latest version
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest

# Or install specific version
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.8
```

### Direct Download
1. Download the appropriate binary for your platform from the assets above
2. Make it executable (Unix-like systems): `chmod +x go-pwr-*`
3. Move to a directory in your PATH or run directly

## ⚠️ Breaking Changes

None - This release is fully backward compatible with previous versions.

## 🔄 Upgrade Instructions

No special upgrade steps required. Simply install the new version.

---

**Full Changelog**: https://github.com/rocketpowerinc/go-pwr/compare/v1.0.7...v1.0.8
