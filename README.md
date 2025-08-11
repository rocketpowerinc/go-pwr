# Go-PWR Application

## ✨ Overview

**`go-pwr`** is a cross-platform launcher for your personal automation scripts. Built with Go and powered by [Charm's](https://github.com/charmbracelet) [Bubble Tea framework](https://github.com/charmbracelet/bubbletea), it delivers a sleek, interactive TUI for browsing, previewing, and running bash and powershell scripts across Windows, macOS, Linux, and server environments.

Features beautiful syntax highlighting for script previews (when `bat` is installed), tag-based script search functionality. It automatically clones the "scriptbin" repository to `$HOME/Downloads/Temp/scriptbin`, providing a centralized location for script access and management that's easily accessible and always up to date with the latest scripts.

---

## 📥 Installation

First, install Go:

- **Windows:**

    `winget install -e --id GoLang.Go`

- **MacOS:**

    `brew install go`

- **Ubuntu:**

    `sudo snap install go --classic`

    Add to PATH if needed:

    `export PATH=/snap/bin:$PATH`


Then install `go-pwr`:

```bash
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest
```

Or build from source:

```bash
git clone https://github.com/rocketpowerinc/go-pwr.git
cd go-pwr
make install
```

---

## 🚀 Usage

After installation, you can run the app from any directory:

```bash
go-pwr

```

Or directly:

```bash
~/go/bin/go-pwr

```

To ensure it's always in your path:

- **Linux (Bash):**

    ```bash
    echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc && source ~/.bashrc

    ```

- **macOS (Zsh):**

    ```bash
    echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc

    ```
---


## ⚡️ Dev Aliases

### Windows (PowerShell)

```powershell
$gooAlias = @"
# Alias to launch latest go-pwr
function goo {
    Set-Location $env:USERPROFILE
    Remove-Item -Recurse -Force go-pwr -ErrorAction SilentlyContinue
    Remove-Item -Force `"$HOME\\go\\bin\\go-pwr.exe`" -ErrorAction SilentlyContinue
    git clone https://github.com/rocketpowerinc/go-pwr.git
    Set-Location go-pwr
    make install
    & "$env:USERPROFILE\\go\\bin\\go-pwr.exe"
}
"@
$gooAlias | Out-File -Append -Encoding UTF8 $PROFILE

```


### MacOS (Zsh)

```bash
cat << 'EOF' >> ~/.zshrc
# Alias to launch latest go-pwr
alias goo='
    cd $HOME
    rm -rf go-pwr && \
    rm -f ~/go/bin/go-pwr && \
    git clone https://github.com/rocketpowerinc/go-pwr.git && \
    cd go-pwr && \
    make install && \
    ~/go/bin/go-pwr'
EOF


```

### Linux (Bash)

```bash
cat << 'EOF' >> ~/.bashrc
# Alias to launch latest go-pwr
function goo() {
  cd $HOME
  rm -rf go-pwr &&
  rm -f ~/go/bin/go-pwr &&
  git clone https://github.com/rocketpowerinc/go-pwr.git &&
  cd go-pwr &&
  make install &&
  ~/go/bin/go-pwr
}
EOF


