// Package config handles configuration management for go-pwr.
package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Config holds the application configuration.
type Config struct {
	ScriptbinPath string `json:"scriptbin_path"`
	RepoURL       string `json:"repo_url"`
	Theme         string `json:"theme"` // Store the theme name
}

// UserConfig represents the persistent user configuration
type UserConfig struct {
	Theme   string `json:"theme"`
	RepoURL string `json:"repo_url,omitempty"` // Custom repository URL
}

// Load loads the application configuration.
func Load() (*Config, error) {
	scriptbinPath, err := getScriptbinPath()
	if err != nil {
		return nil, err
	}

	// Default repository URL
	defaultRepoURL := "https://github.com/rocketpowerinc/scriptbin.git"

	config := &Config{
		ScriptbinPath: scriptbinPath,
		RepoURL:       defaultRepoURL,
		Theme:         "Ocean Breeze", // Default theme
	}

	// Load user preferences
	if userConfig, err := loadUserConfig(); err == nil {
		if userConfig.Theme != "" {
			config.Theme = userConfig.Theme
		}
		// Use custom repo URL if set, otherwise keep default
		if userConfig.RepoURL != "" {
			config.RepoURL = userConfig.RepoURL
		}
	}

	return config, nil
}

// SaveTheme saves the user's theme preference
func SaveTheme(themeName string) error {
	userConfig, _ := loadUserConfig() // Load existing config or create new
	if userConfig == nil {
		userConfig = &UserConfig{}
	}
	
	userConfig.Theme = themeName
	return saveUserConfig(userConfig)
}

// SaveRepoURL saves the user's custom repository URL
func SaveRepoURL(repoURL string) error {
	userConfig, _ := loadUserConfig() // Load existing config or create new
	if userConfig == nil {
		userConfig = &UserConfig{}
	}
	
	userConfig.RepoURL = repoURL
	return saveUserConfig(userConfig)
}

// ResetToDefaultRepo resets the repository to the default scriptbin
func ResetToDefaultRepo() error {
	userConfig, _ := loadUserConfig() // Load existing config or create new
	if userConfig == nil {
		userConfig = &UserConfig{}
	}
	
	userConfig.RepoURL = "" // Empty string means use default
	return saveUserConfig(userConfig)
}

// GetDefaultRepoURL returns the default repository URL
func GetDefaultRepoURL() string {
	return "https://github.com/rocketpowerinc/scriptbin.git"
}

// saveUserConfig saves the complete user configuration
func saveUserConfig(userConfig *UserConfig) error {
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

// ValidateRepoURL validates if a repository URL is properly formatted
func ValidateRepoURL(repoURL string) error {
	if repoURL == "" {
		return fmt.Errorf("repository URL cannot be empty")
	}

	// Parse the URL
	u, err := url.Parse(repoURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	// Check if it's a supported scheme
	if u.Scheme != "https" && u.Scheme != "http" && u.Scheme != "git" && u.Scheme != "ssh" {
		return fmt.Errorf("unsupported URL scheme: %s (supported: https, http, git, ssh)", u.Scheme)
	}

	// Basic validation for git repositories
	if !strings.HasSuffix(strings.ToLower(repoURL), ".git") {
		return fmt.Errorf("URL should end with .git for git repositories")
	}

	// Additional validation for GitHub URLs
	if strings.Contains(u.Host, "github.com") {
		pathParts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(pathParts) != 2 {
			return fmt.Errorf("GitHub URLs should be in format: https://github.com/owner/repo.git")
		}
	}

	return nil
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
