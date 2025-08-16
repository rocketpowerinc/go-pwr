#!/bin/bash
# macOS Bootstrap Script for go-pwr
# Automates installation of go-pwr and all dependencies on macOS

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Default options
INSTALL_DEPENDENCIES=false
SKIP_PATH_SETUP=false
SHOW_HELP=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --install-dependencies)
            INSTALL_DEPENDENCIES=true
            shift
            ;;
        --skip-path-setup)
            SKIP_PATH_SETUP=true
            shift
            ;;
        --help|-h)
            SHOW_HELP=true
            shift
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Show help
if [ "$SHOW_HELP" = true ]; then
    echo -e "${GREEN}macOS Bootstrap Script for go-pwr${NC}"
    echo ""
    echo -e "${YELLOW}USAGE:${NC}"
    echo "  ./bootstrap-macos.sh [--install-dependencies] [--skip-path-setup] [--help]"
    echo ""
    echo -e "${YELLOW}OPTIONS:${NC}"
    echo "  --install-dependencies  Install all dependencies (Homebrew, Go, etc.)"
    echo "  --skip-path-setup      Skip PATH configuration"
    echo "  --help, -h             Show this help message"
    echo ""
    echo -e "${YELLOW}EXAMPLES:${NC}"
    echo "  ./bootstrap-macos.sh                        # Install go-pwr only"
    echo "  ./bootstrap-macos.sh --install-dependencies # Install everything"
    echo ""
    exit 0
fi

echo -e "${GREEN}🚀 macOS Bootstrap Script for go-pwr${NC}"
echo -e "${GREEN}====================================${NC}"
echo ""

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    if [ "$INSTALL_DEPENDENCIES" = true ]; then
        echo -e "${YELLOW}📦 Installing Homebrew...${NC}"
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

        # Add Homebrew to PATH for Apple Silicon Macs
        if [[ $(uname -m) == "arm64" ]]; then
            echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
            eval "$(/opt/homebrew/bin/brew shellenv)"
        fi
        echo -e "${GREEN}✓ Homebrew installed successfully${NC}"
    else
        echo -e "${RED}❌ Homebrew not found. Install it first or use --install-dependencies${NC}"
        echo "Visit: https://brew.sh/"
        exit 1
    fi
fi

if [ "$INSTALL_DEPENDENCIES" = true ]; then
    echo -e "${YELLOW}📦 Installing Dependencies...${NC}"
    echo ""

    # Install Go
    echo -e "${CYAN}Installing Go...${NC}"
    if brew install go; then
        echo -e "${GREEN}✓ Go installed successfully${NC}"
    else
        echo -e "${RED}⚠️ Failed to install Go${NC}"
    fi

    # Install core dependencies
    echo -e "${CYAN}Installing core dependencies...${NC}"
    if brew install git gh jq make bat curl wget glow gum; then
        echo -e "${GREEN}✓ Core dependencies installed successfully${NC}"
    else
        echo -e "${RED}⚠️ Some core dependencies may have failed to install${NC}"
    fi

    # Install optional PowerShell
    echo -e "${CYAN}Installing PowerShell 7+ (optional)...${NC}"
    if brew install powershell; then
        echo -e "${GREEN}✓ PowerShell 7+ installed successfully${NC}"
    else
        echo -e "${YELLOW}⚠️ PowerShell 7+ installation failed (optional dependency)${NC}"
    fi

    echo ""
fi

# PATH Configuration
if [ "$SKIP_PATH_SETUP" != true ]; then
    echo -e "${YELLOW}🔧 Configuring PATH...${NC}"

    # Determine shell and profile file
    if [[ $SHELL == *"zsh"* ]]; then
        PROFILE_FILE="$HOME/.zshrc"
        SHELL_NAME="zsh"
    elif [[ $SHELL == *"bash"* ]]; then
        PROFILE_FILE="$HOME/.bashrc"
        SHELL_NAME="bash"
    else
        PROFILE_FILE="$HOME/.profile"
        SHELL_NAME="shell"
    fi

    # Add Go bin to PATH if not already present
    if ! grep -q 'export PATH="$HOME/go/bin:$PATH"' "$PROFILE_FILE" 2>/dev/null; then
        echo 'export PATH="$HOME/go/bin:$PATH"' >> "$PROFILE_FILE"
        echo -e "${GREEN}✓ Added Go bin to PATH in $PROFILE_FILE${NC}"

        # Apply to current session
        export PATH="$HOME/go/bin:$PATH"
    else
        echo -e "${GREEN}✓ Go bin already in PATH${NC}"
    fi
fi

# Install go-pwr
echo -e "${YELLOW}📥 Installing go-pwr...${NC}"
echo ""

# Try normal installation first
echo -e "${CYAN}Attempting normal installation...${NC}"
if go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest; then
    echo -e "${GREEN}✓ go-pwr installed successfully${NC}"
else
    echo -e "${YELLOW}⚠️ Normal installation failed, trying workaround...${NC}"

    # Try with workaround for checksum issues
    echo -e "${CYAN}Clearing module cache and retrying with direct proxy...${NC}"
    if go clean -modcache && GOPROXY=direct GOSUMDB=off go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest; then
        echo -e "${GREEN}✓ go-pwr installed successfully with workaround${NC}"
    else
        echo -e "${RED}❌ Failed to install go-pwr${NC}"
        echo ""
        echo -e "${YELLOW}Manual installation options:${NC}"
        echo -e "${WHITE}1. Download binary from: https://github.com/rocketpowerinc/go-pwr/releases/latest${NC}"
        echo -e "${WHITE}2. Build from source: git clone + make install${NC}"
        exit 1
    fi
fi

# Test installation
echo ""
echo -e "${YELLOW}🔍 Testing installation...${NC}"
if command -v go-pwr &> /dev/null; then
    VERSION_OUTPUT=$(go-pwr --version)
    echo -e "${GREEN}✓ Installation verified!${NC}"
    echo -e "${WHITE}$VERSION_OUTPUT${NC}"
else
    echo -e "${YELLOW}⚠️ Could not verify installation. Try running: go-pwr --version${NC}"
    echo -e "${YELLOW}You may need to restart your terminal or run: source $PROFILE_FILE${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Bootstrap complete!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo -e "${WHITE}1. Restart your terminal (or run: source $PROFILE_FILE)${NC}"
echo -e "${WHITE}2. Run: go-pwr${NC}"
echo ""
echo -e "${CYAN}For help: go-pwr --help${NC}"

# Known issues warning
echo ""
echo -e "${YELLOW}📝 Known Issues:${NC}"
echo -e "${WHITE}⚠️ macOS default Terminal has display issues with borders/syntax highlighting${NC}"
echo -e "${WHITE}💡 Solution: Use iTerm2 for the best experience${NC}"
