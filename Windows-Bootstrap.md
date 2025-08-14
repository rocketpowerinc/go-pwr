# Windows Bootstrap Guide

Complete installation guide for **go-pwr** on Windows.

## 📋 Prerequisites

### Go Installation

```bash
winget install -e GoLang.Go
```

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
# Install all required dependencies
winget install -e Git.Git sharkdp.bat GnuWin32.Make charmbracelet.glow charmbracelet.gum GitHub.cli jqlang.jq
```


## 🚀 Installation Methods

### Method 1: Go Install (Recommended)

```powershell
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
make install
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

### Terminal Recommendations

- **Windows Terminal** (recommended)
- **PowerShell 7+**
- Avoid Command Prompt for best experience
