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
// It now uses a fixed location in $HOME/Downloads/Temp/scriptbin
func getScriptbinPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}

	// Use fixed location in Downloads/Temp/scriptbin
	scriptbinPath := filepath.Join(homeDir, "Downloads", "Temp", "scriptbin")
	
	// Ensure the parent directories exist
	parentDir := filepath.Dir(scriptbinPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directories: %v", err)
	}

	return scriptbinPath, nil
}
