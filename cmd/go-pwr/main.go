// Package main is the entry point for the go-pwr application.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/rocketpowerinc/go-pwr/internal/app"
)

func main() {
	// Show tmux warning for Linux users
	if runtime.GOOS == "linux" {
		showTmuxWarning()
	}

	// Check if we should run in tmux on Linux
	if shouldRunInTmux() {
		if err := runInTmux(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start in tmux: %v\n", err)
			// Fall back to normal execution
		} else {
			return // Successfully started in tmux, exit this process
		}
	}

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// showTmuxWarning displays a prominent warning for Linux users about tmux
func showTmuxWarning() {
	// Don't show warning if already in tmux
	if os.Getenv("TMUX") != "" {
		return
	}

	// Don't show warning if explicitly disabled
	if os.Getenv("GO_PWR_NO_TMUX_WARNING") != "" {
		return
	}

	fmt.Println("\n" + strings.Repeat("═", 70))
	fmt.Println("⚠️  IMPORTANT: For the best experience on Linux, run go-pwr in tmux!")
	fmt.Println("")
	fmt.Println("   Quick start with tmux:")
	fmt.Println("   $ tmux new-session go-pwr")
	fmt.Println("")
	fmt.Println("   Benefits:")
	fmt.Println("   • Session persistence (survive SSH disconnects)")
	fmt.Println("   • Better script execution handling")
	fmt.Println("   • Background operation support")
	fmt.Println("")
	fmt.Println("   To disable this warning: export GO_PWR_NO_TMUX_WARNING=1")
	fmt.Println(strings.Repeat("═", 70) + "\n")
	
	// Give user a moment to read the warning
	fmt.Print("Press Enter to continue without tmux, or Ctrl+C to exit and use tmux...")
	fmt.Scanln()
	fmt.Println()
}

// shouldRunInTmux checks if we should automatically start in tmux
func shouldRunInTmux() bool {
	// Only on Linux
	if runtime.GOOS != "linux" {
		return false
	}

	// Don't run in tmux if we're already in tmux
	if os.Getenv("TMUX") != "" {
		return false
	}

	// Don't run in tmux if explicitly disabled
	if os.Getenv("GO_PWR_NO_TMUX") != "" {
		return false
	}

	// Check if tmux is available
	if _, err := exec.LookPath("tmux"); err != nil {
		return false
	}

	// Check if we're in a desktop environment - if so, let the normal GUI terminal logic handle it
	if isDesktopEnvironment() {
		return false
	}

	return true
}

// runInTmux starts the application in a new tmux session
func runInTmux() error {
	// Get the current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Create a new tmux session with go-pwr
	sessionName := "go-pwr-main"
	
	// Set environment variable to prevent recursive tmux launching
	env := append(os.Environ(), "GO_PWR_NO_TMUX=1")
	
	cmd := exec.Command("tmux", "new-session", "-s", sessionName, execPath)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// isDesktopEnvironment checks if we're running in a desktop environment
func isDesktopEnvironment() bool {
	// Check for DISPLAY environment variable (X11)
	if display := os.Getenv("DISPLAY"); display != "" {
		return true
	}

	// Check for Wayland environment
	if wayland := os.Getenv("WAYLAND_DISPLAY"); wayland != "" {
		return true
	}

	// Check for common desktop session variables
	if session := os.Getenv("DESKTOP_SESSION"); session != "" {
		return true
	}

	if xdg := os.Getenv("XDG_SESSION_TYPE"); xdg != "" {
		return true
	}

	return false
}
