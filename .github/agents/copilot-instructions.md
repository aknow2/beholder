# beholder Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-01-24

## Active Technologies
- Go 1.24 (existing codebase), Shell script (Bash/sh), PowerShell 5.1+ + GitHub Actions (workflow automation), Go toolchain (cross-compilation), GitHub Artifacts API (binary distribution) (003-beholder-installer)
- GitHub Artifacts for versioned binaries (90-day retention minimum); local filesystem for installed binaries (`~/.beholder/bin` or `/usr/local/bin` or `C:\Program Files\beholder`) (003-beholder-installer)

- Go 1.24 + Copilot SDK (go), kbinani/screenshot, systray, modernc.org/sqlite (001-screenshot-daily-report)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.24

## Code Style

Go 1.24: Follow standard conventions

## Recent Changes
- 003-beholder-installer: Added Go 1.24 (existing codebase), Shell script (Bash/sh), PowerShell 5.1+ + GitHub Actions (workflow automation), Go toolchain (cross-compilation), GitHub Artifacts API (binary distribution)

- 001-screenshot-daily-report: Added Go 1.24 + Copilot SDK (go), kbinani/screenshot, systray, modernc.org/sqlite

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
