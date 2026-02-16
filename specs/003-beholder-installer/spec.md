# Feature Specification: Beholder Installer & Uninstaller

**Feature Branch**: `003-beholder-installer`  
**Created**: 2026-02-02  
**Status**: Draft  
**Input**: User description: "配布するためにインストーラーを作成したい. 最終的にはterminalやcmdでbeholderコマンドが使えるようにインストールでき、アンインストールも出来るようにしたい"

## Clarifications

### Session 2026-02-02

- Q: One-line install URLs should use which hosting target? → A: raw.githubusercontent.com on the main branch.

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Install and Run the `beholder` Command (Priority: P1)

A user can obtain a distributable installer, run it, and then use the `beholder` command from their usual command environment (terminal/shell on macOS and Command Prompt/PowerShell on Windows).

**Why this priority**: This is the core distribution value: a new user must be able to install and successfully run the product.

**Independent Test**: Can be fully tested on a clean machine/user account by running the installer and confirming the `beholder` command is available and launches successfully.

**Acceptance Scenarios**:

1. **Given** a machine where `beholder` is not installed, **When** the user runs the installer with default options, **Then** the installation completes successfully and the user can run `beholder` from their command environment.
2. **Given** a successful installation, **When** the user runs `beholder` without additional setup, **Then** the command displays a helpful response (e.g., usage/help) and does not crash.

---

### User Story 2 - Uninstall Cleanly (Priority: P2)

A user can uninstall the application so the `beholder` command is no longer present, and installed application files are removed.

**Why this priority**: Uninstallability is a distribution requirement for trust and operational hygiene (especially in managed environments).

**Independent Test**: Can be fully tested by installing, uninstalling, and verifying the command is removed and the application is not listed/registered as installed.

**Acceptance Scenarios**:

1. **Given** `beholder` is installed, **When** the user runs the uninstaller (or uses the OS uninstall flow), **Then** the uninstall completes successfully and `beholder` can no longer be invoked.
2. **Given** `beholder` is installed and has created user data, **When** the user uninstalls, **Then** the user is clearly informed what happens to user data and can choose whether to remove it.

---

### User Story 3 - Automated Build and Release Pipeline (Priority: P2)

The development team can trigger a build, produce platform-specific binaries, and store them in a centralized location for distribution without manual compilation or file management on each developer's machine.

**Why this priority**: Consistent, reproducible builds and reliable distribution are foundational to end-user installation; core to the entire feature's viability.

**Independent Test**: Can be tested by verifying that a code commit triggers automated builds, binaries are generated for all supported platforms, and are accessible from artifact storage.

**Acceptance Scenarios**:

1. **Given** a commit to the main/release branch, **When** GitHub Actions workflow is triggered, **Then** binaries are compiled for Windows and macOS and are available in GitHub Releases assets.
2. **Given** binaries are stored in GitHub Releases assets, **When** the installation script downloads a binary, **Then** the download succeeds and the binary is executable and correct.

---

### User Story 4 - One-Line Installation via curl/PowerShell (Priority: P3)

A user can install `beholder` using a single command line (via curl on Unix-like systems or Invoke-WebRequest on Windows), without needing to manually download a file or navigate a GUI installer.

**Why this priority**: Developer convenience and discoverability; aligns with modern distribution patterns (pnpm, rustup, etc.).

**Independent Test**: Can be tested by running a curl or PowerShell one-liner on a fresh machine and verifying installation completes and `beholder` command is available.

**Acceptance Scenarios**:

1. **Given** a clean Unix-like system with curl installed, **When** the user runs `curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.sh | sh`, **Then** the installation completes and `beholder` is available.
2. **Given** a clean Windows system with PowerShell, **When** the user runs `Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.ps1 -UseBasicParsing | Invoke-Expression`, **Then** the installation completes and `beholder` is available.
3. **Given** the user has completed a one-line install, **When** the user runs a newer installer script, **Then** the existing installation is detected and upgraded safely.

---

### User Story 5 - Upgrade Existing Installations (Priority: P3)

A user with an existing installation can install a newer version without manual cleanup and without losing their existing settings/data by default.

**Why this priority**: Smooth upgrades reduce support burden and enable regular releases.

**Independent Test**: Can be tested by installing version A, creating some user configuration/data, then installing version B and confirming the command works and existing user data remains.

**Acceptance Scenarios**:

1. **Given** an older version is installed, **When** the user runs the newer installer, **Then** the system completes the upgrade and the `beholder` command now corresponds to the newer version.
2. **Given** user data exists from the previous installation, **When** the upgrade completes, **Then** existing user settings/data remain available (unless the user explicitly chooses otherwise).

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right edge cases.
-->

