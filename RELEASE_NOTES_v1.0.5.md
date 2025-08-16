# Release Notes - go-pwr v1.0.5

# Release Notes - go-pwr v1.0.5

## 🐛 Bug Fixes

### Windows Makefile Compatibility

- **Fixed Windows date command compatibility**: Updated Makefile to use PowerShell-based date formatting instead of Unix `date` command
- **Fixed shell redirection syntax**: Changed from Unix-style `2>/dev/null` to Windows-style `2>nul`
- **Improved git commit detection**: Enhanced git command compatibility for Windows make environment
- **Resolved "system cannot find path" errors**: Windows users can now successfully run `make install` with proper version information

### Build System Improvements

- **PowerShell compatibility**: Updated to use `powershell -NoProfile` for reliable cross-platform execution
- **Version injection**: Restored proper git commit hash and build date injection on Windows builds
- **Shell command reliability**: Improved shell command execution in Windows make environment

## 📚 Documentation Updates

- Cleaned up Windows Bootstrap guide by removing unnecessary troubleshooting sections
- Updated Makefile comments to reflect Windows compatibility improvements

## 🔧 Technical Details

- **Git Commit**: 59176c0
- **Build Date**: 2025-08-16
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
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.5
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

**Full Changelog**: https://github.com/rocketpowerinc/go-pwr/compare/v1.0.4...v1.0.5
