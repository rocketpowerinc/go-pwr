# macOS Bootstrap Guide

Complete installation guide for **go-pwr** on macOS.

## Install Homebrew Package Manager

Visit: https://brew.sh/

- Run install

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

- Add Paths to zsh

```
  echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> /Users/rocket/.zprofile
eval "$(/opt/homebrew/bin/brew shellenv)"
```

## 📋 Prerequisites

### Go Installation

```bash
brew install go
```

## 🔧 Path Configuration

### Add Go bin to PATH (Zsh - default on modern macOS)

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc
```

## 📋 Dependencies

```bash
brew install git gh jq make bat curl wget glow gum
```

## 🚀 Installation Methods

### Method 1: Go Install (Recommended)

```bash
# Install latest version
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest

# Or install specific version (if @latest doesn't show newest)
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.4
```

**Note**: If `@latest` installs an older version, use the specific version or clear the module cache:
```bash
go clean -modcache
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest
```

### Method 2: Download Binary

1. Visit [Releases](https://github.com/rocketpowerinc/go-pwr/releases/latest)
2. Download:
   - **Intel Macs**: `go-pwr-darwin-amd64`
   - **Apple Silicon**: `go-pwr-darwin-arm64`
3. Make executable: `chmod +x go-pwr-darwin-*`
4. Move to PATH: `sudo mv go-pwr-darwin-* /usr/local/bin/go-pwr`

### Method 3: Build from Source

```bash
git clone https://github.com/rocketpowerinc/go-pwr.git
cd go-pwr
make install
```

## 🚀 Usage

After installation:

```bash
go-pwr
```

Or with full path:

```bash
~/go/bin/go-pwr
```

## ⚡ Dev Alias (Advanced Users)

Add this alias to your shell profile for quick updates:

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

Then reload your shell:

```bash
source ~/.zshrc
```

### Known Issues

- ⚠️ **macOS default Terminal** has display issues with borders/syntax highlighting
- **Solution**: Use iTerm2 for the best experience

## 🔒 Security Notes

macOS may show security warnings for downloaded binaries:

1. Go to **System Preferences > Security & Privacy**
2. Click **"Allow Anyway"** for go-pwr
3. Or use the Go install method which avoids this issue
