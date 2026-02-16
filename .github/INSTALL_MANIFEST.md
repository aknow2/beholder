## Installation Manifest Template

The installation manifest is a simple text file created at `~/.beholder/install-manifest.txt` after successful installation.

### Format

```
INSTALLED_VERSION=v1.0.0
INSTALLATION_PATH=/home/user/.local/bin/beholder
INSTALLATION_DATE=2026-02-06T10:30:45Z
USER_DATA_PATH=/home/user/.beholder
```

### Fields

- **INSTALLED_VERSION**: Version of the installed binary (e.g., `v1.0.0`)
- **INSTALLATION_PATH**: Full path to the installed binary executable
- **INSTALLATION_DATE**: ISO 8601 timestamp of installation
- **USER_DATA_PATH**: Location where user data (config, database, images) is stored

### Usage

- **Detection**: Installation scripts read this manifest to detect existing installations
- **Upgrade**: Version comparison is done by reading `INSTALLED_VERSION`
- **Uninstall**: Provides the exact binary path and user data location for cleanup

### Example Script Creation (Bash)

```bash
cat > ~/.beholder/install-manifest.txt << EOF
INSTALLED_VERSION=${VERSION}
INSTALLATION_PATH=${INSTALL_PATH}
INSTALLATION_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
USER_DATA_PATH=~/.beholder
EOF
```

### Example Script Reading (Bash)

```bash
if [ -f ~/.beholder/install-manifest.txt ]; then
  source ~/.beholder/install-manifest.txt
  echo "Detected installation at: $INSTALLATION_PATH"
  echo "Version: $INSTALLED_VERSION"
fi
```
