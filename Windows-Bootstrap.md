# Windows Bootstrap Guide

Complete installation guide for **go-pwr** on Windows.

## Install Flatpak and Snap Package Managers

Visit: Winget should be install by default in windows 11 now
Visit: https://chocolatey.org/install
Visit: https://github.com/ScoopInstaller/Scoop

## 📋 Prerequisites

### Go Installation

```bash
winget install -e GoLang.Go
```

### PowerShell 7+ Installation (Required)

**`go-pwr` requires PowerShell 7+ (pwsh) for running PowerShell scripts.**

```powershell
# Install PowerShell 7+
winget install Microsoft.PowerShell
```

**Note**: Windows PowerShell 5.1 (built-in) is not sufficient. You need PowerShell 7+ (`pwsh.exe`) for full compatibility.

## 🔧 Path Configuration

### Add Go bin to PATH (PowerShell Profile)

```powershell
if (-not (Select-String -Path $PROFILE -Pattern '\$HOME\\go\\bin' -Quiet)) {
    Add-Content -Path $PROFILE -Value '$env:PATH = "$HOME\go\bin;" + $env:PATH'
    . $PROFILE
}
```

## 📋 Dependencies

```powershell
# Install all required dependencies (including PowerShell 7+ if not already installed)
winget install -e Git.Git Microsoft.PowerShell sharkdp.bat GnuWin32.Make charmbracelet.glow charmbracelet.gum GitHub.cli jqlang.jq
```

**Core Dependencies:**

- **Git** - Repository operations
- **PowerShell 7+** - Required for .ps1 script execution
- **Go** - For building from source

**Optional Dependencies:**

- **bat** - Syntax highlighting for script previews
- **Make** - For using Makefile build commands
- **glow/gum** - Enhanced CLI experience
- **GitHub CLI** - For automated releases

## 🚀 Installation Methods

### Method 1: Go Install (Recommended)

```powershell
# Install latest version 
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest

# Or install specific version
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.7
```

**Note**: If you encounter checksum mismatch errors (common with Go module updates), use this workaround:

```powershell
# Temporary workaround for checksum mismatch errors
$env:GOPROXY="direct"; $env:GOSUMDB="off"; go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest

# Alternative: Clear module cache and retry
go clean -modcache
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest
```

### Method 2: Download Binary

1. Visit [Releases](https://github.com/rocketpowerinc/go-pwr/releases/latest)
2. Download `go-pwr-windows-amd64.exe`
3. Place in a directory in your PATH or run directly

### Method 3: Build from Source

```powershell
git clone https://github.com/rocketpowerinc/go-pwr.git
cd go-pwr

# If you encounter checksum mismatch errors, clear the module cache first:
go clean -modcache

# Then run make install
make install
```

**Note**: If you still encounter checksum mismatch errors when building from source, use this workaround:

```powershell
# Clear module cache and use direct proxy
go clean -modcache
$env:GOPROXY="direct"; $env:GOSUMDB="off"; make install
```

## 🚀 Usage

After installation:

```powershell
go-pwr
```

Or with full path:

```powershell
$env:USERPROFILE\go\bin\go-pwr.exe
```

## ⚡ Dev Alias (Advanced Users)

Add this function to your PowerShell profile for quick updates:

```powershell
$gooAlias = @"
# Alias to launch latest go-pwr
function goo {
    Set-Location $env:USERPROFILE
    Remove-Item -Recurse -Force go-pwr -ErrorAction SilentlyContinue
    Remove-Item -Force `"$HOME\\go\\bin\\go-pwr.exe`" -ErrorAction SilentlyContinue
    git clone https://github.com/rocketpowerinc/go-pwr.git
    Set-Location go-pwr
    make install
    & "$env:USERPROFILE\\go\\bin\\go-pwr.exe"
}
"@
$gooAlias | Out-File -Append -Encoding UTF8 $PROFILE
```

Then reload your profile:

```powershell
. $PROFILE
```

## 🐞 Troubleshooting

### Common Issues

- **"go command not found"**: Restart your terminal after installing Go
- **PATH issues**: Use the path configuration steps above
- **Execution Policy**: Run `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser` if needed
- **"pwsh command not found"**: Install PowerShell 7+ with `winget install Microsoft.PowerShell`
- **PowerShell scripts fail**: Ensure you have PowerShell 7+ (`pwsh.exe`), not just Windows PowerShell 5.1

### PowerShell Version Check

Verify you have the correct PowerShell version:

```powershell
# Check PowerShell 7+ is installed
pwsh --version

# Should show version 7.0 or higher
# If command not found, install with: winget install Microsoft.PowerShell
```

### Terminal Recommendations

- **Windows Terminal** (recommended)
- **PowerShell 7+**
- Avoid Command Prompt for best experience
