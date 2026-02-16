# Data Model: Beholder Installer

**Feature**: 003-beholder-installer  
**Date**: 2026-02-02  
**Phase**: 1 - Design

## Overview

This document defines the data entities involved in the installation, upgrade, and uninstallation workflows. These entities are technology-agnostic and describe **what** information is tracked, not **how** it is implemented.

---

## Entities

### 1. Build Artifact

**Description**: Represents a compiled binary produced by the CI/CD pipeline for a specific platform and architecture.

**Attributes**:
- `version`: Semantic version string (e.g., "v1.2.3")
- `platform`: Operating system (darwin, windows)
- `architecture`: CPU architecture (amd64, arm64)
- `checksum`: SHA256 hash of the binary file (for integrity verification)
- `download_url`: HTTPS URL where the binary can be retrieved
- `build_timestamp`: ISO 8601 timestamp of when the build completed
- `file_size`: Size of the binary in bytes

**Lifecycle**:
- **Created**: By GitHub Actions workflow on each release/tag
- **Published**: To GitHub Releases as a downloadable asset
- **Consumed**: By installation scripts during download phase

**Example**:
```json
{
  "version": "v0.1.0",
  "platform": "darwin",
  "architecture": "amd64",
  "checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "download_url": "https://github.com/aknow2/beholder/releases/download/v0.1.0/beholder-darwin-amd64",
  "build_timestamp": "2026-02-02T10:30:00Z",
  "file_size": 8388608
}
```

---

### 2. Installation Receipt

**Description**: Local record of an installed Beholder instance, used by uninstaller to track what was installed and where.

**Attributes**:
- `installed_version`: Version string of the installed binary (e.g., "v1.2.3")
- `installation_path`: Absolute filesystem path where binary was placed
- `installation_date`: ISO 8601 timestamp of when installation completed
- `user_data_locations`: List of directories containing user-created data (config, database, screenshots)
- `path_modified`: Boolean indicating if shell PATH was modified
- `profile_files`: List of shell profile files that were modified (for PATH cleanup during uninstall)

**Lifecycle**:
- **Created**: During installation script execution
- **Persisted**: As a text file (e.g., `~/.beholder/install-manifest.txt`)
- **Read**: By uninstaller to determine what to remove
- **Deleted**: During uninstallation (or preserved with user data if user declines data removal)

**Example**:
```text
installed_version=v0.1.0
installation_path=/home/user/.local/bin/beholder
installation_date=2026-02-02T15:45:00Z
user_data_locations=/home/user/.beholder
path_modified=true
profile_files=/home/user/.bashrc
```

---

### 3. Version Manifest

**Description**: Remote index of all available Beholder releases, used by installation scripts to discover the latest version and available platforms.

**Attributes**:
- `latest_version`: Version string of the most recent stable release
- `releases`: Array of release metadata (version, date, changelog URL)
- `artifacts_by_version`: Mapping of version → array of Build Artifacts

**Lifecycle**:
- **Generated**: Implicitly by GitHub Releases API
- **Queried**: By installation scripts to determine latest version
- **Updated**: Automatically when new releases are published

**Example** (conceptual; actual data comes from GitHub API):
```json
{
  "latest_version": "v0.1.0",
  "releases": [
    {
      "version": "v0.1.0",
      "published_date": "2026-02-02T10:30:00Z",
      "changelog_url": "https://github.com/aknow2/beholder/releases/tag/v0.1.0"
    }
  ],
  "artifacts_by_version": {
    "v0.1.0": [
      {
        "platform": "darwin",
        "architecture": "amd64",
        "download_url": "https://github.com/aknow2/beholder/releases/download/v0.1.0/beholder-darwin-amd64"
      },
      {
        "platform": "darwin",
        "architecture": "arm64",
        "download_url": "https://github.com/aknow2/beholder/releases/download/v0.1.0/beholder-darwin-arm64"
      },
      {
        "platform": "windows",
        "architecture": "amd64",
        "download_url": "https://github.com/aknow2/beholder/releases/download/v0.1.0/beholder-windows-amd64.exe"
      }
    ]
  }
}
```

---

## Entity Relationships

```
Version Manifest (1) ──contains──> (N) Build Artifact
       │
       │ queried by
       ↓
Installation Script ──downloads──> Build Artifact
       │
       │ creates
       ↓
Installation Receipt ──references──> Build Artifact.version
       │
       │ read by
       ↓
Uninstaller Script
```

---

## State Transitions

### Installation States

1. **Not Installed**: No Installation Receipt exists, `beholder` command not in PATH
2. **Installing**: Installation script running, binary downloaded but not yet in final location
3. **Installed**: Installation Receipt exists, `beholder` command available in PATH
4. **Upgrade Available**: Installed version < latest version in Version Manifest
5. **Uninstalling**: Uninstaller script running, binary being removed

### Installation Flow

```
Not Installed
    ↓ (user runs install script)
Installing
    ↓ (download success, binary validated, moved to install path, PATH updated)
Installed
    ↓ (user runs install script again, newer version available)
Installing (upgrade)
    ↓ (existing binary replaced, Installation Receipt updated)
Installed (newer version)
```

### Uninstallation Flow

```
Installed
    ↓ (user runs uninstall script)
Uninstalling
    ↓ (binary removed, PATH cleaned)
    ├─ (user declines data removal)
    │   → Not Installed (user data preserved)
    └─ (user accepts data removal)
        → Not Installed (user data removed)
```

---

## Validation Rules

### Build Artifact
- `version` MUST match semantic versioning format (vX.Y.Z)
- `platform` MUST be one of: darwin, windows
- `architecture` MUST be one of: amd64, arm64
- `checksum` MUST be a valid SHA256 hash (64 hex characters)
- `download_url` MUST be an HTTPS URL
- `file_size` MUST be > 0 and < 100MB (sanity check)

### Installation Receipt
- `installed_version` MUST match semantic versioning format
- `installation_path` MUST be an absolute filesystem path
- `installation_date` MUST be a valid ISO 8601 timestamp
- `user_data_locations` entries MUST be absolute filesystem paths

### Version Manifest
- `latest_version` MUST match semantic versioning format
- `releases` array MUST be sorted by `published_date` descending
- Each artifact in `artifacts_by_version` MUST have a unique (platform, architecture) pair per version

---

## Privacy & Security Considerations

- Installation Receipt contains only local filesystem paths (no PII)
- No telemetry or analytics data is collected during installation/uninstallation
- Build Artifacts checksums enable integrity verification (prevents tampered binaries)
- Download URLs use HTTPS to prevent man-in-the-middle attacks

---

## Next Steps

1. Use this data model to design API contracts in `contracts/github-artifacts-api.yaml`
2. Implement Installation Receipt persistence in installation scripts
3. Query Version Manifest via GitHub Releases API in installation scripts
4. Validate Build Artifacts against this schema in CI/CD pipeline
