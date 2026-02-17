#!/bin/sh
# Hopsule CLI Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/Hopsule/cli-tool/main/install.sh | bash
#
# This script detects your OS and architecture, downloads the latest
# Hopsule CLI binary from GitHub Releases, and installs it to /usr/local/bin.

set -e

REPO="Hopsule/cli-tool"
BINARY_NAME="hopsule"
INSTALL_DIR="/usr/local/bin"

# Colors (only if terminal supports it)
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    CYAN='\033[0;36m'
    NC='\033[0m'
else
    RED=''
    GREEN=''
    YELLOW=''
    CYAN=''
    NC=''
fi

info() {
    printf "${CYAN}==>${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}==>${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}WARNING:${NC} %s\n" "$1"
}

error() {
    printf "${RED}ERROR:${NC} %s\n" "$1" >&2
    exit 1
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) error "Unsupported operating system: $(uname -s)" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get latest version from GitHub API
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | \
            grep '"tag_name"' | head -1 | sed 's/.*"v\([^"]*\)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | \
            grep '"tag_name"' | head -1 | sed 's/.*"v\([^"]*\)".*/\1/'
    else
        error "Neither curl nor wget found. Please install one and try again."
    fi
}

# Download file
download() {
    local url="$1"
    local output="$2"

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "$output" "$url"
    elif command -v wget >/dev/null 2>&1; then
        wget -q -O "$output" "$url"
    fi
}

# Verify checksum
verify_checksum() {
    local file="$1"
    local checksums_url="$2"
    local expected_name="$3"

    local tmpcheck
    tmpcheck=$(mktemp)

    if download "$checksums_url" "$tmpcheck" 2>/dev/null; then
        local expected
        expected=$(grep "$expected_name" "$tmpcheck" | awk '{print $1}')
        if [ -n "$expected" ]; then
            local actual
            if command -v sha256sum >/dev/null 2>&1; then
                actual=$(sha256sum "$file" | awk '{print $1}')
            elif command -v shasum >/dev/null 2>&1; then
                actual=$(shasum -a 256 "$file" | awk '{print $1}')
            else
                warn "No SHA256 tool found, skipping checksum verification."
                rm -f "$tmpcheck"
                return 0
            fi

            if [ "$actual" = "$expected" ]; then
                info "Checksum verified."
            else
                rm -f "$tmpcheck"
                error "Checksum mismatch! Expected: ${expected}, Got: ${actual}"
            fi
        else
            warn "Could not find checksum for ${expected_name}, skipping verification."
        fi
    else
        warn "Could not download checksums, skipping verification."
    fi

    rm -f "$tmpcheck"
}

main() {
    local os
    local arch
    local version
    os=$(detect_os)
    arch=$(detect_arch)

    printf "\n"
    info "Hopsule CLI Installer"
    printf "\n"

    # Get version
    info "Fetching latest version..."
    version=$(get_latest_version)

    if [ -z "$version" ]; then
        error "Could not determine latest version. Check https://github.com/${REPO}/releases"
    fi

    info "Latest version: v${version}"
    info "Platform: ${os}/${arch}"

    # Build download URL
    local ext="tar.gz"
    if [ "$os" = "windows" ]; then
        ext="zip"
    fi
    local archive_name="${BINARY_NAME}-${os}-${arch}.${ext}"
    local url="https://github.com/${REPO}/releases/download/v${version}/${archive_name}"
    local checksums_url="https://github.com/${REPO}/releases/download/v${version}/checksums.txt"

    # Download
    info "Downloading ${archive_name}..."
    local tmpdir
    tmpdir=$(mktemp -d)
    local tmpfile="${tmpdir}/${archive_name}"

    if ! download "$url" "$tmpfile"; then
        rm -rf "$tmpdir"
        error "Download failed. URL: ${url}"
    fi

    # Verify checksum
    verify_checksum "$tmpfile" "$checksums_url" "$archive_name"

    # Extract
    info "Extracting..."
    if [ "$ext" = "tar.gz" ]; then
        tar xzf "$tmpfile" -C "$tmpdir"
    else
        unzip -q "$tmpfile" -d "$tmpdir"
    fi

    # Find binary
    local binary_path="${tmpdir}/${BINARY_NAME}"
    if [ ! -f "$binary_path" ]; then
        rm -rf "$tmpdir"
        error "Binary not found in archive."
    fi

    chmod +x "$binary_path"

    # Install
    info "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."

    if [ -w "$INSTALL_DIR" ]; then
        mv "$binary_path" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        info "Requesting sudo access to install to ${INSTALL_DIR}..."
        sudo mv "$binary_path" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    # Cleanup
    rm -rf "$tmpdir"

    # Verify installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        printf "\n"
        success "hopsule v${version} installed successfully!"
        printf "\n"
        info "Run 'hopsule' to launch the interactive dashboard."
        info "Run 'hopsule --help' for available commands."
        printf "\n"
    else
        printf "\n"
        success "Binary installed to ${INSTALL_DIR}/${BINARY_NAME}"
        warn "Make sure ${INSTALL_DIR} is in your PATH."
        printf "\n"
    fi
}

main "$@"
