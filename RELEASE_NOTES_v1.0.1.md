# Release Notes - go-pwr v1.0.1

## 🆕 New Features

### Theme Persistence

- **Save theme preferences** - Your selected color scheme is now saved and remembered between sessions
- **Persistent configuration** - Theme settings are stored in user config file
- **Automatic theme loading** - Your preferred theme loads automatically on startup

### Linux Enhancements

- **Tmux integration improvements** - Better handling of tmux sessions for Linux users
- **Enhanced terminal compatibility** - Improved experience across different Linux terminals

## 🔧 Improvements

### User Experience

- **Theme memory** - No need to reselect your preferred color scheme every time
- **Configuration management** - Proper config file handling for user preferences
- **Startup experience** - Faster loading with saved preferences

### Technical

- **Configuration system** - Added robust config file management
- **Cross-platform compatibility** - Enhanced support for different operating systems

## 📦 Available Downloads

- **Windows**: `go-pwr-windows-amd64.exe`
- **macOS Intel**: `go-pwr-darwin-amd64`
- **macOS Apple Silicon**: `go-pwr-darwin-arm64`
- **Linux AMD64**: `go-pwr-linux-amd64`

## 🔧 Installation

### Go Install

```bash
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.1
```

### Direct Download

Download the appropriate binary for your platform from the assets above.

## ⚠️ Breaking Changes

None - This release is fully backward compatible with v1.0.0.

---

**Full Changelog**: https://github.com/rocketpowerinc/go-pwr/compare/v1.0.0...v1.0.1
