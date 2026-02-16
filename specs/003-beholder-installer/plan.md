# Implementation Plan: Beholder Installer & Uninstaller

**Branch**: `003-beholder-installer` | **Date**: 2026-02-02 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `/specs/003-beholder-installer/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command.

## Summary

Provide distributable installation and uninstallation mechanisms for the Beholder CLI tool across Windows and macOS. Primary requirements include: (1) automated build pipeline via GitHub Actions producing platform-specific binaries, (2) one-line installation scripts (curl/PowerShell) that download and install the binary from GitHub Releases assets, (3) PATH management to enable `beholder` command invocation, and (4) clean uninstallation with user data handling. Technical approach leverages GitHub Actions for cross-compilation, shell/PowerShell scripts for installation automation, and platform-specific conventions for binary placement and PATH configuration.

## Technical Context

**Language/Version**: Go 1.24 (existing codebase), Shell script (Bash/sh), PowerShell 5.1+  
**Primary Dependencies**: GitHub Actions (workflow automation), Go toolchain (cross-compilation), GitHub Releases API (binary distribution)  
**Storage**: GitHub Releases assets for versioned binaries (90-day retention minimum); local filesystem for installed binaries (`~/.local/bin` or `%USERPROFILE%\.beholder\bin`)  
**Testing**: Go test for smoke tests of built binaries; shell script validation on target platforms (macOS/Windows)  
**Target Platform**: Cross-platform - macOS (Intel x86-64 + Apple Silicon arm64), Windows (x86-64)  
**Project Type**: CLI distribution (installation scripts + CI/CD automation)  
**Performance Goals**: Build completion within 15 minutes; installation script execution within 2 minutes on typical network  
**Constraints**: No external runtime dependencies (statically linked Go binaries); installation must work without admin/root privileges (user-local install); scripts must be idempotent (re-run safely)  
**Scale/Scope**: Support concurrent installations across multiple platforms; handle 100+ releases over product lifetime; maintain backward compatibility for upgrade paths

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Modular Architecture

✅ **PASS** - Distribution infrastructure is cleanly separated from core application:
- GitHub Actions workflow: `.github/workflows/release.yml` (isolated from internal/)
- Installation scripts: `scripts/install.sh`, `scripts/install.ps1` (separate from cmd/beholder)
- No changes to existing internal/ package boundaries
- Distribution logic does not create circular dependencies

### II. CLI-First Interface

✅ **PASS** - Installer does not modify existing CLI contract:
- The `beholder` command interface remains unchanged
- Installation scripts place the binary in PATH but do not alter command behavior
- Uninstaller does not introduce new CLI commands (uses system-level uninstall mechanisms)

### III. Configuration-Driven

✅ **PASS** - Installation respects existing configuration model:
- Installed binary uses existing `~/.beholder/config.yaml` discovery
- Installation does not modify or override user configuration
- Default configuration initialization remains handled by the application on first run

### IV. Local-First & Privacy

✅ **PASS** - Installation preserves local-first principle:
- Binaries are downloaded once and stored locally
- No telemetry or tracking in installation scripts
- User data handling during uninstall is explicit and user-controlled (FR-006)
- Installation scripts detect existing installations without remote calls

### V. Incremental Enhancement

✅ **PASS** - Feature follows staged delivery:
- P1: Basic installer (manual binary placement)
- P2: Automated builds (GitHub Actions)
- P3: One-line scripts (curl/PowerShell)
- Each priority level is independently testable and deployable

**Constitution Compliance**: APPROVED - No violations. Distribution infrastructure is additive and does not compromise existing architectural principles.

## Project Structure

### Documentation (this feature)

```text
specs/003-beholder-installer/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output: best practices for installers, GitHub Actions patterns
├── data-model.md        # Phase 1 output: installation state, artifact metadata
├── quickstart.md        # Phase 1 output: developer setup for testing installers
├── contracts/           # Phase 1 output: API contracts for artifact download
│   └── github-artifacts-api.yaml
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Distribution infrastructure (new)
.github/
└── workflows/
    └── release.yml           # GitHub Actions workflow for cross-platform builds

scripts/
├── install.sh                # POSIX installation script (curl | sh pattern)
├── install.ps1               # PowerShell installation script (Invoke-WebRequest pattern)
├── uninstall.sh              # POSIX uninstallation script
└── uninstall.ps1             # PowerShell uninstallation script

# Existing application (unchanged structure)
cmd/
└── beholder/
    ├── main.go
    └── cli.go

internal/
├── app/
├── classify/
├── config/
├── scheduler/
├── storage/
└── summary/

