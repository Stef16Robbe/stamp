#!/bin/sh
set -e

# Stamp installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/stef16robbe/stamp/main/install.sh | sh

REPO="stef16robbe/stamp"
BINARY="stamp"
INSTALL_DIR="${HOME}/.local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    printf "${BLUE}==>${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}==>${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}==>${NC} %s\n" "$1"
}

error() {
    printf "${RED}==>${NC} %s\n" "$1"
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
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get latest version from GitHub
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" |
        grep '"tag_name":' |
        sed -E 's/.*"([^"]+)".*/\1/'
}

main() {
    echo ""
    printf "${GREEN}"
    echo "  _____ _                        "
    echo " / ____| |                       "
    echo "| (___ | |_ __ _ _ __ ___  _ __  "
    echo " \\___ \\| __/ _\` | '_ \` _ \\| '_ \\ "
    echo " ____) | || (_| | | | | | | |_) |"
    echo "|_____/ \\__\\__,_|_| |_| |_| .__/ "
    echo "                          | |    "
    echo "                          |_|    "
    printf "${NC}"
    echo ""

    OS=$(detect_os)
    ARCH=$(detect_arch)

    info "Detected OS: ${OS}"
    info "Detected arch: ${ARCH}"

    info "Fetching latest version..."
    VERSION=$(get_latest_version)

    if [ -z "$VERSION" ]; then
        error "Could not determine latest version"
    fi

    info "Latest version: ${VERSION}"

    # Construct download URL
    FILENAME="${BINARY}_${OS}_${ARCH}"
    if [ "$OS" = "windows" ]; then
        FILENAME="${FILENAME}.exe"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

    info "Downloading ${DOWNLOAD_URL}..."

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf ${TMP_DIR}" EXIT

    # Download binary
    if ! curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${BINARY}"; then
        error "Download failed. Please check if the release exists for your platform."
    fi

    # Make executable
    chmod +x "${TMP_DIR}/${BINARY}"

    # Create install directory if needed
    if [ ! -d "$INSTALL_DIR" ]; then
        info "Creating ${INSTALL_DIR}..."
        mkdir -p "$INSTALL_DIR"
    fi

    # Install
    info "Installing to ${INSTALL_DIR}/${BINARY}..."
    mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"

    success "Successfully installed stamp ${VERSION}"
    echo ""

    # Check if install dir is in PATH
    case ":${PATH}:" in
        *":${INSTALL_DIR}:"*) ;;
        *)
            warn "${INSTALL_DIR} is not in your PATH"
            echo ""
            echo "  Add this to your shell profile (.bashrc, .zshrc, etc.):"
            echo ""
            echo "    export PATH=\"\${HOME}/.local/bin:\${PATH}\""
            echo ""
            ;;
    esac

    echo "  Run 'stamp --help' to get started"
    echo ""
}

main
