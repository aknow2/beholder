#!/bin/sh
# Beholder Uninstallation Script for macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/main/scripts/uninstall.sh | sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
DATA_DIR="$HOME/.beholder"
MANIFEST_FILE="$DATA_DIR/install-manifest.txt"
FALLBACK_INSTALL_DIR="$HOME/.local/bin"

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

ask() {
    printf "${YELLOW}?${NC} %s " "$1"
}

# Read installation manifest
read_manifest() {
    if [ -f "$MANIFEST_FILE" ]; then
        info "Reading installation manifest..."
        # shellcheck disable=SC1090
        . "$MANIFEST_FILE"
        
        if [ -n "$INSTALLATION_PATH" ]; then
            info "Found installation: $INSTALLATION_PATH"
            return 0
        fi
    fi
    
    # Fallback: try to find beholder in PATH
    warn "Installation manifest not found, searching for beholder..."
    if command -v beholder >/dev/null 2>&1; then
        INSTALLATION_PATH=$(command -v beholder)
        info "Found beholder at: $INSTALLATION_PATH"
        USER_DATA_PATH="$DATA_DIR"
        return 0
    fi
    
    error "Could not locate beholder installation. Please specify the path manually or try 'which beholder'."
}

# Remove binary
remove_binary() {
    if [ -z "$INSTALLATION_PATH" ]; then
        warn "Installation path not specified, skipping binary removal"
        return
    fi
    
    if [ ! -f "$INSTALLATION_PATH" ]; then
        warn "Binary not found at $INSTALLATION_PATH, may already be removed"
        return
    fi
    
    info "Removing binary from $INSTALLATION_PATH..."
    rm -f "$INSTALLATION_PATH"
    success "Binary removed"
}

# Clean PATH from shell profile
clean_path() {
    info "Cleaning PATH from shell profiles..."
    
    # Detect shell profile files
    PROFILES=""
    [ -f "$HOME/.zshrc" ] && PROFILES="$PROFILES $HOME/.zshrc"
    [ -f "$HOME/.bashrc" ] && PROFILES="$PROFILES $HOME/.bashrc"
    [ -f "$HOME/.bash_profile" ] && PROFILES="$PROFILES $HOME/.bash_profile"
    [ -f "$HOME/.profile" ] && PROFILES="$PROFILES $HOME/.profile"
    
    if [ -z "$PROFILES" ]; then
        info "No shell profiles found to clean"
        return
    fi
    
    # Remove beholder PATH entries
    CLEANED=0
    for profile in $PROFILES; do
        if grep -q "# Added by beholder installer" "$profile" 2>/dev/null; then
            # Create backup
            cp "$profile" "$profile.beholder-backup"
            
            # Remove beholder entries (line with comment and next line)
            sed -i.tmp '/# Added by beholder installer/,+1d' "$profile"
            rm -f "$profile.tmp"
            
            success "Cleaned PATH from $(basename "$profile")"
            CLEANED=1
        fi
    done
    
    if [ $CLEANED -eq 0 ]; then
        info "No beholder PATH entries found in shell profiles"
    fi
}

# Handle user data
handle_user_data() {
    if [ ! -d "$USER_DATA_PATH" ]; then
        info "No user data found at $USER_DATA_PATH"
        return
    fi
    
    # Check if running interactively
    if [ -t 0 ]; then
        echo ""
        ask "Remove user data in $USER_DATA_PATH? (y/N): "
        read -r RESPONSE
        echo ""
    else
        # Non-interactive: preserve data by default
        RESPONSE="n"
        info "Running non-interactively, preserving user data"
    fi
    
    case "$RESPONSE" in
        [Yy]*)
            info "Removing user data..."
            rm -rf "$USER_DATA_PATH"
            success "User data removed"
            ;;
        *)
            info "User data preserved at $USER_DATA_PATH"
            warn "To remove manually later: rm -rf $USER_DATA_PATH"
            ;;
    esac
}

# Main uninstallation flow
main() {
    echo ""
    echo "Beholder Uninstaller"
    echo "===================="
    echo ""
    
    read_manifest
    remove_binary
    clean_path
    handle_user_data
    
    echo ""
    success "Uninstallation complete!"
    echo ""
    info "To verify removal: command -v beholder"
    if [ -d "$USER_DATA_PATH" ]; then
        info "User data preserved at: $USER_DATA_PATH"
    fi
    echo ""
}

main
