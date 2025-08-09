# Go-PWR Application

## ✨ Overview

**`go-pwr`** is a cross-platform launcher for your personal automation scripts. Built with Go and powered by [Charm’s](https://github.com/charmbracelet) [Bubble Tea framework](https://github.com/charmbracelet/bubbletea), it delivers a sleek, interactive TUI for browsing, previewing, and running  bash and powershell scripts across Windows, macOS, Linux, and server environments. On first run, it automatically clones my “scriptbin” repository into a local  directory along with the `go-pwr` TUI, making script access and management simple and seamless.

---

## 📥 Installation

First, install Go:

- **Windows:**

    `winget install -e --id GoLang.Go`

- **macOS:**

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

## 🛠️ Development

This project follows Go best practices with a clean, modular architecture:

```
go-pwr/
├── cmd/go-pwr/        # Application entry point
├── internal/          # Private application code
│   ├── app/           # Core application logic
│   ├── config/        # Configuration management
│   ├── git/           # Git operations
│   ├── scripts/       # Script discovery and management
│   └── ui/            # User interface components
├── pkg/platform/      # Platform-specific utilities
└── Makefile          # Build automation
```

### Build Commands

- `make build` - Build the application
- `make install` - Install to GOPATH/bin
- `make dev` - Build and run in development mode
- `make test` - Run tests
- `make clean` - Clean build artifacts

See `ARCHITECTURE.md` for detailed documentation.

---

## ⚡️ Dev Aliases

### Windows (PowerShell)

```powershell
$gooAlias = @"
function goo {
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

### Linux (Bash)

```bash
cat << 'EOF' >> ~/.bashrc
# Alias to reinstall go-pwr
alias goo='
  rm -rf go-pwr && \\
  rm -f ~/go/bin/go-pwr && \\
  git clone https://github.com/rocketpowerinc/go-pwr.git && \\
  cd go-pwr && \\
  make install && \\
  ~/go/bin/go-pwr'
EOF

```

### macOS (Zsh)

```bash
cat << 'EOF' >> ~/.zshrc
# Alias to reinstall go-pwr
alias goo='
  rm -rf go-pwr && \\
  rm -f ~/go/bin/go-pwr && \\
  git clone <https://github.com/rocketpowerinc/go-pwr.git> && \\
  cd go-pwr && \\
  go install && \\
  ~/go/bin/go-pwr'
EOF

```

Then reload your shell:

- `. $PROFILE`
- `source ~/.bashrc`
- `source ~/.zshrc`

---

## 🐞 Known Issues / Bugs

- ⚠️ macOS default Terminal has issues with borders/syntax in `gum`; use **iTerm2** instead