# Build outputs (not committed)
dist/                         # Local build artifacts for testing
└── beholder-{version}-{platform}{.exe}
```

**Structure Decision**: Distribution infrastructure is separated into `.github/workflows/` for automation and `scripts/` for installation logic. This maintains a clean separation between application code (existing `cmd/` and `internal/`) and distribution mechanics, adhering to the Modular Architecture principle.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

*No violations - this section is empty.*

## Phase 0: Research & Unknowns

**Objective**: Resolve all NEEDS CLARIFICATION markers from Technical Context and establish best practices for installer implementation.

### Research Tasks

1. **GitHub Actions cross-compilation patterns for Go**
   - Research: How to build Go binaries for multiple platforms (GOOS/GOARCH) in a single workflow
   - Research: Best practices for matrix builds (parallel vs sequential)
   - Research: Artifact naming conventions for versioned releases
   - Output: Document recommended workflow structure in `research.md`

2. **GitHub Artifacts API usage from shell scripts**
   - Research: How to download artifacts from public repositories without authentication
   - Research: URL patterns for accessing latest release artifacts vs tagged versions
   - Research: Fallback strategies when artifact download fails
   - Output: Document artifact download patterns and error handling in `research.md`

3. **Platform-specific binary installation conventions**
   - Research: Standard installation paths for user-local vs system-wide installs on macOS and Windows
   - Research: PATH management strategies (profile modification vs symlinks vs user bin directories)
   - Research: Permission requirements and privilege escalation patterns
   - Output: Document installation paths and PATH strategies per platform in `research.md`

4. **Installer script best practices**
   - Research: Idempotency patterns (detecting existing installations, safe overwrites)
   - Research: Error handling and rollback strategies for failed installations
   - Research: User feedback patterns (progress indicators, success/failure messaging)
   - Research: Security considerations (verifying checksums, avoiding arbitrary code execution)
   - Output: Document installer script template patterns in `research.md`

5. **Uninstaller implementation patterns**
   - Research: How to track installed files for clean removal (manifest files, installation receipts)
   - Research: User data handling conventions (prompt vs preserve vs remove)
   - Research: How to remove PATH entries safely across different shells
   - Output: Document uninstaller patterns and user data handling in `research.md`

**Deliverable**: `research.md` with all decisions documented and alternatives evaluated.

## Phase 1: Design & Contracts

**Prerequisites**: Phase 0 research complete

### Data Model

**File**: `data-model.md`

Entities:
- **Build Artifact**: Metadata for compiled binaries (version, platform, architecture, checksum, download URL, build timestamp)
- **Installation Receipt**: Local record of installation (installed version, installation path, installation date, user data locations)
- **Version Manifest**: Remote index of available versions (latest stable, all releases, checksums)

### API Contracts

**Directory**: `contracts/`

Files:
- `github-artifacts-api.yaml`: OpenAPI/informal spec for GitHub Artifacts API endpoints used by installation scripts
  - GET artifact download URLs by version/platform
  - Response schemas for artifact metadata
  - Error responses and retry patterns

### Developer Quickstart

**File**: `quickstart.md`

Content:
- How to test the GitHub Actions workflow locally (using `act` or manual cross-compilation)
- How to test installation scripts on each platform (VM setup, testing procedure)
- How to manually trigger builds and download artifacts
- How to verify installed binary works correctly

### Agent Context Update

Run `.specify/scripts/bash/update-agent-context.sh copilot` to update agent-specific context with new technologies:
- GitHub Actions workflow definitions
- Shell scripting patterns
- PowerShell scripting patterns
- Cross-platform installation conventions

**Deliverables**: `data-model.md`, `contracts/github-artifacts-api.yaml`, `quickstart.md`, updated agent context

## Phase 2: Implementation Planning (Task Breakdown)

**Note**: This phase is handled by the `/speckit.tasks` command and produces `tasks.md`. It is NOT part of this plan document.

The tasks will be organized by user story priority:
1. P1: Basic binary installation (manual download)
2. P2: Automated GitHub Actions builds
3. P2: Uninstaller scripts
4. P3: One-line installation scripts (curl/PowerShell)
5. P3: Upgrade detection and handling

---

**Next Steps**:
1. Review and approve this plan
2. Execute Phase 0 research (generates `research.md`)
3. Execute Phase 1 design (generates `data-model.md`, `contracts/`, `quickstart.md`)
4. Re-validate Constitution Check after design
5. Run `/speckit.tasks` to generate task breakdown in `tasks.md`
