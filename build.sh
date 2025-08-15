#!/bin/bash
# Bash build script for go-pwr with git commit information

# Get git commit hash (short)
gitCommit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Get current date/time in ISO format
buildDate=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags to inject version information
ldflags="-X main.gitCommit=$gitCommit -X main.buildDate=$buildDate"

echo "Building go-pwr..."
echo "Git commit: $gitCommit"
echo "Build date: $buildDate"

# Build the application
go build -ldflags "$ldflags" -o go-pwr cmd/go-pwr/main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Testing version output:"
    ./go-pwr -v
else
    echo "Build failed!"
    exit 1
fi
