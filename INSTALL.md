# Beholder Installation Guide

This guide provides detailed installation instructions for Beholder on macOS and Windows.

## Table of Contents

- [Quick Install](#quick-install)
- [Manual Installation](#manual-installation)
- [Upgrading](#upgrading)
- [Uninstallation](#uninstallation)
- [Troubleshooting](#troubleshooting)

---

## Quick Install

### macOS

Open Terminal and run:

```bash
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.sh | sh
```

The installer will:
1. Detect your platform (Intel or Apple Silicon)
2. Download the latest version from GitHub Releases
3. Install to `~/.local/bin/beholder`
4. Add `~/.local/bin` to your PATH
5. Create installation manifest at `~/.beholder/install-manifest.txt`

**Note**: You may need to restart your terminal or run `source ~/.zshrc` (or `~/.bashrc`) for PATH changes to take effect.

### Windows

Open PowerShell and run:

```powershell
Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.ps1 -UseBasicParsing | Invoke-Expression
```

The installer will:
1. Download the latest Windows version
2. Install to `%USERPROFILE%\.beholder\bin\beholder.exe`
3. Add installation directory to user PATH
4. Create installation manifest at `%USERPROFILE%\.beholder\install-manifest.txt`

**Note**: You may need to restart your terminal for PATH changes to take effect.

---

## Manual Installation

If you prefer to install manually or the automated installer doesn't work:

### macOS

1. Visit [GitHub Releases](https://github.com/aknow2/beholder/releases/latest)
2. Download the appropriate binary:
   - **Intel Mac**: `beholder-v{VERSION}-darwin-amd64`
   - **Apple Silicon (M1/M2/M3)**: `beholder-v{VERSION}-darwin-arm64`
3. Make it executable and move to a directory in your PATH:

```bash
chmod +x beholder-v{VERSION}-darwin-{ARCH}
mkdir -p ~/.local/bin
mv beholder-v{VERSION}-darwin-{ARCH} ~/.local/bin/beholder
```

4. Add `~/.local/bin` to your PATH (if not already):

```bash
# For zsh (default on macOS)
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# For bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

5. Verify installation:

```bash
beholder --version
```

### Windows

1. Visit [GitHub Releases](https://github.com/aknow2/beholder/releases/latest)
2. Download `beholder-v{VERSION}-windows-amd64.exe`
3. Create installation directory and move binary:

```powershell
New-Item -ItemType Directory -Path "$env:USERPROFILE\.beholder\bin" -Force
Move-Item beholder-v{VERSION}-windows-amd64.exe "$env:USERPROFILE\.beholder\bin\beholder.exe"
```

4. Add to PATH:

```powershell
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
$newPath = "$currentPath;$env:USERPROFILE\.beholder\bin"
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
```

5. Restart your terminal and verify:

```powershell
beholder --version
```

---

## Upgrading

To upgrade to the latest version, simply run the installation script again:

### macOS

```bash
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.sh | sh
```

The installer will detect your existing installation and upgrade it automatically. Your user data in `~/.beholder/` will be preserved.

### Windows

```powershell
Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.ps1 -UseBasicParsing | Invoke-Expression
```

---

## Uninstallation

### macOS

Run the uninstall script:

```bash
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/main/scripts/uninstall.sh | sh
```

The uninstaller will:
1. Remove the binary from `~/.local/bin/beholder`
2. Clean PATH entries from shell profiles
3. Prompt whether to preserve or remove user data in `~/.beholder/`

**Non-interactive uninstall** (preserves user data):

```bash
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/main/scripts/uninstall.sh | sh < /dev/null
```

### Windows

Run the uninstall script:

```powershell
Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/main/scripts/uninstall.ps1 -UseBasicParsing | Invoke-Expression
```

The uninstaller will:
1. Remove the binary from `%USERPROFILE%\.beholder\bin\beholder.exe`
2. Clean PATH entries from user environment
3. Prompt whether to preserve or remove user data in `%USERPROFILE%\.beholder\`

### Manual Uninstallation

If you installed manually or the uninstall script doesn't work:

**macOS:**

```bash
# Remove binary
rm -f ~/.local/bin/beholder

# Remove user data (optional)
rm -rf ~/.beholder

# Clean PATH from your shell profile (~/.zshrc or ~/.bashrc)
# Remove lines containing "# Added by beholder installer"
```

**Windows:**

```powershell
# Remove binary
Remove-Item "$env:USERPROFILE\.beholder\bin\beholder.exe"

# Remove user data (optional)
Remove-Item -Recurse "$env:USERPROFILE\.beholder"

# Clean PATH manually via System Environment Variables GUI
# or use PowerShell to remove the PATH entry
```

---

## Troubleshooting

### "command not found: beholder" after installation

**Cause**: The installation directory is not in your PATH, or PATH changes haven't taken effect yet.

**Solution**:
- **macOS**: Restart your terminal or run `source ~/.zshrc` (or `~/.bashrc`)
- **Windows**: Restart your terminal or PowerShell window

If the issue persists, verify PATH:

```bash
# macOS
echo $PATH | grep -o ~/.local/bin

# Windows
$env:Path -split ';' | Select-String '.beholder'
```

### Download fails with 404 error

**Cause**: No releases have been published yet, or you're trying to access a non-existent version.

**Solution**:
- Check [GitHub Releases](https://github.com/aknow2/beholder/releases) to see if any releases exist
- Verify you're using the correct repository URL
- Try installing a specific version manually

### Permission denied errors

**Cause**: Installation directory requires elevated permissions.

**Solution**:
- The installer uses user-local directories (`~/.local/bin` on macOS, `%USERPROFILE%\.beholder\bin` on Windows) which should not require sudo/admin privileges
- If you encounter permission errors, check directory ownership:

```bash
# macOS
ls -la ~/.local/bin

# Windows
Get-Acl "$env:USERPROFILE\.beholder\bin"
```

### Binary validation fails

**Cause**: Downloaded binary is corrupted or incomplete.

**Solution**:
1. Try running the installer again (it will re-download)
2. Check your internet connection
3. Manually download from GitHub Releases and verify file size (should be 10-15 MB)

### Platform not supported error

**Cause**: The installer detected an unsupported operating system or architecture.

**Currently supported platforms**:
- macOS Intel (darwin-amd64)
- macOS Apple Silicon (darwin-arm64)
- Windows x64 (windows-amd64)

**Solution**:
- If you're on Linux, build from source using the development setup instructions in [README.md](README.md)
- For other platforms, please open an issue on GitHub

### Upgrade doesn't seem to work

**Cause**: The new binary might not be in your PATH, or you're running a cached version.

**Solution**:

```bash
# macOS - verify which binary is being used
which beholder
beholder --version

# Windows
Get-Command beholder
beholder --version
```

If the version doesn't match the latest release:
1. Restart your terminal
2. Check PATH configuration
3. Manually remove old binary and reinstall

### Uninstaller can't find installation

**Cause**: Installation manifest is missing or corrupted.

**Solution**:
- The uninstaller will fall back to searching PATH: `command -v beholder` (macOS) or `Get-Command beholder` (Windows)
- If that also fails, perform manual uninstallation (see above)

---

## Getting Help

If you encounter issues not covered in this guide:

1. Check existing [GitHub Issues](https://github.com/aknow2/beholder/issues)
2. Review the [main README](README.md) for configuration and usage
3. Open a new issue with:
   - Your operating system and version
   - Installation method used
   - Complete error messages
   - Output of `beholder --version` (if available)
