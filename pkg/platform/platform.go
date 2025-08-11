// Package platform provides platform-specific utilities for go-pwr.
package platform

import (
	"fmt"
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
		// Linux: try common terminals
		term := ""
		for _, candidate := range []string{"gnome-terminal", "konsole", "x-terminal-emulator"} {
			if _, err := exec.LookPath(candidate); err == nil {
				term = candidate
				break
			}
		}
		if term == "" {
			return fmt.Errorf("no supported terminal emulator found")
		}
		if strings.HasSuffix(scriptName, ".ps1") {
			cmd = exec.Command(term, "--", "bash", "-l", "-c", "clear; pwsh "+scriptPath+"; echo; read -p 'Press Enter to exit'")
		} else {
			cmd = exec.Command(term, "--", "bash", "-l", "-c", "clear; bash "+scriptPath+"; echo; read -p 'Press Enter to exit'")
		}
	}

	return cmd.Start()
}
