# PowerShell script to build go-pwr for all platforms

param(
  [string]$Version = "1.0.4"
)

# Get git commit hash (short)
$gitCommit = git rev-parse --short HEAD
if ($LASTEXITCODE -ne 0) {
  $gitCommit = "unknown"
}

# Get current date/time in ISO format
$buildDate = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"

# Build flags to inject version information
$ldflags = "-X main.gitCommit=$gitCommit -X main.buildDate=$buildDate"

# Create build directory
if (!(Test-Path "build")) {
  New-Item -ItemType Directory -Path "build" | Out-Null
}

Write-Host "Building go-pwr v$Version for all platforms..." -ForegroundColor Green
Write-Host "Git commit: $gitCommit" -ForegroundColor Yellow
Write-Host "Build date: $buildDate" -ForegroundColor Yellow
Write-Host ""

# Build for Windows
Write-Host "Building for Windows..." -ForegroundColor Cyan
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags $ldflags -o "build/go-pwr-windows-amd64.exe" cmd/go-pwr/main.go
if ($LASTEXITCODE -eq 0) {
  Write-Host "✓ Windows AMD64 build successful" -ForegroundColor Green
}
else {
  Write-Host "✗ Windows AMD64 build failed" -ForegroundColor Red
}

# Build for macOS AMD64
Write-Host "Building for macOS AMD64..." -ForegroundColor Cyan
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags $ldflags -o "build/go-pwr-darwin-amd64" cmd/go-pwr/main.go
if ($LASTEXITCODE -eq 0) {
  Write-Host "✓ macOS AMD64 build successful" -ForegroundColor Green
}
else {
  Write-Host "✗ macOS AMD64 build failed" -ForegroundColor Red
}

# Build for macOS ARM64
Write-Host "Building for macOS ARM64..." -ForegroundColor Cyan
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -ldflags $ldflags -o "build/go-pwr-darwin-arm64" cmd/go-pwr/main.go
if ($LASTEXITCODE -eq 0) {
  Write-Host "✓ macOS ARM64 build successful" -ForegroundColor Green
}
else {
  Write-Host "✗ macOS ARM64 build failed" -ForegroundColor Red
}

# Build for Linux AMD64
Write-Host "Building for Linux AMD64..." -ForegroundColor Cyan
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags $ldflags -o "build/go-pwr-linux-amd64" cmd/go-pwr/main.go
if ($LASTEXITCODE -eq 0) {
  Write-Host "✓ Linux AMD64 build successful" -ForegroundColor Green
}
else {
  Write-Host "✗ Linux AMD64 build failed" -ForegroundColor Red
}

# Build for Linux ARM64
Write-Host "Building for Linux ARM64..." -ForegroundColor Cyan
$env:GOOS = "linux"
$env:GOARCH = "arm64"
go build -ldflags $ldflags -o "build/go-pwr-linux-arm64" cmd/go-pwr/main.go
if ($LASTEXITCODE -eq 0) {
  Write-Host "✓ Linux ARM64 build successful" -ForegroundColor Green
}
else {
  Write-Host "✗ Linux ARM64 build failed" -ForegroundColor Red
}

# Reset environment variables
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "Build complete! Files created in build/ directory:" -ForegroundColor Green
Get-ChildItem "build/" | ForEach-Object {
  $size = [math]::Round($_.Length / 1MB, 2)
  Write-Host "  $($_.Name) (${size} MB)" -ForegroundColor White
}

Write-Host ""
Write-Host "Testing Windows build..." -ForegroundColor Cyan
& "build/go-pwr-windows-amd64.exe" -v
