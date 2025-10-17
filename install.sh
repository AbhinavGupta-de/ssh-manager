#!/bin/bash

# SSH Manager Installation Script
# Works on macOS and Linux

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔑 SSH Manager Installation${NC}"
echo "================================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed!${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓${NC} Go is installed ($(go version))"

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
esac

echo -e "${GREEN}✓${NC} Detected platform: ${OS}/${ARCH}"

# Build the binary
echo -e "\n${YELLOW}Building SSH Manager...${NC}"
go build -o sshm main.go

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}✓${NC} Build successful"

# Determine installation directory
INSTALL_DIR=""
if [ "$OS" = "darwin" ]; then
    # macOS
    if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    elif [ -d "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    else
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
    fi
elif [ "$OS" = "linux" ]; then
    # Linux
    if [ -d "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
    else
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
    fi
fi

# Try to install with sudo if /usr/local/bin
echo -e "\n${YELLOW}Installing to ${INSTALL_DIR}...${NC}"

if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
    if sudo -n true 2>/dev/null; then
        sudo mv sshm "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/sshm"
    else
        echo -e "${YELLOW}This requires sudo permissions...${NC}"
        sudo mv sshm "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/sshm"
    fi
else
    mv sshm "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/sshm"
fi

echo -e "${GREEN}✓${NC} Installed to ${INSTALL_DIR}/sshm"

# Check if directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "\n${YELLOW}⚠️  ${INSTALL_DIR} is not in your PATH${NC}"
    
    SHELL_CONFIG=""
    if [ -n "$ZSH_VERSION" ]; then
        SHELL_CONFIG="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        if [ -f "$HOME/.bash_profile" ]; then
            SHELL_CONFIG="$HOME/.bash_profile"
        else
            SHELL_CONFIG="$HOME/.bashrc"
        fi
    fi
    
    if [ -n "$SHELL_CONFIG" ]; then
        echo -e "${BLUE}Adding ${INSTALL_DIR} to PATH in ${SHELL_CONFIG}...${NC}"
        echo "" >> "$SHELL_CONFIG"
        echo "# Added by SSH Manager" >> "$SHELL_CONFIG"
        echo "export PATH=\"\$PATH:${INSTALL_DIR}\"" >> "$SHELL_CONFIG"
        echo -e "${GREEN}✓${NC} Updated ${SHELL_CONFIG}"
        echo -e "${YELLOW}Please run: source ${SHELL_CONFIG}${NC}"
        echo -e "${YELLOW}Or restart your terminal${NC}"
    fi
else
    echo -e "${GREEN}✓${NC} ${INSTALL_DIR} is already in PATH"
fi

echo -e "\n${GREEN}✅ Installation complete!${NC}"
echo -e "\n${BLUE}Usage:${NC}"
echo "  sshm new              # Create a new profile"
echo "  sshm list             # List all profiles"
echo "  sshm switch <n>    # Switch profiles"
echo "  sshm current          # Show current profile"
echo ""
echo -e "${YELLOW}💡 Try: sshm new${NC}"