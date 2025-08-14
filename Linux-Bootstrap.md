# Server Bootstrap Guide

# 📋 Prerequisites

### Go Installation
```bash

# Ubuntu
sudo snap install go --classic

# Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

```
#### Add Go bin to PATH
```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

# 📋 Dependencies

### Ubuntu/Debian
```bash
sudo apt update && sudo apt install -y git gh jq make bat tmux curl glow gum
```

### Fedora
```bash
sudo dnf install -y git gh jq make bat tmux curl glow gum
```

#### Why These Tools Are Required
- **`bat`**: Enhanced syntax highlighting for script previews
- **`make`**: Build system for compiling from source
- **`glow`**: Beautiful markdown rendering in terminal
- **`gum`**: Interactive UI components and prompts
- **`git`**: Version control and repository management
- **`gh`**: GitHub CLI for repository operations
- **`jq`**: JSON processing for API responses
- **`tmux`**: Session management and persistence (Linux only)



## 🚀 Installation Methods

### Method 1: Go Install (Recommended)
```bash
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
make install
```


## 🚀 Usage
- From terminal `go-pwr`

