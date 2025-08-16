# Linux Bootstrap Guide

Complete installation guide for **go-pwr** on Linux.

## Install Flatpak and Snap Package Managers

Visit: https://flathub.org/setup
Visit: https://snapcraft.io/docs/installing-snapd

## 📋 Prerequisites

### Go Installation

```bash

# Ubuntu
sudo snap install go --classic
```

```bash
# Debian
sudo apt install golang-go
```

```bash
# Fedora
sudo dnf install golang
```

#### Add Go bin to PATH

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

## 📋 Dependencies

```bash
# Ubuntu/Debian - Core dependencies
sudo apt update && sudo apt install -y git gh jq make bat tmux curl wget glow gum

# Optional: PowerShell 7+ for .ps1 script support
wget -q https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb
sudo dpkg -i packages-microsoft-prod.deb
sudo apt update && sudo apt install -y powershell
```

```bash
# Fedora - Core dependencies
sudo dnf install -y git gh jq make bat tmux curl wget glow gum

# Optional: PowerShell 7+ for .ps1 script support
sudo dnf install -y powershell
```

## 🚀 Installation Methods

### Method 1: Go Install (Recommended)

```bash
# Install latest version
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest

# Or install specific version
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.7
```

**Note**: If you encounter checksum mismatch errors (common with Go module updates), use this workaround:

```bash
# Temporary workaround for checksum mismatch errors
GOPROXY=direct GOSUMDB=off go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest

# Alternative: Clear module cache and retry
go clean -modcache
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest
```

### Method 2: Download Binary

```bash
# Download latest release
wget https://github.com/rocketpowerinc/go-pwr/releases/latest/download/go-pwr-linux-amd64

# Make executable and install
chmod +x go-pwr-linux-amd64
sudo mv go-pwr-linux-amd64 /usr/local/bin/go-pwr
```

### Method 3: Build from Source

```bash
git clone https://github.com/rocketpowerinc/go-pwr.git
cd go-pwr

# If you encounter checksum mismatch errors, clear the module cache first:
go clean -modcache

# Then run make install
make install
```

**Note**: If you still encounter checksum mismatch errors when building from source, use this workaround:

```bash
# Clear module cache and use direct proxy
go clean -modcache
GOPROXY=direct GOSUMDB=off make install
```

## ⚡ Dev Alias (Advanced Users)

Add this alias to your shell profile for quick updates:

```bash
cat << 'EOF' >> ~/.bashrc
# Alias to launch latest go-pwr
function goo() {
    tmux new-session bash -c "
        cd \$HOME &&
        rm -rf go-pwr &&
        rm -f ~/go/bin/go-pwr &&
        git clone https://github.com/rocketpowerinc/go-pwr.git &&
        cd go-pwr &&
        go clean -modcache &&
        rm -f go.sum &&
        go mod tidy &&
        make install &&
        ~/go/bin/go-pwr;
        exec bash"
}
EOF
```

Then reload your shell:

```bash
source ~/.zshrc
```

## 🚀 Usage

### Basic Usage

```bash
go-pwr
```

Or with full path:

```bash
~/go/bin/go-pwr
```

**⚠️ Important**: The app will show a prominent warning recommending tmux usage on Linux.

### Recommended Usage (with tmux)

```bash
# Start in a new tmux session
tmux new-session go-pwr

# Or start tmux first, then run go-pwr
tmux
go-pwr
```

### Disable Warning (if needed)

- Add to bashrc so it persists

```bash
export GO_PWR_NO_TMUX_WARNING=1
go-pwr
```
