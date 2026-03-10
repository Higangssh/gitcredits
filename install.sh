#!/bin/sh
set -e

REPO="Higangssh/gitcredits"
INSTALL_DIR="/usr/local/bin"
BINARY="gitcredits"

# Detect OS
OS_RAW=$(uname -s)
case "$OS_RAW" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
  *) echo "Error: unsupported OS: $OS_RAW"; exit 1 ;;
esac

# Detect architecture
ARCH_RAW=$(uname -m)
case "$ARCH_RAW" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Error: unsupported architecture: $ARCH_RAW"; exit 1 ;;
esac

# Windows + arm64 not supported
if [ "$OS" = "windows" ] && [ "$ARCH" = "arm64" ]; then
  echo "Error: Windows ARM64 is not supported"
  exit 1
fi

# Get latest version
echo "Fetching latest version..."
VERSION=$(curl -sL -o /dev/null -w '%{url_effective}' "https://github.com/$REPO/releases/latest" 2>/dev/null | sed 's|.*/tag/||')

if [ -z "$VERSION" ] || echo "$VERSION" | grep -q "releases$"; then
  echo "Error: could not determine latest version (no releases found)"
  exit 1
fi

# Strip 'v' prefix for filename
VERSION_NUM=$(echo "$VERSION" | sed 's/^v//')

# Build download URL
EXT="tar.gz"
if [ "$OS" = "windows" ]; then
  EXT="tar.gz"
fi

FILENAME="${BINARY}_${VERSION_NUM}_${OS}_${ARCH}.${EXT}"
URL="https://github.com/$REPO/releases/download/${VERSION}/${FILENAME}"

echo "Downloading $BINARY $VERSION for ${OS}/${ARCH}..."

# Create temp directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Download and extract
curl -sL "$URL" -o "$TMP_DIR/$FILENAME"

if [ $? -ne 0 ]; then
  echo "Error: download failed"
  exit 1
fi

cd "$TMP_DIR"
tar xzf "$FILENAME"

# Install
if [ -w "$INSTALL_DIR" ]; then
  mv "$BINARY" "$INSTALL_DIR/$BINARY"
else
  echo "Installing to $INSTALL_DIR (requires sudo)..."
  sudo mv "$BINARY" "$INSTALL_DIR/$BINARY"
fi

chmod +x "$INSTALL_DIR/$BINARY"

echo ""
echo "✓ $BINARY $VERSION installed to $INSTALL_DIR/$BINARY"
echo ""
echo "Run 'gitcredits' in any Git repo to see the credits!"
