#!/bin/bash

# Install script for ResponseWatch CLI
# Usage: curl -sSfL https://response-watch.web.app/install.sh | sh

set -e

APP_NAME="rwcli"
INSTALL_DIR="/usr/local/bin"

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
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    darwin|linux)
        ;;
    mingw*|msys*|cygwin*)
        OS="windows"
        APP_NAME="${APP_NAME}.exe"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

echo "Installing $APP_NAME for $OS/$ARCH..."

# Download URL from Firebase Hosting
BINARY_NAME="${APP_NAME}_${OS}_${ARCH}"
DOWNLOAD_URL="https://response-watch.web.app/cli/${BINARY_NAME}"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download binary
echo "Downloading from $DOWNLOAD_URL..."
if command -v curl &> /dev/null; then
    curl -fsSL -o "$TMP_DIR/$APP_NAME" "$DOWNLOAD_URL" || {
        echo "Error: Failed to download from $DOWNLOAD_URL"
        echo "Please check if the binary exists at that URL"
        exit 1
    }
elif command -v wget &> /dev/null; then
    wget -q -O "$TMP_DIR/$APP_NAME" "$DOWNLOAD_URL" || {
        echo "Error: Failed to download from $DOWNLOAD_URL"
        exit 1
    }
else
    echo "Error: curl or wget is required"
    exit 1
fi

# Make executable
chmod +x "$TMP_DIR/$APP_NAME"

# Check if binary works
if ! "$TMP_DIR/$APP_NAME" version &> /dev/null; then
    echo "Error: Downloaded binary is not valid"
    exit 1
fi

# Install
echo "Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$APP_NAME" "$INSTALL_DIR/"
else
    echo "Need sudo access to install to $INSTALL_DIR"
    sudo mv "$TMP_DIR/$APP_NAME" "$INSTALL_DIR/"
fi

echo ""
echo "✓ Installation complete!"
echo ""
echo "Run '$APP_NAME --help' to get started"