- Installation is started without sufficient permissions (system-wide install) and must clearly guide the user to a successful path.
- The `beholder` command name conflicts with an existing command or an older manual installation.
- Installation is interrupted (power loss / forced shutdown) and must avoid leaving a broken half-installed state.
- Uninstall is attempted while the application is running.
- Upgrade is attempted across multiple previously installed versions.

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: The product MUST provide distributable installer artifacts for supported operating systems. (Acceptance: A user can download an installer appropriate for their operating system.)
- **FR-002**: Installation MUST result in the `beholder` command being usable from the user’s normal command environment without requiring manual file copying. (Acceptance: After install, `beholder` can be invoked in that environment.)
- **FR-003**: The installer MUST confirm successful installation and provide a clear next step to verify it. (Acceptance: The installer communicates success and tells the user how to verify the command works.)
- **FR-004**: The product MUST provide a supported uninstall flow for each supported operating system. (Acceptance: A user can start an uninstall using common OS conventions.)
- **FR-005**: Uninstall MUST remove the installed application and ensure the `beholder` command is no longer available. (Acceptance: After uninstall, invoking `beholder` fails in the same environment it previously worked.)
- **FR-006**: Uninstall MUST clearly communicate how user-created data is handled and MUST allow the user to opt into removing user data. (Acceptance: The uninstall flow presents a clear choice and executes it correctly.)
- **FR-007**: The installer MUST detect an existing installation and provide a safe default path to upgrade/repair. (Acceptance: Re-running the installer offers an upgrade/repair path and completes without manual cleanup.)
- **FR-008**: Upgrades MUST preserve user settings/data by default. (Acceptance: After upgrade with default options, prior user settings/data remain available.)
- **FR-009**: The installer and uninstaller MUST support unattended (non-interactive) execution. (Acceptance: Installation and uninstall can be executed without user interaction and still provide success/failure results.)
- **FR-010**: Installation and uninstall operations MUST produce user-visible outcomes that are actionable (clear success/failure and next steps). (Acceptance: Errors are understandable and include next steps; success is explicit.)
- **FR-011**: The product MUST provide shell installation scripts (Unix-like systems) and PowerShell scripts (Windows) that can be downloaded and piped directly into interpreters. (Acceptance: `curl | sh` and `Invoke-WebRequest | Invoke-Expression` patterns work as documented.)
- **FR-012**: The installation scripts MUST be remotely hosted at well-known, stable URLs (e.g., `https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.sh` and `https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.ps1`). (Acceptance: The documented URLs are accessible and return the correct script.)
- **FR-013**: Installation scripts MUST detect and handle existing installations by offering upgrade/repair options. (Acceptance: Running the script on an existing installation detects it and provides safe upgrade.)
- **FR-014**: Installation scripts MUST fail safely and provide clear error messages if prerequisites are unmet or the environment is incompatible. (Acceptance: Missing dependencies or incompatible OS are reported with remediation steps.)

### Automated Build and Artifact Management

- **FR-015**: The project MUST include a GitHub Actions workflow that automatically builds binaries for supported platforms (Windows, macOS) on each commit/release. (Acceptance: Each commit to the release branch triggers a build; no manual intervention required.)
- **FR-016**: Built binaries MUST be stored as GitHub Releases assets with consistent naming that allows installation scripts to locate and download the correct version and platform. (Acceptance: Release assets follow a naming convention (e.g., `beholder-v{version}-{platform}{.exe}`) and are retrievable via URL.)
- **FR-017**: GitHub Actions workflow MUST validate that built binaries are functional before storing them as artifacts (basic smoke test). (Acceptance: A binary is tested for existence and executability; if tests fail, the artifact is not published.)
- **FR-018**: Release assets MUST be retained for at least 90 days to support installation and upgrades. (Acceptance: GitHub Releases assets are retained; if GitHub Actions artifacts are used, retention is configured to at least 90 days.)
- **FR-019**: GitHub Actions workflow MUST tag and version builds consistently (e.g., using git tags or semantic versioning). (Acceptance: Each artifact is associated with a clear version identifier that installation scripts can reference.)

### Key Entities *(include if feature involves data)*

- **Distribution Artifact**: A downloadable item a user runs to install the product (includes version, supported operating system).
- **Build Artifact**: A compiled binary or installer file produced by GitHub Actions CI/CD (stored in GitHub Releases assets, versioned and platform-specific).
- **Installation Record**: What the system uses to recognize the product as installed (installed version, install scope such as per-user/system).
- **User Data**: User-created settings and local data generated while using the product (retention/removal behavior must be explicit).

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: On a clean machine/user account, users can install and successfully run `beholder` within 5 minutes (standard installer or one-liner script).
- **SC-002**: On a clean machine/user account, users can install via one-liner script (curl or Invoke-WebRequest) within 2 minutes.
- **SC-003**: On an installed machine, users can fully uninstall within 3 minutes and `beholder` is no longer invokable afterward.
- **SC-004**: At least 95% of upgrade attempts (from the previous released version) complete successfully without user data loss when default options are used.
- **SC-005**: Installation-related support requests drop by at least 30% within one release cycle after shipping the installer/uninstaller.
- **SC-006**: Build artifacts are produced for all supported platforms (Windows, macOS) within 15 minutes of code commit/merge.
- **SC-007**: At least 99% of builds complete successfully and produce functional, executable binaries for all platforms.

## Assumptions

- Supported operating systems include Windows and macOS for the initial release.
- Default behavior preserves user-created data across upgrades and uninstalls (unless the user opts into removal).
- The primary distribution goal is to enable end users to install and remove the product without developer tools.
- One-liner scripts (curl/PowerShell) require network access to the remote script hosting service.
- Scripts will be hosted with HTTPS and appropriate security measures (signing/verification recommended but not a hard requirement for MVP).
- GitHub Actions CI/CD is available for the repository (standard for public GitHub repositories).
- Binaries are built for Windows (x86-64) and macOS (Intel and Apple Silicon) as supported platforms.
- GitHub Releases assets are used as the authoritative source for versioned binaries; installation scripts download from GitHub Releases.
- Build reproducibility is ensured through pinned dependencies and consistent build environment (Go version, etc.).

## Out of Scope

- Providing integrations with third-party package managers (can be added later).
- Changing the product’s core features unrelated to installation/uninstallation.

## Dependencies

- A defined release process that produces versioned artifacts for distribution.
- Any required organizational approvals for distributing installable software (e.g., security review/signing policies).
