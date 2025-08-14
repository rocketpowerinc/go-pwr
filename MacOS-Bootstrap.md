# macOS Bootstrap Guide

Complete installation guide for **go-pwr** on macOS.

## 📋 Prerequisites

### Install All Dependencies (One Command)
```bash
# Install all required dependencies
brew install go git make bat glow gum gh jq
```

### Individual Installations (Alternative)
```bash
# Core tools
brew install go
brew install git
brew install make

# CLI tools  
brew install bat
brew install glow
brew install gum
brew install gh
brew install jq
```

### Install Homebrew (if not already installed)
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

## 🚀 Installation Methods

### Method 1: Go Install (Recommended)
```bash
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

## 🔧 Path Configuration

### Add Go bin to PATH (Zsh - default on modern macOS)
```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc
```

### For Bash users
```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc && source ~/.bashrc
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

## 🐞 Troubleshooting

### Common Issues
- **"go command not found"**: Restart your terminal after installing Go
- **PATH issues**: Use the path configuration steps above
- **Permission denied**: Make sure binary is executable (`chmod +x`)

### Terminal Recommendations
- **iTerm2** (highly recommended - better than default Terminal.app)
- **Warp**
- **Alacritty**

### Known Issues
- ⚠️ **macOS default Terminal** has display issues with borders/syntax highlighting
- **Solution**: Use iTerm2 for the best experience

## 🔒 Security Notes

macOS may show security warnings for downloaded binaries:
1. Go to **System Preferences > Security & Privacy**
2. Click **"Allow Anyway"** for go-pwr
3. Or use the Go install method which avoids this issue
