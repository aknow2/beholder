---
description: "Task list for Beholder Installer & Uninstaller implementation"
---

# Tasks: Beholder Installer & Uninstaller

**Feature**: 003-beholder-installer  
**Input**: Design documents from `/specs/003-beholder-installer/`  
**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, contracts/ ✓, quickstart.md ✓

**Tests**: Not explicitly requested in feature specification - focus on manual testing per quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `- [ ] [ID] [P?] [Story?] Description with file path`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and directory structure

- [x] T001 Create `scripts/` directory at repository root for installation/uninstallation scripts
- [x] T002 Create `.github/workflows/` directory for GitHub Actions workflow
- [x] T003 [P] Add `dist/` to `.gitignore` for local build artifacts

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Add version metadata to `cmd/beholder/main.go` (ldflags injection support for -X main.version)
- [x] T005 Add `--version` flag to CLI in `cmd/beholder/cli.go` to display version string
- [x] T006 Create installation manifest template structure (defines Installation Receipt schema)
- [x] T007 Document binary naming convention in `.github/RELEASE_NAMING.md` (e.g., `beholder-{os}-{arch}`)

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 3 - Automated Build and Release Pipeline (Priority: P2) 🎯

**Goal**: GitHub Actions workflow automatically builds and publishes platform-specific binaries on release

**Independent Test**: Push a git tag, verify workflow runs, and all platform binaries are available as GitHub Release assets

### Implementation for User Story 3

- [x] T008 [US3] Create GitHub Actions workflow file in `.github/workflows/release.yml`
- [x] T009 [US3] Define workflow triggers (on push of tags matching `v*` pattern) in `.github/workflows/release.yml`
- [x] T010 [US3] Configure matrix build strategy for platforms (darwin-amd64, darwin-arm64, windows-amd64) in `.github/workflows/release.yml`
- [x] T011 [US3] Add Go setup step (actions/setup-go@v5, version 1.24) in `.github/workflows/release.yml`
- [x] T012 [US3] Add build step with GOOS/GOARCH environment variables and ldflags for version injection in `.github/workflows/release.yml`
- [x] T013 [US3] Add binary validation step (smoke test: run `beholder --version`) in `.github/workflows/release.yml`
- [x] T014 [US3] Add artifact upload step (actions/upload-artifact@v4) with platform-specific naming in `.github/workflows/release.yml`
- [x] T015 [US3] Add GitHub Release creation step (creates release from tag with all binaries as assets) in `.github/workflows/release.yml`
- [x] T016 [US3] Test workflow locally using manual cross-compilation per `quickstart.md` instructions
- [x] T017 [US3] Create test release (tag v0.0.1-test) and validate all artifacts are published correctly

**Checkpoint**: At this point, automated builds should be fully functional and binaries available in GitHub Releases

---

## Phase 4: User Story 1 - Install and Run the `beholder` Command (Priority: P1)

**Goal**: Users can manually download a binary and place it in their PATH to use the `beholder` command

**Independent Test**: Download platform-appropriate binary from GitHub Release, move to installation path, run `beholder --version`

### Implementation for User Story 1

- [x] T018 [P] [US1] Create POSIX installation script skeleton in `scripts/install.sh`
- [x] T019 [P] [US1] Create PowerShell installation script skeleton in `scripts/install.ps1`
- [x] T020 [US1] Implement platform/architecture detection in `scripts/install.sh` (uname -s, uname -m, map to Go naming)
- [x] T021 [US1] Implement platform/architecture detection in `scripts/install.ps1` (hardcoded to windows-amd64)
- [x] T022 [US1] Implement GitHub Releases API query for latest version in `scripts/install.sh` (curl to api.github.com/repos/aknow2/beholder/releases/latest)
- [x] T023 [US1] Implement GitHub Releases API query for latest version in `scripts/install.ps1` (Invoke-RestMethod)
- [x] T024 [US1] Implement binary download with retry logic in `scripts/install.sh` (curl with exponential backoff, 3 retries)
- [x] T025 [US1] Implement binary download with retry logic in `scripts/install.ps1` (Invoke-WebRequest with retry)
- [x] T026 [US1] Implement binary validation (file exists, size check, execute --version) in `scripts/install.sh`
- [x] T027 [US1] Implement binary validation in `scripts/install.ps1`
- [x] T028 [US1] Implement installation to user-local path (`~/.local/bin` on Unix, `%USERPROFILE%\.beholder\bin` on Windows) in `scripts/install.sh`
- [x] T029 [US1] Implement installation to user-local path in `scripts/install.ps1`
- [x] T030 [US1] Implement PATH modification (append to shell profile) in `scripts/install.sh` (detect .bashrc/.zshrc/.profile)
- [x] T031 [US1] Implement PATH modification (update user environment variable) in `scripts/install.ps1`
- [x] T032 [US1] Implement installation manifest creation in `scripts/install.sh` (write to ~/.beholder/install-manifest.txt)
- [x] T033 [US1] Implement installation manifest creation in `scripts/install.ps1`
- [x] T034 [US1] Add user feedback messages (download progress, success, verification instructions) in `scripts/install.sh`
- [x] T035 [US1] Add user feedback messages in `scripts/install.ps1`
- [x] T036 [US1] Add error handling with actionable messages (download failed, invalid binary, permission denied) in `scripts/install.sh`
- [x] T037 [US1] Add error handling in `scripts/install.ps1`
- [ ] T038 [US1] Test installation on macOS (Intel and Apple Silicon if available) per `quickstart.md` procedure
- [ ] T039 [US1] Test installation on Windows 11 per `quickstart.md` procedure
- [ ] T040 [US1] Test installation on Windows 11 (Command Prompt and PowerShell) per `quickstart.md` procedure

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 5: User Story 2 - Uninstall Cleanly (Priority: P2)

