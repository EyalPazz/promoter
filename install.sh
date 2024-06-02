#!/bin/bash

# Define the download URL
URL="https://github.com/EyalPazz/promoter/releases/download/v0.1.4/promoter"

# Define the install path
INSTALL_PATH="/usr/local/bin/promoter"

# Download the binary
curl -L $URL -o promoter

# Make it executable
chmod +x promoter

# Move it to the install path
sudo mv promoter $INSTALL_PATH

echo "promoter installed successfully!"

