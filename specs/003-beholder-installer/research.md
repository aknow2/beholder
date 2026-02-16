# Research: Beholder Installer & Uninstaller

**Feature**: 003-beholder-installer  
**Date**: 2026-02-02  
**Status**: Phase 0 Complete

## Overview

This document consolidates research findings for implementing a cross-platform installer/uninstaller system for the Beholder CLI tool. All NEEDS CLARIFICATION markers from the Technical Context have been resolved.

---

## 1. GitHub Actions Cross-Compilation for Go

### Decision: Matrix-Based Parallel Builds

**Chosen Approach**: Use GitHub Actions matrix strategy to build all platform/architecture combinations in parallel.

**Workflow Structure**:
```yaml
name: Release
on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - goos: darwin
            goarch: amd64
            output: beholder
          - goos: darwin
            goarch: arm64
            output: beholder
          - goos: windows
            goarch: amd64
            output: beholder.exe
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Build
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: |
          go build -ldflags="-s -w" -o dist/${{ matrix.output }} ./cmd/beholder
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: beholder-${{ matrix.goos }}-${{ matrix.goarch }}
          path: dist/${{ matrix.output }}
```

**Artifact Naming Convention**: `beholder-{os}-{arch}` (e.g., `beholder-darwin-arm64`, `beholder-windows-amd64`)

**Rationale**:
- Matrix builds execute in parallel, reducing total build time
- Each platform/arch combination is isolated, preventing cross-contamination
- Go's built-in cross-compilation support via GOOS/GOARCH requires no additional tooling
- Artifacts are uploaded individually, allowing partial success if one platform fails

**Alternatives Considered**:
- **Sequential builds**: Rejected due to longer total build time (5+ minutes per platform)
- **Self-hosted runners per platform**: Rejected due to maintenance overhead and complexity for a simple CLI tool
- **Docker-based builds**: Rejected as Go's native cross-compilation is simpler and faster

---

## 2. GitHub Artifacts API Usage from Shell Scripts

### Decision: GitHub Releases API + Direct Binary Downloads

**Chosen Approach**: Publish built binaries as GitHub Release assets (not workflow artifacts) and download via public HTTPS URLs.

**Download Pattern**:
```bash
# Latest release
VERSION=$(curl -fsSL https://api.github.com/repos/aknow2/beholder/releases/latest | grep '"tag_name"' | cut -d'"' -f4)
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then ARCH="amd64"; fi
if [ "$ARCH" = "aarch64" ]; then ARCH="arm64"; fi

URL="https://github.com/aknow2/beholder/releases/download/${VERSION}/beholder-${PLATFORM}-${ARCH}"
curl -fsSL "$URL" -o /tmp/beholder
chmod +x /tmp/beholder
```

**PowerShell Pattern**:
```powershell
$LatestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/aknow2/beholder/releases/latest"
$Version = $LatestRelease.tag_name
$Url = "https://github.com/aknow2/beholder/releases/download/$Version/beholder-windows-amd64.exe"
Invoke-WebRequest -Uri $Url -OutFile "$env:TEMP\beholder.exe"
```

**Rationale**:
- GitHub Releases provides permanent, publicly accessible URLs for binaries
- No authentication required for public repositories
- Releases are semantically versioned and immutable
- Simpler than workflow artifacts API (which requires authentication and has 90-day expiration)
- Users can manually download from the Releases page as a fallback

**Alternatives Considered**:
- **Workflow Artifacts API**: Rejected because it requires GitHub token authentication and artifacts expire after 90 days
- **External CDN (S3, CloudFront)**: Rejected to minimize dependencies and costs; GitHub Releases is sufficient for CLI distribution
- **Package registries (npm, Homebrew)**: Out of scope for P1/P2; can be added in future iterations

**Fallback Strategy**:
1. Attempt download from latest release
2. If 404, fall back to a hardcoded stable version URL
3. If download fails entirely, display error with manual download instructions

---

## 3. Platform-Specific Binary Installation Conventions

### Decision: User-Local Installation by Default

