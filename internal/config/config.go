// Package config handles configuration management for go-pwr.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration.
type Config struct {
	ScriptbinPath string
	RepoURL       string
}

// Load loads the application configuration.
func Load() (*Config, error) {
	scriptbinPath, err := getScriptbinPath()
	if err != nil {
		return nil, err
	}

	return &Config{
		ScriptbinPath: scriptbinPath,
		RepoURL:       "https://github.com/rocketpowerinc/scriptbin.git",
	}, nil
}

// getScriptbinPath returns the path where scriptbin should be stored
// It tries to find the go-pwr installation directory and places scriptbin there
func getScriptbinPath() (string, error) {
	// Get the executable path
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}

	// Get the directory containing the executable
	exeDir := filepath.Dir(exePath)

	// Try to find the go-pwr source directory
	// Check if we're running from the development directory (has main.go)
	if _, err := os.Stat(filepath.Join(exeDir, "main.go")); err == nil {
		// We're in the development directory
		return filepath.Join(exeDir, "scriptbin"), nil
	}

	// Check if there's a go-pwr directory in common locations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}

	// Common locations where go-pwr might be cloned
	possiblePaths := []string{
		filepath.Join(homeDir, "go-pwr"),
		filepath.Join(homeDir, "Github-pwr", "go-pwr"),
		filepath.Join(homeDir, "projects", "go-pwr"),
		filepath.Join(homeDir, "src", "go-pwr"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(filepath.Join(path, "main.go")); err == nil {
			return filepath.Join(path, "scriptbin"), nil
		}
	}

	// If we can't find the source directory, create scriptbin next to the executable
	// This handles the case where go-pwr is installed via go install
	return filepath.Join(exeDir, "scriptbin"), nil
}