```

Then reload your shell:

- `. $PROFILE`
- `source ~/.bashrc`
- `source ~/.zshrc`

---
## 🐞 Known Issues / Bugs

- ⚠️ macOS default Terminal has issues with borders/syntax in `gum`; use **iTerm2** instead

---

## 🖥️ Server Environment Support

**`go-pwr`** works seamlessly on both desktop and server environments:

**Desktop Environments:**
- Opens scripts in new terminal windows (Windows Terminal, iTerm2, gnome-terminal, etc.)
- Full GUI experience with syntax highlighting

**Server Environments (Headless/SSH):**
- Automatically detects server environments (no DISPLAY variable)
- Uses `tmux` sessions to run scripts in separate windows
- Falls back to direct execution if tmux is not available
- Perfect for remote server management and automation

**Server Usage Tips:**
- Install `tmux` for the best experience: `sudo apt install tmux`
- Scripts run in named sessions: `tmux list-sessions` to see active sessions
- Use `Ctrl+B, D` to detach from running scripts
- tmux allows scripts to continue running even if you disconnect

---

### 🎨 Syntax Highlighting

`go-pwr` automatically provides beautiful syntax highlighting in script previews using [`bat`](https://github.com/sharkdp/bat) when available. **If `bat` is not installed, the application gracefully falls back to plain text previews with helpful installation instructions.**

**Optional: Install `bat` for enhanced syntax highlighting:**

- **Windows:**
  ```bash
  winget install sharkdp.bat
  ```

- **macOS:**
  ```bash
  brew install bat
  ```

- **Ubuntu/Debian:**
  ```bash
  sudo apt install bat
  ```
  *Note: On Ubuntu, the command is known as `batcat` but `go-pwr` automatically detects this.*

---
### 🔍 Tag-Based Search

**`go-pwr`** includes powerful tag-based search functionality to help you quickly find the right scripts for your needs. Scripts can be tagged with:

**Usage:**
- Press `Ctrl+F` to activate search mode
- Type multiple tags separated by spaces (e.g., `bash linux ubuntu`)
- Search results update in real-time as you type
- Press `Enter` to apply search or `Escape` to cancel
- Press `Escape` again to clear search and show all scripts


## 🏷️ Tagging Your Scripts

To make your scripts searchable, add tags at the top of your script files using this format:

### Bash/Shell Scripts (.sh)
```bash
#!/usr/bin/env bash
set -euo pipefail

#*Tags:
# Languages: bash zsh
# Platforms: Linux Mac WSL
# Distros: Ubuntu Debian
# Categories: utility development
# PackageManagers: apt brew

# Your script content here...
```

### PowerShell Scripts (.ps1)
```powershell
#! ADMIN NOT REQUIRED
#! Description: Your script description

#*Tags:
# Languages: pwsh powershell
# Platforms: Windows
# Categories: utility demo
# PackageManagers: winget chocolatey

# Your script content here...
```

**Tagging Guidelines:**
- Start the tags section with `#*Tags:` on its own line
- Each category starts with `# CategoryName:` followed by space-separated tags
- Use lowercase for consistency (parser handles case-insensitivity)
- Common categories: `Languages`, `Platforms`, `Distros`, `Categories`, `PackageManagers`, `DesktopEnvironments`, `Architectures`
- Add as many or as few tags as appropriate for your script


## 🔄 Recursive vs Directory Mode

**`go-pwr`** supports two viewing modes:

- **Directory Mode** (default): Browse scripts folder by folder, just like a file manager
- **Recursive Mode**: Search and view ALL scripts from all subdirectories at once

**Toggle between modes:**
- Press `Ctrl+R` to switch between Directory and Recursive modes
- In **Recursive Mode**:
  - See all scripts from every subdirectory in one list
  - Perfect for searching across your entire script collection
  - Directory navigation is disabled (no need to browse folders)
  - Scripts show their relative path from the root
- In **Directory Mode**:
  - Browse one folder at a time
  - Use arrow keys to navigate into/out of directories
  - Traditional file manager experience

---

## ⌨️ Keyboard Shortcuts

**Navigation:**
- `Tab` / `Shift+Tab` - Switch between tabs (Scripts, Options, About)
- `Ctrl+Tab` - Alternative tab switching (useful when Tab key is intercepted)
- `↑` / `↓` - Navigate through lists
- `←` / `→` - Navigate directories (Directory Mode only)
- `Enter` - Run script or enter directory

**Pane Switching (when preview is available):**
- `Ctrl+←` / `Ctrl+→` - Switch between left panel and preview pane
- `Shift+←` / `Shift+→` - Alternative pane switching
- `Ctrl+H` / `Ctrl+L` - Vim-style pane switching (works better on servers)

**Search and Modes:**
- `Ctrl+F` or `/` - Activate tag-based search
- `Ctrl+R` - Toggle between Recursive and Directory viewing modes
- `Escape` - Cancel search or clear search terms

**Preview Navigation:**
- `Page Up` / `Page Down` - Scroll through script preview
- Arrow keys when preview pane is focused

**General:**
- `q` or `Ctrl+C` - Quit application

---