**Goal**: Users can uninstall the beholder binary and optionally remove user data

**Independent Test**: Install, uninstall with data preservation, verify binary removed and data remains; then install, uninstall with data removal, verify all removed

### Implementation for User Story 2

- [x] T041 [P] [US2] Create POSIX uninstallation script skeleton in `scripts/uninstall.sh`
- [x] T042 [P] [US2] Create PowerShell uninstallation script skeleton in `scripts/uninstall.ps1`
- [x] T043 [US2] Implement installation manifest reading in `scripts/uninstall.sh` (source ~/.beholder/install-manifest.txt)
- [x] T044 [US2] Implement installation manifest reading in `scripts/uninstall.ps1`
- [x] T045 [US2] Implement binary removal (delete from installation path) in `scripts/uninstall.sh`
- [x] T046 [US2] Implement binary removal in `scripts/uninstall.ps1`
- [x] T047 [US2] Implement PATH cleanup (remove entries from shell profiles) in `scripts/uninstall.sh`
- [x] T048 [US2] Implement PATH cleanup (remove from user environment variable) in `scripts/uninstall.ps1`
- [x] T049 [US2] Implement user data handling prompt (preserve vs remove ~/.beholder) in `scripts/uninstall.sh`
- [x] T050 [US2] Implement user data handling prompt in `scripts/uninstall.ps1`
- [x] T051 [US2] Implement user data removal logic (rm -rf ~/.beholder if user confirms) in `scripts/uninstall.sh`
- [x] T052 [US2] Implement user data removal logic in `scripts/uninstall.ps1`
- [x] T053 [US2] Add user feedback messages (binary removed, PATH cleaned, data preserved/removed) in `scripts/uninstall.sh`
- [x] T054 [US2] Add user feedback messages in `scripts/uninstall.ps1`
- [x] T055 [US2] Add graceful handling when installation manifest is missing (fallback to `which beholder`) in `scripts/uninstall.sh`
- [x] T056 [US2] Add graceful handling when installation manifest is missing in `scripts/uninstall.ps1`
- [ ] T057 [US2] Test uninstallation on macOS (preserve data scenario) per `quickstart.md` procedure
- [ ] T058 [US2] Test uninstallation on macOS (remove data scenario) per `quickstart.md` procedure
- [ ] T059 [US2] Test uninstallation on Windows (preserve data scenario) per `quickstart.md` procedure
- [ ] T060 [US2] Test uninstallation on Windows (remove data scenario) per `quickstart.md` procedure

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently

---

## Phase 6: User Story 4 - One-Line Installation via curl/PowerShell (Priority: P3)

**Goal**: Users can install with a single curl | sh or Invoke-WebRequest | Invoke-Expression command

**Independent Test**: Run one-liner command on fresh machine, verify installation completes and beholder command works

### Implementation for User Story 4

- [x] T061 [US4] Configure repository to serve raw installation scripts via raw.githubusercontent.com
- [x] T062 [US4] Update `scripts/install.sh` to support unattended execution (no interactive prompts by default)
- [x] T063 [US4] Update `scripts/install.ps1` to support unattended execution
- [x] T064 [US4] Add `--interactive` flag to `scripts/install.sh` for explicit prompt mode
- [x] T065 [US4] Add `-Interactive` parameter to `scripts/install.ps1`
- [ ] T066 [US4] Test one-liner installation on Ubuntu 22.04: `curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.sh | sh`
- [ ] T067 [US4] Test one-liner installation on macOS
- [ ] T068 [US4] Test one-liner installation on Windows: `Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.ps1 -UseBasicParsing | Invoke-Expression`
- [ ] T069 [US4] Update README.md with one-liner installation instructions
- [ ] T070 [US4] Create installation documentation page with platform-specific instructions

