# Release Notes - go-pwr v1.0.0

## 🎉 Initial Release

This is the first public release of go-pwr, a cross-platform script launcher built with Go and Bubble Tea.

## ✨ Features

### Core Functionality

- **Interactive TUI** powered by Charm's Bubble Tea framework
- **Cross-platform support** for Windows, macOS, and Linux
- **Script browsing** with directory navigation
- **Script execution** in new terminal windows
- **Syntax highlighting** for script previews (when `bat` is installed)

### User Interface

- **Tabbed interface** with Scripts, Options, and About tabs
- **Responsive design** that adapts to terminal size
- **Keyboard navigation** with intuitive shortcuts
- **Script preview pane** with syntax highlighting

### Script Management

- **Automatic repository cloning** to `$HOME/Downloads/Temp/scriptbin`
- **Directory-based organization** for scripts
- **Support for multiple script types** (bash, PowerShell, etc.)

## 📦 Available Downloads

- **Windows**: `go-pwr-windows-amd64.exe`
- **macOS Intel**: `go-pwr-darwin-amd64`
- **macOS Apple Silicon**: `go-pwr-darwin-arm64`
- **Linux AMD64**: `go-pwr-linux-amd64`

## 🔧 Installation

### Go Install

```bash
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.0
```

### Direct Download

Download the appropriate binary for your platform from the assets above.

## 🚀 Getting Started

1. Install go-pwr using one of the methods above
2. Run `go-pwr` to start the interactive interface
3. Browse and execute scripts from the RocketPowerInc scriptbin

---

**Initial Release**: First public version of go-pwr