**Installation Paths**:

| Platform | User-Local Path | System-Wide Path (requires sudo/admin) |
|----------|----------------|----------------------------------------|
| macOS    | `~/.local/bin` | `/usr/local/bin` |
| Windows  | `%USERPROFILE%\.beholder\bin` | `C:\Program Files\beholder` |

**PATH Management**:

**macOS**:
```bash
# Add to appropriate profile file
PROFILE_FILE="$HOME/.profile"
if [ -f "$HOME/.bashrc" ]; then PROFILE_FILE="$HOME/.bashrc"; fi
if [ -f "$HOME/.zshrc" ]; then PROFILE_FILE="$HOME/.zshrc"; fi

if ! grep -q ".local/bin" "$PROFILE_FILE"; then
  echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$PROFILE_FILE"
  echo "Added ~/.local/bin to PATH in $PROFILE_FILE"
  echo "Run: source $PROFILE_FILE"
fi
```

**Windows (PowerShell)**:
```powershell
$BinPath = "$env:USERPROFILE\.beholder\bin"
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($CurrentPath -notlike "*$BinPath*") {
    [Environment]::SetEnvironmentVariable("Path", "$CurrentPath;$BinPath", "User")
    Write-Host "Added $BinPath to user PATH"
    Write-Host "Restart your terminal for changes to take effect"
}
```

**Rationale**:
- User-local installation requires no elevated privileges (lower barrier to entry)
- `~/.local/bin` is increasingly common on macOS and avoids elevated privileges
- Windows user PATH modification is reversible and does not require admin rights
- System-wide installation can be offered as an opt-in for advanced users

**Alternatives Considered**:
- **System-wide only**: Rejected due to privilege escalation friction
- **Symlinks to /usr/local/bin**: Rejected because it still requires sudo on macOS
- **Current directory installation**: Rejected as it doesn't solve PATH management

---

## 4. Installer Script Best Practices

### Decision: Idempotent, Fail-Safe Shell Scripts

**Idempotency Pattern**:
```bash
# Detect existing installation
if command -v beholder >/dev/null 2>&1; then
  EXISTING_VERSION=$(beholder --version 2>/dev/null | grep -o 'v[0-9.]\+' || echo "unknown")
  echo "Existing installation detected: $EXISTING_VERSION"
  read -p "Upgrade to latest version? (y/N): " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Installation cancelled"
    exit 0
  fi
fi
```

**Error Handling**:
```bash
set -e  # Exit on any error
trap 'echo "Installation failed. Check errors above."; exit 1' ERR

# Validate download
if [ ! -f /tmp/beholder ]; then
  echo "Error: Download failed"
  exit 1
fi

# Validate binary is executable
if ! /tmp/beholder --version >/dev/null 2>&1; then
  echo "Error: Downloaded binary is not valid"
  rm -f /tmp/beholder
  exit 1
fi
```

**User Feedback**:
```bash
echo "→ Downloading beholder $VERSION for $PLATFORM..."
echo "→ Installing to $INSTALL_DIR..."
echo "✓ Installation complete!"
echo ""
echo "Verify installation:"
echo "  beholder --version"
```

**Security Considerations**:
- Use HTTPS for all downloads (prevent MITM attacks)
- Verify binary executes successfully before moving to final location
- Use `set -e` and `trap` to prevent partial installations
- Avoid `eval` or dynamic code execution from downloaded content

**Rationale**:
- Idempotency allows safe re-runs (upgrade scenario)
- Early validation prevents broken installations
- Clear feedback reduces support burden
- Security-first approach builds user trust

**Alternatives Considered**:
- **Checksum verification**: Deferred to P3 (requires publishing checksums alongside binaries)
- **Digital signatures**: Out of scope for MVP; can add later with code signing certificates
- **Interactive prompts**: Minimized to preserve unattended install capability (FR-009)

---

## 5. Uninstaller Implementation Patterns

### Decision: Installation Manifest + Interactive User Data Handling

