# Server Environment Improvements

## Changes Made

### 1. ✅ Enhanced Pane Switching
- **Added**: `Ctrl+H` (left pane) and `Ctrl+L` (right pane) - Vim-style navigation
- **Why**: These keys work better on servers where arrow key combinations are intercepted by SSH clients
- **Kept**: Original `Ctrl+Left/Right` and `Shift+Left/Right` for compatibility

### 2. ✅ Simplified to tmux-only
- **Removed**: screen support (as requested)
- **Streamlined**: Server script execution now uses tmux exclusively
- **Fallback**: Direct execution if tmux is not available
- **Benefits**: Simpler codebase, better user experience with consistent tmux workflow

### 3. ✅ Updated Documentation
- **README**: Updated to reflect tmux-only approach
- **Help text**: Added Ctrl+H/L key combinations
- **Server tips**: Focused on tmux installation and usage

## Key Combinations for Pane Switching (try these on your server):

1. `Ctrl+H` - Switch to left pane (list)
2. `Ctrl+L` - Switch to right pane (preview)
3. `Ctrl+Left` - Switch to left pane (if supported)
4. `Ctrl+Right` - Switch to right pane (if supported)
5. `Shift+Left` - Switch to left pane (alternative)
6. `Shift+Right` - Switch to right pane (alternative)

## Testing on Ubuntu Server

1. Transfer the new `build/go-pwr.exe` to your server (rename to `go-pwr`)
2. Make it executable: `chmod +x go-pwr`
3. Install tmux: `sudo apt install tmux`
4. Run `./go-pwr`
5. Test pane switching with `Ctrl+H` and `Ctrl+L`

The `Ctrl+H/L` combinations are based on Vim navigation and should work reliably across different SSH clients and terminal configurations.
