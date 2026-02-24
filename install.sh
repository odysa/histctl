#!/bin/sh
set -e

REPO="odysa/histctl"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr A-Z a-z)
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
esac

URL="https://github.com/${REPO}/releases/latest/download/histctl-${OS}-${ARCH}"

echo "Downloading histctl for ${OS}/${ARCH}..."
curl -fsSL "$URL" -o histctl
chmod +x histctl

if [ -w "$INSTALL_DIR" ]; then
  mv histctl "$INSTALL_DIR/histctl"
else
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo mv histctl "$INSTALL_DIR/histctl"
fi

echo "histctl installed to ${INSTALL_DIR}/histctl"
