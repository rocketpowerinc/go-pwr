// Package platform provides platform-specific utilities for go-pwr.
package platform

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// IsWindows returns true if running on Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsMac returns true if running on macOS.
func IsMac() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux returns true if running on Linux.
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// GetOSName returns a human-readable OS name.
func GetOSName() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows"
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	default:
		return runtime.GOOS
	}
}

// ExecuteScript runs a script in a new terminal window based on the platform.
func ExecuteScript(scriptPath, scriptName string) error {
	var cmd *exec.Cmd

	if IsWindows() {
		if strings.HasSuffix(scriptName, ".ps1") {
			cmd = exec.Command("cmd", "/C", "start", "powershell", "-NoExit", "-Command", "Clear-Host; "+scriptPath+"; Write-Host ''; Read-Host 'Press Enter to exit'")
		} else {
			cmd = exec.Command("cmd", "/C", "start", "cmd", "/K", "cls && bash -l "+scriptPath+" & pause")
		}
	} else if IsMac() {
		// Improved macOS terminal handling
		scriptCmd := scriptPath
		if strings.HasSuffix(scriptName, ".ps1") {
			scriptCmd = "pwsh " + scriptPath
		} else {
			scriptCmd = "bash " + scriptPath
		}
		osaCmd := fmt.Sprintf(`tell application "Terminal"
    do script "clear; %s; echo; read -n 1 -s -r -p 'Press any key to exit...'"
    activate
end tell`, scriptCmd)
		cmd = exec.Command("osascript", "-e", osaCmd)
	} else if IsLinux() {
		// Check if we're in a server/headless environment (no DISPLAY)
		if !IsDesktopEnvironment() {
			// Server environment: run in current terminal with tmux/screen if available
			return ExecuteInCurrentTerminal(scriptPath, scriptName)
		}
		
		// Desktop environment: try common GUI terminals
		term := ""
		for _, candidate := range []string{"gnome-terminal", "konsole", "x-terminal-emulator", "xterm"} {
			if _, err := exec.LookPath(candidate); err == nil {
				term = candidate
				break
			}
		}
		if term == "" {
			// Fallback to current terminal execution if no GUI terminal found
			return ExecuteInCurrentTerminal(scriptPath, scriptName)
		}
		
		if strings.HasSuffix(scriptName, ".ps1") {
			cmd = exec.Command(term, "--", "bash", "-l", "-c", "clear; pwsh "+scriptPath+"; echo; read -p 'Press Enter to exit'")
		} else {
			cmd = exec.Command(term, "--", "bash", "-l", "-c", "clear; bash "+scriptPath+"; echo; read -p 'Press Enter to exit'")
		}
	}

	return cmd.Start()
}

// IsDesktopEnvironment checks if we're running in a desktop environment
func IsDesktopEnvironment() bool {
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

// ExecuteInCurrentTerminal executes a script in the current terminal using tmux or direct execution
func ExecuteInCurrentTerminal(scriptPath, scriptName string) error {
	// Check if we're already in tmux
	if os.Getenv("TMUX") != "" {
		// We're in tmux - create a new window
		var cmd *exec.Cmd
		if strings.HasSuffix(scriptName, ".ps1") {
			cmd = exec.Command("tmux", "new-window", "-n", scriptName, "bash", "-c", 
				fmt.Sprintf("clear; echo 'Running: %s'; if command -v bat &>/dev/null; then bat --theme=\"DarkNeon\" --style=numbers --color=always '%s'; elif command -v batcat &>/dev/null; then batcat --theme=\"DarkNeon\" --style=numbers --color=always '%s'; else cat '%s'; fi; echo; pwsh '%s'; echo; read -p 'Press Enter to close this window...'", scriptName, scriptPath, scriptPath, scriptPath, scriptPath))
		} else {
			cmd = exec.Command("tmux", "new-window", "-n", scriptName, "bash", "-c", 
				fmt.Sprintf("clear; echo 'Running: %s'; if command -v bat &>/dev/null; then bat --theme=\"DarkNeon\" --style=numbers --color=always '%s'; elif command -v batcat &>/dev/null; then batcat --theme=\"DarkNeon\" --style=numbers --color=always '%s'; else cat '%s'; fi; echo; bash '%s'; echo; read -p 'Press Enter to close this window...'", scriptName, scriptPath, scriptPath, scriptPath, scriptPath))
		}
		return cmd.Start()
	}
	
	// Check if tmux is available and start a new session
	if _, err := exec.LookPath("tmux"); err == nil {
		var cmd *exec.Cmd
		sessionName := fmt.Sprintf("go-pwr-%s", strings.ReplaceAll(scriptName, ".", "-"))
		if strings.HasSuffix(scriptName, ".ps1") {
			cmd = exec.Command("tmux", "new-session", "-d", "-s", sessionName, "bash", "-c", 
				fmt.Sprintf("clear; echo 'Running: %s'; echo 'Use Ctrl+B then D to detach, or exit to close'; if command -v bat &>/dev/null; then bat --theme=\"DarkNeon\" --style=numbers --color=always '%s'; elif command -v batcat &>/dev/null; then batcat --theme=\"DarkNeon\" --style=numbers --color=always '%s'; else cat '%s'; fi; echo; pwsh '%s'; echo; read -p 'Press Enter to close this session...'", scriptName, scriptPath, scriptPath, scriptPath, scriptPath))
		} else {
			cmd = exec.Command("tmux", "new-session", "-d", "-s", sessionName, "bash", "-c", 
				fmt.Sprintf("clear; echo 'Running: %s'; echo 'Use Ctrl+B then D to detach, or exit to close'; if command -v bat &>/dev/null; then bat --theme=\"DarkNeon\" --style=numbers --color=always '%s'; elif command -v batcat &>/dev/null; then batcat --theme=\"DarkNeon\" --style=numbers --color=always '%s'; else cat '%s'; fi; echo; bash '%s'; echo; read -p 'Press Enter to close this session...'", scriptName, scriptPath, scriptPath, scriptPath, scriptPath))
		}
		if err := cmd.Start(); err == nil {
			// Attach to the session
			attachCmd := exec.Command("tmux", "attach-session", "-t", sessionName)
			return attachCmd.Start()
		}
	}
	
	// Last resort: direct execution with warning
	fmt.Printf("\n=== Executing script directly (tmux not available) ===\n")
	fmt.Printf("Script: %s\n", scriptName)
	fmt.Printf("Path: %s\n", scriptPath)
	fmt.Printf("Install tmux for better experience: sudo apt install tmux\n")
	fmt.Printf("========================================================\n\n")
	
	var cmd *exec.Cmd
	if strings.HasSuffix(scriptName, ".ps1") {
		cmd = exec.Command("pwsh", scriptPath)
	} else {
		cmd = exec.Command("bash", scriptPath)
	}
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	err := cmd.Run()
	
	fmt.Printf("\n========================================================\n")
	fmt.Printf("Script execution completed. Press Enter to continue...")
	fmt.Scanln()
	
	return err
}
