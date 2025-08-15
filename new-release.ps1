# PowerShell script to generate release notes template

param(
  [Parameter(Mandatory = $true)]
  [string]$Version,
  [string]$PreviousVersion = ""
)

# Validate version format
if ($Version -notmatch '^\d+\.\d+\.\d+$') {
  Write-Host "Invalid version format: $Version. Expected format: x.y.z (e.g., 1.0.5)" -ForegroundColor Red
  exit 1
}

# Auto-detect previous version if not provided
if ($PreviousVersion -eq "") {
  $tags = git tag --sort=-version:refname | Where-Object { $_ -match '^v\d+\.\d+\.\d+$' }
  if ($tags.Count -gt 0) {
    $PreviousVersion = $tags[0] -replace '^v', ''
    Write-Host "Detected previous version: $PreviousVersion" -ForegroundColor Yellow
  }
  else {
    $PreviousVersion = "1.0.0"
    Write-Host "No previous tags found, using: $PreviousVersion" -ForegroundColor Yellow
  }
}

$releaseNotesFile = "RELEASE_NOTES_v$Version.md"

# Get git commit hash
$gitCommit = git rev-parse --short HEAD
$buildDate = Get-Date -Format "yyyy-MM-dd"

# Generate template
$template = @"
# Release Notes - go-pwr v$Version

## 🆕 New Features

### [Add new features here]
- Feature 1
- Feature 2

## 🔧 Improvements

### [Add improvements here]
- Improvement 1
- Improvement 2

## 🐛 Bug Fixes

### [Add bug fixes here]
- Fix 1
- Fix 2

## 📚 Documentation Updates

- Updated documentation
- Added examples

## 🔧 Technical Details

- **Git Commit**: $gitCommit
- **Build Date**: $buildDate
- **Go Version**: $(go version)

## 📦 Available Downloads

- **Windows**: ``go-pwr-windows-amd64.exe``
- **macOS Intel**: ``go-pwr-darwin-amd64``
- **macOS Apple Silicon**: ``go-pwr-darwin-arm64``
- **Linux AMD64**: ``go-pwr-linux-amd64``
- **Linux ARM64**: ``go-pwr-linux-arm64``

## 🔧 Installation

### Go Install (Recommended)
``````bash
# Install latest version
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest

# Or install specific version
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v$Version
``````

### Direct Download
1. Download the appropriate binary for your platform from the assets above
2. Make it executable (Unix-like systems): ``chmod +x go-pwr-*``
3. Move to a directory in your PATH or run directly

## ⚠️ Breaking Changes

None - This release is fully backward compatible with previous versions.

## 🔄 Upgrade Instructions

No special upgrade steps required. Simply install the new version.

---

**Full Changelog**: https://github.com/rocketpowerinc/go-pwr/compare/v$PreviousVersion...v$Version
"@

# Write template to file
$template | Out-File -FilePath $releaseNotesFile -Encoding UTF8

Write-Host "✓ Release notes template created: $releaseNotesFile" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Edit $releaseNotesFile with your changes"
Write-Host "2. Update version in cmd/go-pwr/main.go if needed"
Write-Host "3. Run: .\build-all.ps1 -Version $Version -Release"
Write-Host ""
Write-Host "Opening release notes file for editing..." -ForegroundColor Cyan
Start-Process notepad.exe $releaseNotesFile
