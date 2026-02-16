# Developer Quickstart: Beholder Installer Testing

**Feature**: 003-beholder-installer  
**Date**: 2026-02-02  
**Audience**: Developers testing installation/uninstallation workflows

## Overview

This guide explains how to test the installer and uninstaller scripts locally and in CI/CD, covering GitHub Actions workflow validation, cross-platform testing, and manual verification procedures.

---

## Prerequisites

- Go 1.24 installed
- Git repository cloned
- Access to test environments (or VMs) for macOS and Windows
- (Optional) [`act`](https://github.com/nektos/act) for local GitHub Actions testing

---

## 1. Testing GitHub Actions Workflow Locally

### Option A: Using `act` (Recommended for macOS)

Install `act`:
```bash
# macOS
brew install act

```

Run the release workflow locally:
```bash
cd /path/to/beholder

# Test the entire workflow
act -W .github/workflows/release.yml

# Test specific job
act -j build -W .github/workflows/release.yml

# Test with a specific tag event
act push -W .github/workflows/release.yml -e <(echo '{"ref":"refs/tags/v0.1.0"}')
```

**Limitations**:
- `act` runs in Docker containers, so artifact uploads may not work identically
- Matrix builds will execute sequentially (not in parallel)
- Useful for syntax validation and basic build testing

### Option B: Manual Cross-Compilation (macOS/Windows)

Simulate what GitHub Actions does by building manually:

```bash
cd /path/to/beholder

# Create output directory
mkdir -p dist

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/beholder-darwin-amd64 ./cmd/beholder

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/beholder-darwin-arm64 ./cmd/beholder

# Windows (amd64)
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/beholder-windows-amd64.exe ./cmd/beholder

# Verify builds
ls -lh dist/
```

**Verify binaries work**:
```bash
# macOS
./dist/beholder-darwin-amd64 --version
./dist/beholder-darwin-arm64 --version

# Windows (on Windows or via Wine)
./dist/beholder-windows-amd64.exe --version
```

Expected output: `beholder version v0.1.0` (or similar)

---

## 2. Testing Installation Scripts

### macOS Testing

**Setup clean environment** (create a new user or use a VM):
```bash
# Create test user
sudo dscl . -create /Users/testuser
sudo dscl . -create /Users/testuser UserShell /bin/zsh
sudo dscl . -create /Users/testuser UniqueID 1001
sudo dscl . -create /Users/testuser PrimaryGroupID 20
sudo dscl . -create /Users/testuser NFSHomeDirectory /Users/testuser
sudo createhomedir -c -u testuser

# Switch to test user
su - testuser
```

**Test installation**:
```bash
# Test with local script
bash /path/to/beholder/scripts/install.sh

# Verify
beholder --version
cat ~/.zshrc | grep beholder
```

**Test uninstall**:
```bash
bash /path/to/beholder/scripts/uninstall.sh

# Verify
command -v beholder && echo "FAIL" || echo "PASS"
```

### Windows Testing (PowerShell)

**Setup clean environment** (use a VM or Windows Sandbox):
```powershell
# Enable execution of local scripts (for testing only)
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy Bypass -Force
```

**Test installation**:
```powershell
# Test local script
cd C:\path\to\beholder
.\scripts\install.ps1

# Verify
beholder --version
$env:Path -split ';' | Select-String "beholder"
```

**Test uninstall**:
```powershell
.\scripts\uninstall.ps1

# Verify
Get-Command beholder -ErrorAction SilentlyContinue
if ($?) { Write-Host "FAIL: binary still present" } else { Write-Host "PASS: binary removed" }
```

---

## 3. Testing GitHub Actions in CI

### Triggering a Test Release

**Create a test tag**:
```bash
git checkout 003-beholder-installer
git tag v0.0.1-test
git push origin v0.0.1-test
```

**Monitor workflow**:
1. Go to https://github.com/aknow2/beholder/actions
2. Find the "Release" workflow run
3. Monitor each matrix job (darwin, windows)
4. Check artifact uploads

**Download and verify release assets**:
```bash
# List releases
gh release list

# Download artifacts
gh release download v0.0.1-test

# Verify binaries
chmod +x beholder-darwin-amd64
./beholder-darwin-amd64 --version

chmod +x beholder-darwin-arm64
./beholder-darwin-arm64 --version

# Windows (on Windows or Wine)
beholder-windows-amd64.exe --version
```

**Clean up test release**:
```bash
gh release delete v0.0.1-test --yes
git push origin --delete v0.0.1-test
git tag -d v0.0.1-test
```

---

## 4. End-to-End Installation Testing

### Test Scenario: Fresh Install → Use → Uninstall

**macOS**:
```bash
# 1. Fresh install
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/003-beholder-installer/scripts/install.sh | sh

# 2. Verify command works
beholder --version
beholder record --oneshot  # Creates config and database

# 3. Check user data was created
ls -la ~/.beholder/
cat ~/.beholder/config.yaml

# 4. Uninstall (preserve user data)
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/003-beholder-installer/scripts/uninstall.sh | sh
# When prompted, answer 'n' to preserve user data

# 5. Verify binary removed but data preserved
command -v beholder && echo "FAIL" || echo "PASS: binary removed"
[ -d ~/.beholder ] && echo "PASS: user data preserved" || echo "FAIL: user data removed"

# 6. Reinstall
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/003-beholder-installer/scripts/install.sh | sh

# 7. Verify data persisted across reinstall
beholder events  # Should show events from step 2

# 8. Final uninstall (remove data)
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/003-beholder-installer/scripts/uninstall.sh | sh
# When prompted, answer 'y' to remove user data

# 9. Verify everything removed
command -v beholder && echo "FAIL" || echo "PASS: binary removed"
[ -d ~/.beholder ] && echo "FAIL: user data still exists" || echo "PASS: user data removed"
```

**Windows (PowerShell)**:
```powershell
# 1. Fresh install
Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/003-beholder-installer/scripts/install.ps1 -UseBasicParsing | Invoke-Expression

# 2. Verify (restart terminal first!)
beholder --version
beholder record --oneshot

# 3. Check user data
Get-ChildItem $env:USERPROFILE\.beholder

# 4. Uninstall (preserve data)
Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/003-beholder-installer/scripts/uninstall.ps1 -UseBasicParsing | Invoke-Expression
# Answer 'n' when prompted

# 5-9. Repeat verification steps as above
```

---

## 5. Testing Upgrade Scenarios

### Test: Install v0.1.0 → Upgrade to v0.2.0

**Simulate old version**:
```bash
# Manually install "v0.1.0" (using older tag)
curl -fsSL https://github.com/aknow2/beholder/releases/download/v0.1.0/beholder-darwin-amd64 -o /tmp/beholder-old
chmod +x /tmp/beholder-old
mv /tmp/beholder-old ~/.local/bin/beholder

# Verify old version
beholder --version  # Should show v0.1.0

# Create installation manifest manually
mkdir -p ~/.beholder
cat > ~/.beholder/install-manifest.txt <<EOF
installed_version=v0.1.0
installation_path=$HOME/.local/bin/beholder
installation_date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
user_data_locations=$HOME/.beholder
EOF
```

**Run upgrade**:
```bash
# Run installer again (should detect existing installation)
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/003-beholder-installer/scripts/install.sh | sh

# When prompted about existing installation, answer 'y' to upgrade

# Verify new version
beholder --version  # Should show v0.2.0 (or latest)

# Verify manifest updated
cat ~/.beholder/install-manifest.txt | grep installed_version
```

---

## 6. Troubleshooting

### Common Issues

**Binary fails to execute after installation**:
```bash
# Check permissions
ls -l $(which beholder)
# Should be -rwxr-xr-x or similar

# Check if it's a valid binary
file $(which beholder)
# Should show "Mach-O" (macOS) or "PE32+" (Windows)

# Test directly
$(which beholder) --version
```

**PATH not updated after installation**:
```bash
# macOS
source ~/.bashrc  # or ~/.zshrc
echo $PATH | grep -o ".local/bin"

# Windows (restart terminal or)
$env:Path -split ';' | Select-String "beholder"
```

**Download fails with 404**:
```bash
# Verify release exists
curl -I https://github.com/aknow2/beholder/releases/latest
# Should return 302 redirect

# Check asset naming
curl -fsSL https://api.github.com/repos/aknow2/beholder/releases/latest | grep browser_download_url
```

**Uninstaller doesn't remove binary**:
```bash
# Check installation manifest
cat ~/.beholder/install-manifest.txt

# Manually verify binary location
which beholder

# If manifest is wrong, manually remove
rm -f $(which beholder)
```

---

## 7. CI/CD Integration Testing

### Automated Test Suite (Future)

```bash
#!/bin/bash
# tests/integration/installer_test.sh

set -e

echo "→ Testing fresh installation..."
./scripts/install.sh
beholder --version || exit 1

echo "→ Testing command execution..."
beholder record --oneshot || exit 1

echo "→ Testing user data creation..."
[ -f ~/.beholder/config.yaml ] || exit 1

echo "→ Testing uninstallation (preserve data)..."
echo "n" | ./scripts/uninstall.sh
command -v beholder && exit 1
[ -d ~/.beholder ] || exit 1

echo "→ Testing reinstallation..."
./scripts/install.sh
beholder --version || exit 1

echo "→ Testing uninstallation (remove data)..."
echo "y" | ./scripts/uninstall.sh
command -v beholder && exit 1
[ -d ~/.beholder ] && exit 1

echo "✓ All tests passed"
```

Run in GitHub Actions:
```yaml
test-installer:
  runs-on: ${{ matrix.os }}
  strategy:
    matrix:
      os: [ubuntu-latest, macos-latest, windows-latest]
  steps:
    - uses: actions/checkout@v4
    - name: Run installer tests
      run: bash tests/integration/installer_test.sh
```

---

## Next Steps

1. Execute Phase 1 testing on all target platforms
2. Document any edge cases or failures in GitHub Issues
3. Iterate on scripts based on test results
4. Proceed to Phase 2 (task breakdown) after testing validates design
