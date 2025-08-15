# PowerShell script to build go-pwr for all platforms and optionally create release

param(
  [string]$Version = "",
  [switch]$Release = $false,
  [switch]$Help = $false
)

# Show help
if ($Help) {
  Write-Host "go-pwr Build Script" -ForegroundColor Green
  Write-Host ""
  Write-Host "USAGE:" -ForegroundColor Yellow
  Write-Host "  .\build-all.ps1 [-Version <version>] [-Release] [-Help]"
  Write-Host ""
  Write-Host "PARAMETERS:" -ForegroundColor Yellow
  Write-Host "  -Version    Specify version (e.g., '1.0.5'). If not provided, will extract from main.go"
  Write-Host "  -Release    Create GitHub release after successful build"
  Write-Host "  -Help       Show this help message"
  Write-Host ""
  Write-Host "EXAMPLES:" -ForegroundColor Yellow
  Write-Host "  .\build-all.ps1                    # Build only"
  Write-Host "  .\build-all.ps1 -Version 1.0.5     # Build with specific version"
  Write-Host "  .\build-all.ps1 -Release           # Build and create GitHub release"
  Write-Host "  .\build-all.ps1 -Version 1.0.5 -Release  # Build specific version and release"
  Write-Host ""
  return
}

# Extract version from main.go if not provided
if ($Version -eq "") {
  $versionLine = Get-Content "cmd/go-pwr/main.go" | Select-String 'const version = "(.*)"'
  if ($versionLine) {
    $Version = $versionLine.Matches[0].Groups[1].Value
    Write-Host "Detected version from main.go: $Version" -ForegroundColor Yellow
  }
  else {
    Write-Host "Could not detect version from main.go. Please specify -Version parameter." -ForegroundColor Red
    exit 1
  }
}

# Validate version format
if ($Version -notmatch '^\d+\.\d+\.\d+$') {
  Write-Host "Invalid version format: $Version. Expected format: x.y.z (e.g., 1.0.4)" -ForegroundColor Red
  exit 1
}# Get git commit hash (short)
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

# Release creation (optional)
if ($Release) {
  Write-Host ""
  Write-Host "Creating GitHub release..." -ForegroundColor Magenta
    
  # Check if gh CLI is available
  $ghAvailable = Get-Command "gh" -ErrorAction SilentlyContinue
  if (-not $ghAvailable) {
    Write-Host "✗ GitHub CLI (gh) not found. Please install it to create releases automatically." -ForegroundColor Red
    Write-Host "  Install: winget install GitHub.cli" -ForegroundColor Yellow
    Write-Host "  Manual release: https://github.com/rocketpowerinc/go-pwr/releases/new" -ForegroundColor Yellow
    return
  }
    
  # Check if release notes file exists
  $releaseNotesFile = "RELEASE_NOTES_v$Version.md"
  if (-not (Test-Path $releaseNotesFile)) {
    Write-Host "✗ Release notes file not found: $releaseNotesFile" -ForegroundColor Red
    Write-Host "  Please create release notes file first." -ForegroundColor Yellow
    return
  }
    
  # Check if tag already exists
  $tagExists = git tag -l "v$Version"
  if (-not $tagExists) {
    Write-Host "Creating git tag v$Version..." -ForegroundColor Yellow
    git tag -a "v$Version" -m "go-pwr v$Version"
    git push origin "v$Version"
    if ($LASTEXITCODE -ne 0) {
      Write-Host "✗ Failed to create/push git tag" -ForegroundColor Red
      return
    }
  }
    
  # Create GitHub release
  Write-Host "Creating GitHub release v$Version..." -ForegroundColor Yellow
  gh release create "v$Version" build/* --title "go-pwr v$Version" --notes-file $releaseNotesFile --latest
    
  if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ GitHub release created successfully!" -ForegroundColor Green
    Write-Host "  View release: https://github.com/rocketpowerinc/go-pwr/releases/tag/v$Version" -ForegroundColor Cyan
        
    # Test go install after a brief delay
    Write-Host ""
    Write-Host "Testing go install after release..." -ForegroundColor Cyan
    Start-Sleep -Seconds 5
        
    $tempDir = Join-Path $env:TEMP "go-pwr-test"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    Push-Location $tempDir
        
    try {
      go install -v "github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v$Version"
      if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ go install test successful!" -ForegroundColor Green
      }
      else {
        Write-Host "⚠ go install test failed (may need time for Go proxy to update)" -ForegroundColor Yellow
      }
    }
    finally {
      Pop-Location
      Remove-Item -Recurse -Force $tempDir -ErrorAction SilentlyContinue
    }
  }
  else {
    Write-Host "✗ Failed to create GitHub release" -ForegroundColor Red
  }
}
else {
  Write-Host ""
  Write-Host "Build complete! To create a release, run:" -ForegroundColor Green
  Write-Host "  .\build-all.ps1 -Version $Version -Release" -ForegroundColor Cyan
}
