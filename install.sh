#!/bin/bash

# Detect the platform and architecture
OS=$(uname -s)
ARCH=$(uname -m)
VERSION=v0.3.0

# Map platform and architecture to the correct binary name
case "$OS" in
  Linux)
    PLATFORM="linux"
    ;;
  Darwin)
    PLATFORM="darwin"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  arm64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Define the download URL based on platform and architecture
URL="https://github.com/EyalPazz/promoter/releases/download/${VERSION}/promoter_${PLATFORM}_${ARCH}"

# Define the install path
INSTALL_PATH="/usr/local/bin/promoter"

# Download the binary
curl -L $URL -o promoter

# Make it executable
chmod +x promoter

# Move it to the install path
sudo mv promoter $INSTALL_PATH

echo "promoter installed successfully!"
