package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/rocketpowerinc/go-pwr/internal/app"
	"github.com/rocketpowerinc/go-pwr/internal/config"
)

const version = "1.0.6"

// These will be set at build time using -ldflags
var (
	gitCommit = "unknown"
	buildDate = "unknown"
)

func main() {
	// Custom usage function to show our help
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "go-pwr v%s - Cross-platform script launcher\n\n", version)
		fmt.Fprintf(os.Stderr, "USAGE:\n")
		fmt.Fprintf(os.Stderr, "  go-pwr [flags]\n\n")
		fmt.Fprintf(os.Stderr, "FLAGS:\n")
		fmt.Fprintf(os.Stderr, "  -h, -help           Show this help message\n")
		fmt.Fprintf(os.Stderr, "  -v, -version        Show version information\n")
		fmt.Fprintf(os.Stderr, "  -show-repo          Show the current repository URL\n")
		fmt.Fprintf(os.Stderr, "  -set-repo string    Set a custom repository URL\n")
		fmt.Fprintf(os.Stderr, "  -reset-repo         Reset to the default repository\n\n")
		fmt.Fprintf(os.Stderr, "EXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "  go-pwr                                           Start the interactive TUI\n")
		fmt.Fprintf(os.Stderr, "  go-pwr -show-repo                               Show current repository\n")
		fmt.Fprintf(os.Stderr, "  go-pwr -set-repo https://github.com/user/repo.git  Set custom repository\n")
		fmt.Fprintf(os.Stderr, "  go-pwr -reset-repo                              Reset to default repository\n\n")
		fmt.Fprintf(os.Stderr, "For more information, visit: https://github.com/rocketpowerinc/go-pwr\n")
	}

	// Parse command line flags
	var setRepo = flag.String("set-repo", "", "Set a custom repository URL")
	var resetRepo = flag.Bool("reset-repo", false, "Reset to the default repository")
	var showRepo = flag.Bool("show-repo", false, "Show the current repository URL")
	var showVersion = flag.Bool("version", false, "Show version information")
	var showVersionShort = flag.Bool("v", false, "Show version information")
	var showHelp = flag.Bool("help", false, "Show help information")
	var showHelpShort = flag.Bool("h", false, "Show help information")

	flag.Parse()

	// Handle help flags
	if *showHelp || *showHelpShort {
		flag.Usage()
		return
	}

	// Handle version flags
	if *showVersion || *showVersionShort {
		fmt.Printf("go-pwr v%s\n", version)
		fmt.Printf("Git commit: %s\n", gitCommit)
		fmt.Printf("Build date: %s\n", buildDate)
		fmt.Printf("Built with Go %s for %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Repository: https://github.com/rocketpowerinc/go-pwr\n")
		return
	}

	// Handle repository management flags
	if *showRepo {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Current repository: %s\n", cfg.RepoURL)
		fmt.Printf("Default repository: %s\n", config.GetDefaultRepoURL())
		return
	}

	if *resetRepo {
		if err := config.ResetToDefaultRepo(); err != nil {
			fmt.Fprintf(os.Stderr, "Error resetting repository: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Repository reset to default: %s\n", config.GetDefaultRepoURL())
		return
	}

	if *setRepo != "" {
		if err := config.ValidateRepoURL(*setRepo); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid repository URL: %v\n", err)
			os.Exit(1)
		}
		if err := config.SaveRepoURL(*setRepo); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving repository URL: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Repository set to: %s\n", *setRepo)
		return
	}

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
