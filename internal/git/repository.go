// Package git handles git repository operations for go-pwr.
package git

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rocketpowerinc/go-pwr/internal/config"
)

// EnsureRepository ensures the scriptbin repository is cloned and up to date.
func EnsureRepository(cfg *config.Config) error {
	// Always remove and re-clone for fresh content
	if _, err := os.Stat(cfg.ScriptbinPath); err == nil {
		if err := os.RemoveAll(cfg.ScriptbinPath); err != nil {
			return fmt.Errorf("failed to remove old scriptbin: %v", err)
		}
	}

	cmd := exec.Command("git", "clone", cfg.RepoURL, cfg.ScriptbinPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone error: %v\n%s", err, string(out))
	}

	return nil
}