**Installation Manifest**:
```bash
# Create manifest during installation
cat > "$HOME/.beholder/install-manifest.txt" <<EOF
installed_version=$VERSION
installed_date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
binary_path=$INSTALL_DIR/beholder
EOF
```

**Uninstall Script**:
```bash
#!/bin/sh
set -e

# Read manifest
if [ -f "$HOME/.beholder/install-manifest.txt" ]; then
  . "$HOME/.beholder/install-manifest.txt"
else
  echo "Warning: No installation manifest found"
  binary_path=$(which beholder 2>/dev/null || echo "")
fi

# Remove binary
if [ -n "$binary_path" ] && [ -f "$binary_path" ]; then
  rm -f "$binary_path"
  echo "✓ Removed binary: $binary_path"
fi

# Handle user data
if [ -d "$HOME/.beholder" ]; then
  echo ""
  echo "User data directory found: $HOME/.beholder"
  echo "This includes your configuration and database."
  read -p "Remove user data? (y/N): " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -rf "$HOME/.beholder"
    echo "✓ Removed user data"
  else
    # Remove only installer artifacts
    rm -f "$HOME/.beholder/install-manifest.txt"
    echo "✓ User data preserved"
  fi
fi

# Remove from PATH (macOS)
for PROFILE in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile"; do
  if [ -f "$PROFILE" ]; then
    sed -i.bak '/\.local\/bin/d' "$PROFILE" 2>/dev/null || true
  fi
done

echo "✓ Uninstallation complete"
```

**Windows Uninstall (PowerShell)**:
```powershell
$BinPath = "$env:USERPROFILE\.beholder\bin"
$DataPath = "$env:USERPROFILE\.beholder"

# Remove binary
if (Test-Path "$BinPath\beholder.exe") {
    Remove-Item "$BinPath\beholder.exe" -Force
    Write-Host "✓ Removed binary"
}

# Handle user data
if (Test-Path $DataPath) {
    $response = Read-Host "Remove user data at $DataPath ? (y/N)"
    if ($response -eq 'y') {
        Remove-Item $DataPath -Recurse -Force
        Write-Host "✓ Removed user data"
    } else {
        Write-Host "✓ User data preserved"
    }
}

# Remove from PATH
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
$NewPath = $CurrentPath -replace [regex]::Escape(";$BinPath"), ""
[Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
Write-Host "✓ Removed from PATH"
```

**Rationale**:
- Installation manifest provides a single source of truth for what was installed
- Interactive prompt for user data respects FR-006 requirement
- Graceful handling when manifest is missing (searches for binary)
- Preserves user data by default (conservative, reversible)

**Alternatives Considered**:
- **Automatic user data removal**: Rejected as too aggressive; conflicts with Local-First principle
- **System uninstall integration** (macOS pkgutil, Windows Programs list): Deferred to future; requires packaging beyond simple script
- **Backup user data before removal**: Deferred to future enhancement

---

## Summary of Decisions

| Topic | Decision | Rationale |
|-------|----------|-----------|
| **CI/CD** | GitHub Actions matrix builds | Parallel execution, native Go support |
| **Artifact Distribution** | GitHub Releases (not workflow artifacts) | Public URLs, permanent retention, no auth required |
| **Installation Path** | User-local (`~/.local/bin` on Unix, `%USERPROFILE%\.beholder\bin` on Windows) | No privilege escalation, standard conventions |
| **PATH Management** | Modify user shell profile (Unix) or user PATH (Windows) | Persistent across sessions, no sudo required |
| **Installer Pattern** | Idempotent shell script with upgrade detection | Safe re-runs, clear error handling |
| **Uninstaller Pattern** | Manifest-based removal with interactive user data prompt | Clean uninstall, user control over data |

---

## Next Steps

Phase 0 research is complete. Proceed to:
1. **Phase 1: Design** - Generate `data-model.md`, `contracts/`, and `quickstart.md` based on these decisions
2. **Re-validate Constitution Check** - Confirm design still complies with architectural principles
3. **Phase 2: Tasks** - Break down implementation into concrete tasks aligned with user story priorities
