// Package config handles configuration management for go-pwr.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration.
type Config struct {
	ScriptbinPath string `json:"scriptbin_path"`
	RepoURL       string `json:"repo_url"`
	Theme         string `json:"theme"` // Store the theme name
}

// UserConfig represents the persistent user configuration
type UserConfig struct {
	Theme string `json:"theme"`
}

// Load loads the application configuration.
func Load() (*Config, error) {
	scriptbinPath, err := getScriptbinPath()
	if err != nil {
		return nil, err
	}

	config := &Config{
		ScriptbinPath: scriptbinPath,
		RepoURL:       "https://github.com/rocketpowerinc/scriptbin.git",
		Theme:         "Ocean Breeze", // Default theme
	}

	// Load user preferences
	if userConfig, err := loadUserConfig(); err == nil {
		if userConfig.Theme != "" {
			config.Theme = userConfig.Theme
		}
	}

	return config, nil
}

// SaveTheme saves the user's theme preference
func SaveTheme(themeName string) error {
	userConfig := UserConfig{
		Theme: themeName,
	}

	configPath, err := getUserConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.MarshalIndent(userConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// loadUserConfig loads the user configuration from disk
func loadUserConfig() (*UserConfig, error) {
	configPath, err := getUserConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err // File doesn't exist or can't be read
	}

	var userConfig UserConfig
	if err := json.Unmarshal(data, &userConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &userConfig, nil
}

// getUserConfigPath returns the path to the user config file
func getUserConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}

	// Use platform-appropriate config directory
	var configDir string
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		configDir = os.Getenv("XDG_CONFIG_HOME")
	} else {
		configDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configDir, "go-pwr", "config.json"), nil
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
