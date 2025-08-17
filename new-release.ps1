# PowerShell script for complete release automation

param(
  [string]$Version = "",
  [string]$PreviousVersion = "",
  [switch]$DryRun = $false,
  [switch]$Help = $false
)

# Show help
if ($Help) {
  Write-Host "go-pwr Release Automation Script" -ForegroundColor Green
  Write-Host ""
  Write-Host "USAGE:" -ForegroundColor Yellow
  Write-Host "  .\new-release.ps1 -Version <version> [-PreviousVersion <version>] [-DryRun] [-Help]"
  Write-Host ""
  Write-Host "PARAMETERS:" -ForegroundColor Yellow
  Write-Host "  -Version         New version to release (e.g., '1.0.8')"
  Write-Host "  -PreviousVersion Previous version for changelog (auto-detected if not provided)"
  Write-Host "  -DryRun          Show what would be done without actually doing it"
  Write-Host "  -Help            Show this help message"
  Write-Host ""
  Write-Host "WHAT THIS SCRIPT DOES:" -ForegroundColor Yellow
  Write-Host "  1. Validates version format and git status"
  Write-Host "  2. Updates version in cmd/go-pwr/main.go"
  Write-Host "  3. Generates release notes template"
  Write-Host "  4. Commits and pushes version changes"
  Write-Host "  5. Builds all platform binaries"
  Write-Host "  6. Creates proper git tag and GitHub release"
  Write-Host "  7. Verifies installation works correctly"
  Write-Host ""
  Write-Host "EXAMPLES:" -ForegroundColor Yellow
  Write-Host "  .\new-release.ps1 -Version 1.0.9                    # Full release"
  Write-Host "  .\new-release.ps1 -Version 1.0.9 -DryRun            # Preview changes"
  Write-Host ""
  return
}

# Validate version parameter
if ($Version -eq "") {
  Write-Host "✗ Version parameter is required" -ForegroundColor Red
  Write-Host "Usage: .\new-release.ps1 -Version <version> [-DryRun] [-Help]" -ForegroundColor Yellow
  Write-Host "Example: .\new-release.ps1 -Version 1.0.9" -ForegroundColor Yellow
  Write-Host ""
  Write-Host "For more information, run: .\new-release.ps1 -Help" -ForegroundColor Cyan
  exit 1
}

# Validate version format
if ($Version -notmatch '^\d+\.\d+\.\d+$') {
  Write-Host "Invalid version format: $Version. Expected format: x.y.z (e.g., 1.0.8)" -ForegroundColor Red
  exit 1
}

Write-Host "🚀 Starting release process for go-pwr v$Version" -ForegroundColor Green
Write-Host ""

# Check if we're in a git repository
if (-not (Test-Path ".git")) {
  Write-Host "✗ Not in a git repository" -ForegroundColor Red
  exit 1
}

# Check git status
$gitStatus = git status --porcelain
if ($gitStatus -and -not $DryRun) {
  Write-Host "✗ Working directory is not clean. Please commit or stash changes first." -ForegroundColor Red
  Write-Host "Uncommitted changes:" -ForegroundColor Yellow
  $gitStatus | ForEach-Object { Write-Host "  $_" -ForegroundColor Yellow }
  exit 1
}

# Auto-detect previous version if not provided
if ($PreviousVersion -eq "") {
  $tags = git tag --sort=-version:refname | Where-Object { $_ -match '^v\d+\.\d+\.\d+$' }
  if ($tags.Count -gt 0) {
    $PreviousVersion = $tags[0] -replace '^v', ''
    Write-Host "✓ Detected previous version: $PreviousVersion" -ForegroundColor Yellow
  }
  else {
    $PreviousVersion = "1.0.0"
    Write-Host "✓ No previous tags found, using: $PreviousVersion" -ForegroundColor Yellow
  }
}

# Check if version already exists
$existingTag = git tag -l "v$Version"
if ($existingTag) {
  Write-Host "✗ Version v$Version already exists as a git tag" -ForegroundColor Red
  exit 1
}

# Validate main.go exists
if (-not (Test-Path "cmd/go-pwr/main.go")) {
  Write-Host "✗ cmd/go-pwr/main.go not found" -ForegroundColor Red
  exit 1
}

# Step 1: Update version in main.go
Write-Host "📝 Step 1: Updating version in main.go..." -ForegroundColor Cyan
$mainGoContent = Get-Content "cmd/go-pwr/main.go" -Raw
$currentVersionPattern = 'const version = "([^"]+)"'
$currentVersionMatch = [regex]::Match($mainGoContent, $currentVersionPattern)

if (-not $currentVersionMatch.Success) {
  Write-Host "✗ Could not find version constant in main.go" -ForegroundColor Red
  exit 1
}

