#!/bin/bash
#
# News4Coder 一键安装脚本
# Usage: curl -sSL https://get.news4coder.dev | bash
#        bash install.sh

set -e

REPO="news4coder/news4coder"
BINARY_NAME="news4coder"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        echo "Please build from source: go build -o news4coder"
        exit 1
        ;;
esac

case "$OS" in
    linux)
        OS="linux"
        ;;
    darwin)
        OS="darwin"
        ;;
    mingw*|cygwin*|msys*)
        OS="windows"
        BINARY_NAME="${BINARY_NAME}.exe"
        ;;
    *)
        echo -e "${RED}Unsupported OS: $OS${NC}"
        echo "Please build from source: go build -o news4coder"
        exit 1
        ;;
esac

DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
fi

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo -e "${GREEN}Installing News4Coder...${NC}"
echo "OS: $OS"
echo "Arch: $ARCH"

# Check if curl or wget is available
if command -v curl >/dev/null 2>&1; then
    FETCH="curl -fsSL"
elif command -v wget >/dev/null 2>&1; then
    FETCH="wget -qO-"
else
    echo -e "${RED}curl or wget is required${NC}"
    exit 1
fi

# Download binary
echo "Downloading from GitHub releases..."
$FETCH "$DOWNLOAD_URL" > "$TMP_DIR/$BINARY_NAME" || {
    echo -e "${YELLOW}Binary download failed. Falling back to building from source...${NC}"
    if ! command -v go >/dev/null 2>&1; then
        echo -e "${RED}Go is not installed. Please install Go 1.25+ and run:${NC}"
        echo "  go install github.com/${REPO}@latest"
        exit 1
    fi
    go install "github.com/${REPO}@latest"
    echo -e "${GREEN}News4Coder installed via go install!${NC}"
    echo "Binary location: $(which $BINARY_NAME 2>/dev/null || echo '$GOPATH/bin/$BINARY_NAME')"
    exit 0
}

chmod +x "$TMP_DIR/$BINARY_NAME"

# Install to target directory
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Installing to $INSTALL_DIR (requires sudo)..."
    sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

echo -e "${GREEN}✓ News4Coder installed successfully!${NC}"
echo "Location: $INSTALL_DIR/$BINARY_NAME"
echo ""
echo "Get started:"
echo "  $BINARY_NAME --help"
echo "  $BINARY_NAME sync"
