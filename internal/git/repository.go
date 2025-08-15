// Package git handles git repository operations for go-pwr.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rocketpowerinc/go-pwr/internal/config"
)

// EnsureRepository ensures the script repository is cloned and up to date.
func EnsureRepository(cfg *config.Config) error {
	// Generate a unique path based on the repository URL
	scriptPath := getRepositoryPath(cfg)
	
	// Always remove and re-clone for fresh content
	if _, err := os.Stat(scriptPath); err == nil {
		if err := os.RemoveAll(scriptPath); err != nil {
			return fmt.Errorf("failed to remove old repository: %v", err)
		}
	}

	// Ensure parent directory exists
	parentDir := filepath.Dir(scriptPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %v", err)
	}

	cmd := exec.Command("git", "clone", cfg.RepoURL, scriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone error: %v\n%s", err, string(out))
	}

	// Update the config with the actual path used
	cfg.ScriptbinPath = scriptPath
	return nil
}

// getRepositoryPath generates a path for the repository based on the URL
func getRepositoryPath(cfg *config.Config) string {
	// If it's the default repository, use the original path
	if cfg.RepoURL == config.GetDefaultRepoURL() {
		return cfg.ScriptbinPath
	}
	
	// For custom repositories, create a unique directory name
	// Extract repository name from URL (remove .git suffix)
	repoName := filepath.Base(cfg.RepoURL)
	if filepath.Ext(repoName) == ".git" {
		repoName = repoName[:len(repoName)-4]
	}
	
	// Use the parent directory of the original scriptbin path
	parentDir := filepath.Dir(cfg.ScriptbinPath)
	return filepath.Join(parentDir, "custom-"+repoName)
}
