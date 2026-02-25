#!/bin/sh
set -e

REPO="odysa/histctl"
INSTALL_DIR="${HISTCTL_INSTALL_DIR:-$HOME/.local/bin}"

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

mkdir -p "$INSTALL_DIR"
mv histctl "$INSTALL_DIR/histctl"

echo "histctl installed to ${INSTALL_DIR}/histctl"

# Check if INSTALL_DIR is in PATH
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) echo "NOTE: $INSTALL_DIR is not in your PATH. Add it with:"
     echo "  export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
esac
