# Binary Naming Convention

## Release Binary Naming

Built binaries for GitHub Releases follow this pattern:

```
beholder-v{VERSION}-{OS}-{ARCH}{EXT}
```

### Components

- **VERSION**: Semantic version (e.g., `1.0.0`, `0.1.0-alpha`)
- **OS**: Operating System
  - `linux` (Linux)
  - `darwin` (macOS)
  - `windows` (Windows)
- **ARCH**: CPU Architecture
  - `amd64` (x86-64 / Intel 64-bit)
  - `arm64` (ARM 64-bit / Apple Silicon)
- **EXT**: File extension
  - `.exe` (Windows only)
  - (empty) (Unix-like systems)

### Examples

- `beholder-v1.0.0-linux-amd64` - Linux x86-64
- `beholder-v1.0.0-darwin-amd64` - macOS Intel
- `beholder-v1.0.0-darwin-arm64` - macOS Apple Silicon
- `beholder-v1.0.0-windows-amd64.exe` - Windows x86-64

## GitHub Releases Assets

All binaries are published as release assets on GitHub Releases with the naming convention above.

Checksums (SHA256) are provided in a separate `beholder-v{VERSION}-checksums.txt` file.

## Installation Scripts

Installation scripts use the naming convention to construct download URLs:

```bash
DOWNLOAD_URL="https://github.com/aknow2/beholder/releases/download/${TAG}/${BINARY_NAME}"
```

Example:
```
https://github.com/aknow2/beholder/releases/download/v1.0.0/beholder-v1.0.0-linux-amd64
```
