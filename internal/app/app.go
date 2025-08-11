// Package app contains the main application logic for go-pwr.
package app

import (
	"github.com/rocketpowerinc/go-pwr/internal/config"
	"github.com/rocketpowerinc/go-pwr/internal/git"
	"github.com/rocketpowerinc/go-pwr/internal/ui"
)

// Run starts the go-pwr application.
func Run() error {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Ensure repository is cloned/updated
	if err := git.EnsureRepository(cfg); err != nil {
		return err
	}

	// Start the UI
	return ui.Start(cfg)
}
