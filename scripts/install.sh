#!/bin/sh
# Beholder Installation Script for macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.sh | sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO_OWNER="aknow2"
REPO_NAME="beholder"
INSTALL_DIR="$HOME/.local/bin"
DATA_DIR="$HOME/.beholder"
MANIFEST_FILE="$DATA_DIR/install-manifest.txt"

# Helper functions
info() {
    printf "${GREEN}→${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}⚠${NC} %s\n" "$1"
}

error() {
    printf "${RED}✗${NC} %s\n" "$1"
    exit 1
}

success() {
    printf "${GREEN}✓${NC} %s\n" "$1"
}

# Detect platform and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$OS" in
        darwin)
            PLATFORM="darwin"
            ;;
        *)
            error "Unsupported operating system: $OS (only macOS is supported)"
            ;;
    esac

    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $ARCH"
            ;;
    esac

    info "Detected platform: $PLATFORM-$ARCH"
}

# Check for existing installation
check_existing_installation() {
    EXISTING_VERSION=""
    
    if [ -f "$MANIFEST_FILE" ]; then
        # Read version from manifest
        EXISTING_VERSION=$(grep "INSTALLED_VERSION=" "$MANIFEST_FILE" | cut -d'=' -f2)
    elif command -v beholder >/dev/null 2>&1; then
        # Try to get version from binary
        EXISTING_VERSION=$(beholder --version 2>/dev/null | grep -o 'v[0-9.]*[0-9]' || echo "unknown")
    fi
    
    if [ -n "$EXISTING_VERSION" ]; then
        info "Existing installation found: $EXISTING_VERSION"
        UPGRADE_MODE=1
    else
        UPGRADE_MODE=0
    fi
}

# Get latest version from GitHub API
get_latest_version() {
    info "Fetching latest version..."
    
    VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" \
        | grep '"tag_name"' \
        | cut -d'"' -f4)
    
    if [ -z "$VERSION" ]; then
        error "Failed to fetch latest version"
    fi
    
    info "Latest version: $VERSION"
}

# Download binary with retry logic
download_binary() {
    BINARY_NAME="beholder-${VERSION}-${PLATFORM}-${ARCH}"
    DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/$BINARY_NAME"
    TEMP_FILE="/tmp/beholder-$$"
    
    info "Downloading from: $DOWNLOAD_URL"
    
    RETRIES=3
    for i in $(seq 1 $RETRIES); do
        if curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_FILE"; then
            success "Download complete"
            return 0
        fi
        
        if [ $i -lt $RETRIES ]; then
            warn "Download failed, retrying ($i/$RETRIES)..."
            sleep 2
        fi
    done
    
    error "Failed to download binary after $RETRIES attempts"
}

# Validate downloaded binary
validate_binary() {
    info "Validating binary..."
    
    if [ ! -f "$TEMP_FILE" ]; then
        error "Binary file not found: $TEMP_FILE"
    fi
    
    FILE_SIZE=$(stat -f%z "$TEMP_FILE" 2>/dev/null || stat -c%s "$TEMP_FILE" 2>/dev/null)
    if [ "$FILE_SIZE" -lt 1000000 ]; then
        error "Binary file is too small ($FILE_SIZE bytes), download may be corrupted"
    fi
    
    chmod +x "$TEMP_FILE"
    
    if ! "$TEMP_FILE" --version >/dev/null 2>&1; then
        error "Binary validation failed: cannot execute --version"
    fi
    
    success "Binary validated"
}

# Install binary to target directory
install_binary() {
    info "Installing to $INSTALL_DIR..."
    
    mkdir -p "$INSTALL_DIR"
    
    if [ -f "$INSTALL_DIR/beholder" ]; then
        if [ "$UPGRADE_MODE" = "1" ]; then
            info "Upgrading from $EXISTING_VERSION to $VERSION..."
        else
            warn "Existing installation found, replacing..."
        fi
        rm -f "$INSTALL_DIR/beholder"
    fi
    
    mv "$TEMP_FILE" "$INSTALL_DIR/beholder"
    chmod +x "$INSTALL_DIR/beholder"
    
    if [ "$UPGRADE_MODE" = "1" ]; then
        success "Upgrade complete"
    else
        success "Binary installed"
    fi
}

# Update PATH in shell profile
update_path() {
    info "Updating PATH..."
    
    # Detect shell profile file
    if [ -n "$ZSH_VERSION" ] || [ -f "$HOME/.zshrc" ]; then
        PROFILE_FILE="$HOME/.zshrc"
    elif [ -f "$HOME/.bashrc" ]; then
        PROFILE_FILE="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
        PROFILE_FILE="$HOME/.bash_profile"
    else
        PROFILE_FILE="$HOME/.profile"
    fi
    
    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        if ! grep -q "$INSTALL_DIR" "$PROFILE_FILE" 2>/dev/null; then
            echo "" >> "$PROFILE_FILE"
            echo "# Added by beholder installer" >> "$PROFILE_FILE"
            echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$PROFILE_FILE"
            success "Added $INSTALL_DIR to PATH in $PROFILE_FILE"
            warn "Please run: source $PROFILE_FILE"
        else
            info "PATH already configured in $PROFILE_FILE"
        fi
    else
        info "PATH already includes $INSTALL_DIR"
    fi
}

# Create installation manifest
create_manifest() {
    info "Creating installation manifest..."
    
    mkdir -p "$DATA_DIR"
    
    cat > "$MANIFEST_FILE" <<EOF
INSTALLED_VERSION=$VERSION
INSTALLATION_PATH=$INSTALL_DIR/beholder
INSTALLATION_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
USER_DATA_PATH=$DATA_DIR
PROFILE_FILE=$PROFILE_FILE
EOF
    
    success "Manifest created at $MANIFEST_FILE"
}

# Main installation flow
main() {
    echo ""
    echo "Beholder Installer"
    echo "=================="
    echo ""
    
    detect_platform
    check_existing_installation
    get_latest_version
    
    # Version comparison feedback
    if [ "$UPGRADE_MODE" = "1" ]; then
        if [ "$EXISTING_VERSION" = "$VERSION" ]; then
            info "Already at latest version ($VERSION), reinstalling..."
        elif [ "$EXISTING_VERSION" = "unknown" ]; then
            info "Upgrading to $VERSION..."
        else
            info "Upgrade available: $EXISTING_VERSION → $VERSION"
        fi
    fi
    
    download_binary
    validate_binary
    install_binary
    update_path
    create_manifest
    
    echo ""
    success "Installation complete!"
    echo ""
    info "Verify installation: beholder --version"
    info "Get started: beholder help"
    echo ""
}

main
