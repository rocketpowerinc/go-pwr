# PowerShell build script for go-pwr with git commit information

# Get git commit hash (short)
$gitCommit = git rev-parse --short HEAD
if ($LASTEXITCODE -ne 0) {
    $gitCommit = "unknown"
}

# Get current date/time in ISO format
$buildDate = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"

# Build flags to inject version information
$ldflags = "-X main.gitCommit=$gitCommit -X main.buildDate=$buildDate"

Write-Host "Building go-pwr..." -ForegroundColor Green
Write-Host "Git commit: $gitCommit" -ForegroundColor Yellow
Write-Host "Build date: $buildDate" -ForegroundColor Yellow

# Build the application
go build -ldflags $ldflags -o go-pwr.exe cmd/go-pwr/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
    Write-Host "Testing version output:" -ForegroundColor Cyan
    .\go-pwr.exe -v
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
