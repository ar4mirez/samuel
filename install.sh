#!/bin/sh
# Samuel CLI Installer
# Usage: curl -sSL https://raw.githubusercontent.com/ar4mirez/samuel/main/install.sh | sh
#
# This script detects your OS and architecture, downloads the appropriate
# binary from GitHub releases, and installs it to /usr/local/bin (or a
# user-specified location).

set -e

# Configuration
GITHUB_REPO="ar4mirez/samuel"
BINARY_NAME="samuel"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print functions
info() {
    printf "${GREEN}→${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}⚠${NC} %s\n" "$1"
}

error() {
    printf "${RED}✗${NC} %s\n" "$1" >&2
    exit 1
}

success() {
    printf "${GREEN}✓${NC} %s\n" "$1"
}

# Detect OS
detect_os() {
    OS="$(uname -s)"
    case "$OS" in
        Linux*)     OS=linux;;
        Darwin*)    OS=darwin;;
        MINGW*|MSYS*|CYGWIN*) OS=windows;;
        *)          error "Unsupported operating system: $OS";;
    esac
    echo "$OS"
}

# Detect architecture
detect_arch() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64)   ARCH=amd64;;
        aarch64|arm64)  ARCH=arm64;;
        *)              error "Unsupported architecture: $ARCH";;
    esac
    echo "$ARCH"
}

# Get latest version from GitHub API
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -sSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Download file
download() {
    URL="$1"
    OUTPUT="$2"

    if command -v curl >/dev/null 2>&1; then
        curl -sSL "$URL" -o "$OUTPUT"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$URL" -O "$OUTPUT"
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Verify checksum
verify_checksum() {
    FILE="$1"
    EXPECTED="$2"

    if command -v sha256sum >/dev/null 2>&1; then
        ACTUAL=$(sha256sum "$FILE" | cut -d ' ' -f 1)
    elif command -v shasum >/dev/null 2>&1; then
        ACTUAL=$(shasum -a 256 "$FILE" | cut -d ' ' -f 1)
    else
        warn "Unable to verify checksum (sha256sum/shasum not found)"
        return 0
    fi

    if [ "$ACTUAL" != "$EXPECTED" ]; then
        error "Checksum verification failed!\nExpected: $EXPECTED\nActual: $ACTUAL"
    fi
}

# Main installation
main() {
    echo ""
    echo " ███████╗ █████╗ ███╗   ███╗██╗   ██╗███████╗██╗     "
    echo " ██╔════╝██╔══██╗████╗ ████║██║   ██║██╔════╝██║     "
    echo " ███████╗███████║██╔████╔██║██║   ██║█████╗  ██║     "
    echo " ╚════██║██╔══██║██║╚██╔╝██║██║   ██║██╔══╝  ██║     "
    echo " ███████║██║  ██║██║ ╚═╝ ██║╚██████╔╝███████╗███████╗"
    echo " ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝ ╚═════╝ ╚══════╝╚══════╝"
    echo ""
    echo " Samuel - AI Coding Framework"
    echo ""

    # Detect platform
    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "Detected platform: ${OS}/${ARCH}"

    # Get version
    VERSION="${VERSION:-$(get_latest_version)}"
    if [ -z "$VERSION" ]; then
        error "Could not determine latest version. Please set VERSION environment variable."
    fi
    info "Installing version: ${VERSION}"

    # Determine file extension and archive format
    if [ "$OS" = "windows" ]; then
        EXT=".zip"
        BINARY_EXT=".exe"
    else
        EXT=".tar.gz"
        BINARY_EXT=""
    fi

    # Build download URL
    ARCHIVE_NAME="${BINARY_NAME}_${VERSION#v}_${OS}_${ARCH}${EXT}"
    DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"
    CHECKSUMS_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/checksums.txt"

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    # Download archive
    info "Downloading ${ARCHIVE_NAME}..."
    download "$DOWNLOAD_URL" "$TMP_DIR/$ARCHIVE_NAME"

    # Download and verify checksum
    info "Verifying checksum..."
    download "$CHECKSUMS_URL" "$TMP_DIR/checksums.txt"
    EXPECTED_CHECKSUM=$(grep "$ARCHIVE_NAME" "$TMP_DIR/checksums.txt" | cut -d ' ' -f 1)
    if [ -n "$EXPECTED_CHECKSUM" ]; then
        verify_checksum "$TMP_DIR/$ARCHIVE_NAME" "$EXPECTED_CHECKSUM"
        success "Checksum verified"
    else
        warn "Could not find checksum for $ARCHIVE_NAME"
    fi

    # Extract archive
    info "Extracting..."
    cd "$TMP_DIR"
    if [ "$OS" = "windows" ]; then
        unzip -q "$ARCHIVE_NAME"
    else
        tar -xzf "$ARCHIVE_NAME"
    fi

    # Install binary
    info "Installing to ${INSTALL_DIR}..."
    if [ ! -d "$INSTALL_DIR" ]; then
        mkdir -p "$INSTALL_DIR" 2>/dev/null || sudo mkdir -p "$INSTALL_DIR"
    fi

    if [ -w "$INSTALL_DIR" ]; then
        mv "${BINARY_NAME}${BINARY_EXT}" "${INSTALL_DIR}/${BINARY_NAME}${BINARY_EXT}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}${BINARY_EXT}"
    else
        sudo mv "${BINARY_NAME}${BINARY_EXT}" "${INSTALL_DIR}/${BINARY_NAME}${BINARY_EXT}"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}${BINARY_EXT}"
    fi

    # Verify installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        success "Installation complete!"
        echo ""
        "$BINARY_NAME" version
        echo ""
        echo "Run 'samuel init' to get started!"
    else
        success "Binary installed to ${INSTALL_DIR}/${BINARY_NAME}${BINARY_EXT}"
        warn "Make sure ${INSTALL_DIR} is in your PATH"
        echo ""
        echo "Add to your shell profile:"
        echo "  export PATH=\"\$PATH:${INSTALL_DIR}\""
    fi
}

# Run main
main "$@"