**Checkpoint**: At this point, one-line installation should work on all platforms

---

## Phase 7: User Story 5 - Upgrade Existing Installations (Priority: P3)

**Goal**: Users can upgrade an existing installation without losing user data

**Independent Test**: Install v1, create user data, run installer again for v2, verify upgrade succeeds and data persists

### Implementation for User Story 5

- [x] T071 [US5] Implement existing installation detection in `scripts/install.sh` (check `command -v beholder` and read manifest)
- [x] T072 [US5] Implement existing installation detection in `scripts/install.ps1`
- [x] T073 [US5] Implement version comparison in `scripts/install.sh` (compare installed vs latest, offer upgrade)
- [x] T074 [US5] Implement version comparison in `scripts/install.ps1`
- [x] T075 [US5] Add upgrade confirmation prompt in `scripts/install.sh` (default: yes in unattended mode)
- [x] T076 [US5] Add upgrade confirmation prompt in `scripts/install.ps1`
- [x] T077 [US5] Implement safe binary replacement (download to temp, validate, then overwrite) in `scripts/install.sh`
- [x] T078 [US5] Implement safe binary replacement in `scripts/install.ps1`
- [x] T079 [US5] Update installation manifest after upgrade in `scripts/install.sh`
- [x] T080 [US5] Update installation manifest after upgrade in `scripts/install.ps1`
- [ ] T081 [US5] Test upgrade scenario on Ubuntu 22.04 (v0.0.1-test → latest) per `quickstart.md` procedure
- [ ] T082 [US5] Test upgrade scenario on macOS per `quickstart.md` procedure
- [ ] T083 [US5] Test upgrade scenario on Windows per `quickstart.md` procedure
- [ ] T084 [US5] Test upgrade preserves user data (~/.beholder/config.yaml and database remain intact)

**Checkpoint**: All user stories should now be independently functional

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T085 [P] Update README.md with installation/uninstallation documentation
- [x] T086 [P] Create INSTALL.md with detailed installation guide (manual download, one-liner, troubleshooting)
- [x] T087 [P] Update agents.md with installer testing workflow
- [x] T088 Add checksum generation to GitHub Actions workflow (SHA256 for each binary)
- [ ] T089 Add checksum verification to installation scripts (optional: verify downloaded binary against published checksum)
- [ ] T090 Add error reporting documentation (common errors and resolutions)
- [ ] T091 Test complete end-to-end workflow per `quickstart.md` section 4 (fresh install → use → uninstall)
- [ ] T092 Test edge case: installation without sufficient disk space
- [ ] T093 Test edge case: installation with conflicting existing binary (manual install of same name)
- [ ] T094 Test edge case: uninstallation while beholder process is running
- [ ] T095 Create release checklist (workflow validation, manual testing, release notes)
- [ ] T096 [P] Code review of all installation/uninstallation scripts
- [ ] T097 Security review of script execution patterns (avoid code injection, validate inputs)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 3 (Phase 3)**: Depends on Foundational (Phase 2) - Must complete BEFORE User Story 1 (scripts need published binaries)
- **User Story 1 (Phase 4)**: Depends on User Story 3 (Phase 3) - Requires binaries to be available
- **User Story 2 (Phase 5)**: Depends on User Story 1 (Phase 4) - Requires installation to be functional first
- **User Story 4 (Phase 6)**: Depends on User Story 1 (Phase 4) - Enhances existing installation scripts
- **User Story 5 (Phase 7)**: Depends on User Story 1 (Phase 4) - Enhances existing installation scripts
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

**CRITICAL ORDER** (different from spec.md priority due to technical dependencies):

1. **User Story 3 (P2)**: GitHub Actions - MUST be first (produces binaries for other stories)
2. **User Story 1 (P1)**: Basic installation - Depends on US3 binaries
3. **User Story 2 (P2)**: Uninstallation - Depends on US1 installation working
4. **User Story 4 (P3)**: One-liner install - Enhances US1
5. **User Story 5 (P3)**: Upgrades - Enhances US1

### Within Each User Story

- Script skeletons before implementation details
- Detection/validation logic before actions (download, install, remove)
- Error handling after core functionality
- Testing after implementation
- Cross-platform testing required for each story completion

### Parallel Opportunities

**Within Phase 1 (Setup)**:
- All tasks are independent and can run in parallel

**Within Phase 2 (Foundational)**:
- T004 and T005 can run in parallel (both modify different files in cmd/beholder)
- T006 and T007 are documentation tasks and can run in parallel with code tasks

**Within Phase 3 (User Story 3 - GitHub Actions)**:
- T008-T015 are sequential (building up a single workflow file)
- T016 can run in parallel with T017 (local testing vs CI testing)

