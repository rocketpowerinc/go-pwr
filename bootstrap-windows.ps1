# Windows Bootstrap Script for go-pwr
# Automates installation of go-pwr and all dependencies on Windows

param(
  [switch]$InstallDependencies = $false,
  [switch]$SkipPathSetup = $false,
  [switch]$Help = $false
)

# Show help
if ($Help) {
  Write-Host "Windows Bootstrap Script for go-pwr" -ForegroundColor Green
  Write-Host ""
  Write-Host "USAGE:" -ForegroundColor Yellow
  Write-Host "  .\bootstrap-windows.ps1 [-InstallDependencies] [-SkipPathSetup] [-Help]"
  Write-Host ""
  Write-Host "PARAMETERS:" -ForegroundColor Yellow
  Write-Host "  -InstallDependencies  Install all dependencies (Go, PowerShell 7+, etc.)"
  Write-Host "  -SkipPathSetup       Skip PATH configuration"
  Write-Host "  -Help                Show this help message"
  Write-Host ""
  Write-Host "EXAMPLES:" -ForegroundColor Yellow
  Write-Host "  .\bootstrap-windows.ps1                        # Install go-pwr only"
  Write-Host "  .\bootstrap-windows.ps1 -InstallDependencies   # Install everything"
  Write-Host ""
  return
}

Write-Host "🚀 Windows Bootstrap Script for go-pwr" -ForegroundColor Green
Write-Host "=======================================" -ForegroundColor Green
Write-Host ""

if ($InstallDependencies) {
  Write-Host "📦 Installing Dependencies..." -ForegroundColor Yellow
  Write-Host ""

  # Install Go
  Write-Host "Installing Go..." -ForegroundColor Cyan
  try {
    winget install -e GoLang.Go --accept-package-agreements --accept-source-agreements
    Write-Host "✓ Go installed successfully" -ForegroundColor Green
  }
  catch {
    Write-Host "⚠️ Failed to install Go: $_" -ForegroundColor Red
  }

  # Install PowerShell 7+
  Write-Host "Installing PowerShell 7+..." -ForegroundColor Cyan
  try {
    winget install Microsoft.PowerShell --accept-package-agreements --accept-source-agreements
    Write-Host "✓ PowerShell 7+ installed successfully" -ForegroundColor Green
  }
  catch {
    Write-Host "⚠️ Failed to install PowerShell 7+: $_" -ForegroundColor Red
  }

  # Install all dependencies
  Write-Host "Installing additional dependencies..." -ForegroundColor Cyan
  try {
    winget install -e Git.Git Microsoft.PowerShell sharkdp.bat GnuWin32.Make charmbracelet.glow charmbracelet.gum GitHub.cli jqlang.jq --accept-package-agreements --accept-source-agreements
    Write-Host "✓ Additional dependencies installed successfully" -ForegroundColor Green
  }
  catch {
    Write-Host "⚠️ Some dependencies may have failed to install: $_" -ForegroundColor Red
  }

  Write-Host ""
  Write-Host "⚠️ Please restart your terminal to ensure Go is in your PATH" -ForegroundColor Yellow
  Write-Host ""
}

# PATH Configuration
if (-not $SkipPathSetup) {
  Write-Host "🔧 Configuring PATH..." -ForegroundColor Yellow

  if (-not (Test-Path $PROFILE)) {
    New-Item -ItemType File -Path $PROFILE -Force | Out-Null
    Write-Host "✓ Created PowerShell profile" -ForegroundColor Green
  }

  if (-not (Select-String -Path $PROFILE -Pattern '\$HOME\\go\\bin' -Quiet)) {
    Add-Content -Path $PROFILE -Value '$env:PATH = "$HOME\go\bin;" + $env:PATH'
    Write-Host "✓ Added Go bin to PATH in PowerShell profile" -ForegroundColor Green

    # Apply to current session
    $env:PATH = "$HOME\go\bin;" + $env:PATH
  }
  else {
    Write-Host "✓ Go bin already in PATH" -ForegroundColor Green
  }
}

# Install go-pwr
Write-Host "📥 Installing go-pwr..." -ForegroundColor Yellow
Write-Host ""

# Try normal installation first
Write-Host "Attempting normal installation..." -ForegroundColor Cyan
try {
  go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest 2>&1 | Out-Null
  if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ go-pwr installed successfully" -ForegroundColor Green
  }
  else {
    throw "Installation failed"
  }
}
catch {
  Write-Host "⚠️ Normal installation failed, trying workaround..." -ForegroundColor Yellow

  # Try with workaround for checksum issues
  Write-Host "Clearing module cache and retrying with direct proxy..." -ForegroundColor Cyan
  try {
    go clean -modcache
    $env:GOPROXY = "direct"
    $env:GOSUMDB = "off"
    go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest
    Write-Host "✓ go-pwr installed successfully with workaround" -ForegroundColor Green
  }
  catch {
    Write-Host "❌ Failed to install go-pwr: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "Manual installation options:" -ForegroundColor Yellow
    Write-Host "1. Download binary from: https://github.com/rocketpowerinc/go-pwr/releases/latest" -ForegroundColor White
    Write-Host "2. Build from source: git clone + make install" -ForegroundColor White
    exit 1
  }
}

# Test installation
Write-Host ""
Write-Host "🔍 Testing installation..." -ForegroundColor Yellow
try {
  $version = & "$env:USERPROFILE\go\bin\go-pwr.exe" --version 2>&1
  if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Installation verified!" -ForegroundColor Green
    Write-Host "$version" -ForegroundColor White
  }
  else {
    throw "Version check failed"
  }
}
catch {
  Write-Host "⚠️ Could not verify installation. Try running: go-pwr --version" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "🎉 Bootstrap complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Restart your terminal (if dependencies were installed)" -ForegroundColor White
Write-Host "2. Run: go-pwr" -ForegroundColor White
Write-Host ""
Write-Host "For help: go-pwr --help" -ForegroundColor Cyan
