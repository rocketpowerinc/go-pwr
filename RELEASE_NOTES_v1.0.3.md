# Release Notes - go-pwr v1.0.3

## 🆕 New Features

### Linux Tmux Integration

- **Tmux warning system** - Prominent warning for Linux users about tmux benefits
- **Automatic tmux detection** - Smart detection of existing tmux sessions
- **Tmux launch capability** - Option to automatically start in tmux on Linux
- **Session persistence** - Better handling of script execution in tmux environments

### Enhanced Linux Experience

- **Improved terminal compatibility** - Better support for various Linux terminal emulators
- **Background operation support** - Enhanced script execution for long-running tasks
- **SSH-friendly operation** - Better experience for remote Linux sessions

## 🔧 Improvements

### User Experience

- **Clear tmux guidance** - Users understand the benefits of using tmux on Linux
- **Flexible tmux usage** - Can disable warnings if preferred
- **Better session management** - Improved handling of terminal sessions

### Documentation

- **Enhanced documentation** - Improved setup and usage instructions
- **Platform-specific guidance** - Better instructions for different operating systems
- **Troubleshooting guides** - Added common issue solutions

## 📚 Documentation Updates

### Bootstrap Guides

- **Comprehensive installation guides** - Detailed setup for Windows, macOS, and Linux
- **Platform-specific instructions** - Tailored guidance for each operating system
- **Dependency management** - Clear instructions for required tools

### Usage Documentation

- **Keyboard shortcuts** - Complete reference for all available shortcuts
- **Feature explanations** - Detailed explanations of all features
- **Best practices** - Recommended usage patterns

## 📦 Available Downloads

- **Windows**: `go-pwr-windows-amd64.exe`
- **macOS Intel**: `go-pwr-darwin-amd64`
- **macOS Apple Silicon**: `go-pwr-darwin-arm64`
- **Linux AMD64**: `go-pwr-linux-amd64`

## 🔧 Installation

### Go Install

```bash
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.3
```

### Direct Download

Download the appropriate binary for your platform from the assets above.

## ⚠️ Breaking Changes

None - This release is fully backward compatible with previous versions.

## 🐧 Linux Users

This release significantly improves the Linux experience. For the best experience on Linux:

```bash
# Start with tmux (recommended)
tmux new-session go-pwr

# Or disable warnings if preferred
export GO_PWR_NO_TMUX_WARNING=1
go-pwr
```

---

**Full Changelog**: https://github.com/rocketpowerinc/go-pwr/compare/v1.0.1...v1.0.3