**Within Phase 4 (User Story 1 - Installation)**:
- T018 and T019 can run in parallel (different scripts)
- T020 and T021 can run in parallel (different scripts)
- T022 and T023 can run in parallel
- T024 and T025 can run in parallel
- T026 and T027 can run in parallel
- T028 and T029 can run in parallel
- T030 and T031 can run in parallel
- T032 and T033 can run in parallel
- T034 and T035 can run in parallel
- T036 and T037 can run in parallel
- T038, T039, T040 can run in parallel (different test environments)

**Within Phase 5 (User Story 2 - Uninstallation)**:
- T041 and T042 can run in parallel (different scripts)
- T043 and T044 can run in parallel
- T045 and T046 can run in parallel
- T047 and T048 can run in parallel
- T049 and T050 can run in parallel
- T051 and T052 can run in parallel
- T053 and T054 can run in parallel
- T055 and T056 can run in parallel
- T057, T058, T059, T060 can run in parallel (different test scenarios/environments)

**Within Phase 6 (User Story 4 - One-Liner)**:
- T062 and T063 can run in parallel (different scripts)
- T064 and T065 can run in parallel
- T066, T067, T068 can run in parallel (different test environments)
- T069 and T070 can run in parallel (different documentation files)

**Within Phase 7 (User Story 5 - Upgrades)**:
- T071 and T072 can run in parallel (different scripts)
- T073 and T074 can run in parallel
- T075 and T076 can run in parallel
- T077 and T078 can run in parallel
- T079 and T080 can run in parallel
- T081, T082, T083 can run in parallel (different test environments)

**Within Phase 8 (Polish)**:
- T085, T086, T087 can run in parallel (different documentation files)
- T092, T093, T094 can run in parallel (different edge case tests)
- T096 and T097 can run in parallel (code review and security review)

---

## Parallel Example: User Story 1 (Installation Scripts)

Assuming two developers are available:

**Developer A (POSIX):**
```bash
# Parallel batch 1: Script skeletons
T018  # Create scripts/install.sh

# Sequential: Implementation
T020  # Platform detection
T022  # API query
T024  # Download with retry
T026  # Validation
T028  # Installation
T030  # PATH modification
T032  # Manifest creation
T034  # User feedback
T036  # Error handling

# Parallel batch 2: Testing
T038  # Test on Ubuntu
T039  # Test on macOS
```

**Developer B (Windows):**
```bash
# Parallel batch 1: Script skeletons
T019  # Create scripts/install.ps1

# Sequential: Implementation
T021  # Platform detection
T023  # API query
T025  # Download with retry
T027  # Validation
T029  # Installation
T031  # PATH modification
T033  # Manifest creation
T035  # User feedback
T037  # Error handling

# Parallel batch 2: Testing
T040  # Test on Windows
```

Both developers work in parallel on their respective scripts, then coordinate for integration testing.

---

## Implementation Strategy

### MVP Scope (Minimum Deliverable)

**Goal**: Basic installation working on at least one platform

**Required Tasks**: T001-T007 (Setup + Foundational), T008-T017 (GitHub Actions), T018-T040 (Installation for at least macOS)

**Time Estimate**: ~3-5 days for single developer

### Incremental Delivery

**Iteration 1**: GitHub Actions + macOS installation (US3 + US1 for macOS only)
- Delivers: Automated builds + macOS users can install

**Iteration 2**: Windows + macOS installation (US1 complete)
- Delivers: All platforms can install

**Iteration 3**: Uninstallation (US2)
- Delivers: Complete install/uninstall lifecycle

**Iteration 4**: One-liner scripts (US4)
- Delivers: Improved developer experience

**Iteration 5**: Upgrades (US5)
- Delivers: Upgrade path for existing users

**Iteration 6**: Polish (Phase 8)
- Delivers: Production-ready with documentation and edge case handling

---

## Testing Validation Checklist

Per [quickstart.md](quickstart.md), validate each user story with:

- [ ] GitHub Actions workflow runs successfully and publishes all platform binaries (US3)
- [ ] Installation on macOS Intel (US1)
- [ ] Installation on macOS Apple Silicon (US1)
- [ ] Installation on Windows 11 (US1)
- [ ] Uninstallation with data preservation (US2)
- [ ] Uninstallation with data removal (US2)
- [ ] One-liner curl installation on macOS (US4)
- [ ] One-liner Invoke-WebRequest installation on Windows (US4)
- [ ] Upgrade from previous version preserving data (US5)
- [ ] All edge cases documented in spec.md (Phase 8)

---

**Next Steps**:
1. Begin with Phase 1 (Setup) to create directory structure
2. Complete Phase 2 (Foundational) to prepare codebase for versioning
3. Start Phase 3 (GitHub Actions) to enable binary distribution
4. Proceed through phases sequentially, testing each user story independently before moving to the next