$currentVersion = $currentVersionMatch.Groups[1].Value
Write-Host "  Current version: $currentVersion" -ForegroundColor White
Write-Host "  New version: $Version" -ForegroundColor White

if (-not $DryRun) {
  $newContent = $mainGoContent -replace $currentVersionPattern, "const version = `"$Version`""
  $newContent | Set-Content "cmd/go-pwr/main.go" -NoNewline
  Write-Host "✓ Version updated in main.go" -ForegroundColor Green
}
else {
  Write-Host "✓ Would update version in main.go" -ForegroundColor Yellow
}

# Step 2: Generate release notes template
Write-Host ""
Write-Host "📋 Step 2: Generating release notes template..." -ForegroundColor Cyan
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

if (-not $DryRun) {
  $template | Out-File -FilePath $releaseNotesFile -Encoding UTF8
  Write-Host "✓ Release notes template created: $releaseNotesFile" -ForegroundColor Green
}
else {
  Write-Host "✓ Would create release notes template: $releaseNotesFile" -ForegroundColor Yellow
}

# Step 3: Clean old binaries
Write-Host ""
Write-Host "🧹 Step 3: Cleaning old binaries..." -ForegroundColor Cyan
if (Test-Path "build") {
  if (-not $DryRun) {
    Remove-Item "build/*" -Force -ErrorAction SilentlyContinue
    Write-Host "✓ Cleaned build directory" -ForegroundColor Green
  }
  else {
    Write-Host "✓ Would clean build directory" -ForegroundColor Yellow
  }
}

if (Test-Path "go-pwr.exe") {
  if (-not $DryRun) {
    Remove-Item "go-pwr.exe" -Force
    Write-Host "✓ Removed local go-pwr.exe" -ForegroundColor Green
  }
  else {
    Write-Host "✓ Would remove local go-pwr.exe" -ForegroundColor Yellow
  }
}

if (-not $DryRun) {
  Write-Host ""
  Write-Host "⏸️  MANUAL STEP REQUIRED" -ForegroundColor Magenta
  Write-Host "Please edit the release notes file: $releaseNotesFile" -ForegroundColor White
  Write-Host "Update the sections with actual changes for this release." -ForegroundColor White
  Write-Host ""
  Write-Host "Opening release notes file for editing..." -ForegroundColor Cyan
  Start-Process notepad.exe $releaseNotesFile
  
  Write-Host ""
  Read-Host "Press Enter when you've finished editing the release notes and are ready to continue"
  
  # Step 4: Commit and push changes
  Write-Host ""
  Write-Host "💾 Step 4: Committing and pushing changes..." -ForegroundColor Cyan
  git add cmd/go-pwr/main.go $releaseNotesFile
  git commit -m "bump: Update version to $Version and add release notes"
  
  if ($LASTEXITCODE -eq 0) {
    git push origin main
    if ($LASTEXITCODE -eq 0) {
      Write-Host "✓ Changes committed and pushed to main" -ForegroundColor Green
    }
    else {
      Write-Host "✗ Failed to push changes" -ForegroundColor Red
      exit 1
    }
  }
  else {
    Write-Host "✗ Failed to commit changes" -ForegroundColor Red
    exit 1
  }
  
  # Step 5: Build all platforms
  Write-Host ""
  Write-Host "🔨 Step 5: Building all platforms..." -ForegroundColor Cyan
  & ".\build-all.ps1" -Version $Version
  
  if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Build failed" -ForegroundColor Red
    exit 1
  }
  
  # Step 6: Create git tag and GitHub release
  Write-Host ""
  Write-Host "🏷️  Step 6: Creating git tag and GitHub release..." -ForegroundColor Cyan
  
  # Create and push tag
  git tag -a "v$Version" -m "go-pwr v$Version"
  git push origin "v$Version"
  
  if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Failed to create/push git tag" -ForegroundColor Red
    exit 1
  }
  
  # Check if gh CLI is available
  $ghAvailable = Get-Command "gh" -ErrorAction SilentlyContinue
  if (-not $ghAvailable) {
    Write-Host "✗ GitHub CLI (gh) not found. Please install it to create releases automatically." -ForegroundColor Red
    Write-Host "  Install: winget install GitHub.cli" -ForegroundColor Yellow
    Write-Host "  Manual release: https://github.com/rocketpowerinc/go-pwr/releases/new" -ForegroundColor Yellow
    exit 1
  }
  
  # Create GitHub release
  gh release create "v$Version" build/* --title "go-pwr v$Version" --notes-file $releaseNotesFile --latest
  
  if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ GitHub release created successfully!" -ForegroundColor Green
    Write-Host "  View release: https://github.com/rocketpowerinc/go-pwr/releases/tag/v$Version" -ForegroundColor Cyan
  }
  else {
    Write-Host "✗ Failed to create GitHub release" -ForegroundColor Red
    exit 1
  }
  
  # Step 7: Update local binaries
  Write-Host ""
  Write-Host "🔄 Step 7: Updating local binaries..." -ForegroundColor Cyan
  
  # Copy built binary to local directory
  Copy-Item "build/go-pwr-windows-amd64.exe" "go-pwr.exe" -Force
  Write-Host "✓ Updated local go-pwr.exe" -ForegroundColor Green
  
  # Copy to Go bin directory if it exists
  $goBinPath = Join-Path $env:GOPATH "bin\go-pwr.exe"
  if (-not $env:GOPATH) {
    $goBinPath = Join-Path $env:USERPROFILE "go\bin\go-pwr.exe"
  }
  
  $goBinDir = Split-Path $goBinPath -Parent
  if (Test-Path $goBinDir) {
    Copy-Item "build/go-pwr-windows-amd64.exe" $goBinPath -Force
    Write-Host "✓ Updated Go bin directory: $goBinPath" -ForegroundColor Green
  }
  
  # Step 8: Verify installation
  Write-Host ""
  Write-Host "✅ Step 8: Verifying installation..." -ForegroundColor Cyan
  
  # Test local binary
  $localVersion = & "./go-pwr.exe" -v 2>$null | Select-String "go-pwr v"
  if ($localVersion -and $localVersion.ToString().Contains("v$Version")) {
    Write-Host "✓ Local binary version correct: $($localVersion.ToString().Trim())" -ForegroundColor Green
  }
  else {
    Write-Host "⚠ Local binary version check failed" -ForegroundColor Yellow
  }
  
  # Test Go install after a brief delay
  Write-Host ""
  Write-Host "Testing go install after release (may take a moment for Go proxy to update)..." -ForegroundColor Cyan
  Start-Sleep -Seconds 10
  
  $tempDir = Join-Path $env:TEMP "go-pwr-release-test"
  New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
  Push-Location $tempDir
  
  try {
    go install -v "github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v$Version"
    if ($LASTEXITCODE -eq 0) {
      Write-Host "✓ go install test successful!" -ForegroundColor Green
    }
    else {
      Write-Host "⚠ go install test failed (Go proxy may need time to update)" -ForegroundColor Yellow
    }
  }
  finally {
    Pop-Location
    Remove-Item -Recurse -Force $tempDir -ErrorAction SilentlyContinue
  }
  
  # Final summary
  Write-Host ""
  Write-Host "🎉 Release v$Version completed successfully!" -ForegroundColor Green
  Write-Host ""
  Write-Host "📋 Summary:" -ForegroundColor Yellow
  Write-Host "  ✓ Version updated in main.go" -ForegroundColor White
  Write-Host "  ✓ Release notes created and edited" -ForegroundColor White
  Write-Host "  ✓ Changes committed and pushed" -ForegroundColor White
  Write-Host "  ✓ All platform binaries built" -ForegroundColor White
  Write-Host "  ✓ Git tag v$Version created and pushed" -ForegroundColor White
  Write-Host "  ✓ GitHub release published" -ForegroundColor White
  Write-Host "  ✓ Local binaries updated" -ForegroundColor White
  Write-Host ""
  Write-Host "🔗 Links:" -ForegroundColor Yellow
  Write-Host "  Release: https://github.com/rocketpowerinc/go-pwr/releases/tag/v$Version" -ForegroundColor Cyan
  Write-Host "  Install: go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v$Version" -ForegroundColor Cyan
  Write-Host ""
  
}
else {
  Write-Host ""
  Write-Host "🔍 DRY RUN COMPLETE" -ForegroundColor Yellow
  Write-Host "The following would be done:" -ForegroundColor White
  Write-Host "  1. Update version in cmd/go-pwr/main.go from $currentVersion to $Version" -ForegroundColor Gray
  Write-Host "  2. Create release notes file: $releaseNotesFile" -ForegroundColor Gray
  Write-Host "  3. Clean old binaries" -ForegroundColor Gray
  Write-Host "  4. Commit and push changes" -ForegroundColor Gray
  Write-Host "  5. Build all platform binaries" -ForegroundColor Gray
  Write-Host "  6. Create git tag v$Version and push" -ForegroundColor Gray
  Write-Host "  7. Create GitHub release with binaries" -ForegroundColor Gray
  Write-Host "  8. Update local binaries" -ForegroundColor Gray
  Write-Host "  9. Verify installation" -ForegroundColor Gray
  Write-Host ""
  Write-Host "To execute the release, run without -DryRun flag:" -ForegroundColor Cyan
  Write-Host "  .\new-release.ps1 -Version $Version" -ForegroundColor White
}
