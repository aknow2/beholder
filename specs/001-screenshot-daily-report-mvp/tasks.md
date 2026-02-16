---
description: "Task list for feature implementation"
---

# Tasks: スクリーンショット日報作成ツール

**Input**: Design documents from `/specs/001-screenshot-daily-report/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/, quickstart.md
**Tests**: 明示的なテスト要件は指定されていないため、テストタスクは含めない

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- All tasks include exact file paths

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create project directories per plan in cmd/beholder/, internal/{capture,classify,scheduler,storage,summary,config,app}/, configs/, tests/{unit,integration}/
- [X] T002 Add default configuration template in configs/default.yaml
- [X] T003 [P] Add local data directory to .gitignore in .gitignore

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure required before any user story work

- [X] T004 Add config models and loader in internal/config/config.go
- [X] T005 Add config validation in internal/config/validate.go
- [X] T006 Add SQLite connection + migration bootstrap in internal/storage/db.go and internal/storage/migrate.go
- [X] T007 [P] Define storage models in internal/storage/models.go
- [X] T008 Implement repositories for events/categories/settings in internal/storage/events.go, internal/storage/categories.go, internal/storage/settings.go
- [X] T009 Add app initialization (config load, DB open, migrate) in internal/app/app.go

**Checkpoint**: Foundation ready - user story implementation can begin

---

## Phase 3: User Story 1 - CLIで手動記録する (Priority: P1) 🎯 MVP

**Goal**: CLIで1件の記録を作成し、分類結果を保存できる

**Independent Test**: `beholder record` を実行してイベントが保存されることを確認

### Implementation for User Story 1

- [X] T010 [P] [US1] Implement full-screen capture in internal/app/capture.go (macOS screencapture + sips resize)
- [X] T011 [P] [US1] Implement Copilot classify client in internal/classify/client.go (Attachments API)
- [X] T012 [US1] Implement record-once flow in internal/app/record_once.go (with error logging)
- [X] T013 [US1] Add CLI command for manual record in cmd/beholder/cli.go
- [X] T014 [US1] Add CLI command to list events in cmd/beholder/cli.go
- [X] T015 [P] [US1] Add auto-create database directory in internal/storage/db_v2.go
- [X] T016 [P] [US1] Add image format support (JPEG/PNG) in internal/app/capture.go
- [X] T017 [US1] Implement screenshot file save to ~/.beholder/imgs/ in internal/app/capture.go

**Checkpoint**: User Story 1 fully functional and independently testable

---

## Phase 4: User Story 2 - 定期収集から日報を作る (Priority: P2)

**Goal**: 定期収集を開始/停止でき、指定日のサマリーを生成できる

**Independent Test**: `beholder start` → 収集後に `beholder summary --date YYYY-MM-DD` を実行

### Implementation for User Story 2

- [X] T018 [P] [US2] Implement scheduler loop in internal/scheduler/scheduler.go
- [X] T019 [US2] Wire scheduler run/stop in internal/app/daemon.go
- [X] T020 [P] [US2] Implement summary generator in internal/summary/generate.go
- [X] T021 [US2] Add CLI command for summary output in cmd/beholder/cli.go
- [X] T022 [US2] Add CLI command to start scheduler in cmd/beholder/cli.go

**Checkpoint**: User Story 2 independently functional

---

## Phase 5: User Story 3 - 設定ファイルでカテゴリと保存方針を管理する (Priority: P3)

**Goal**: 設定ファイルでカテゴリ/保存方針を管理し、再起動で反映できる

**Independent Test**: 設定ファイルを変更→再起動→分類候補と保存方針が反映されることを確認

### Implementation for User Story 3

- [ ] T023 [P] [US3] Implement config sync (categories/settings upsert) in internal/app/config_sync.go
- [ ] T024 [US3] Enforce saveImages policy in internal/app/record_once.go
- [ ] T025 [US3] Add CLI command to validate config in cmd/beholder/main.go

**Checkpoint**: User Story 3 independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements affecting multiple stories

- [ ] T026 [P] Update quickstart CLI usage in specs/001-screenshot-daily-report/quickstart.md
- [ ] T027 Add structured logging helpers in internal/app/logging.go
- [ ] T028 Run quickstart verification checklist in specs/001-screenshot-daily-report/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **User Stories (Phase 3–5)** → **Polish (Phase 6)**

### User Story Dependencies

- **US1** starts after Phase 2
- **US2** starts after Phase 2
- **US3** starts after Phase 2

### Within Each User Story

- Core modules before CLI wiring
- Storage/config dependencies before application flows

---

## Parallel Execution Examples

### User Story 1

- T010 and T011 can run in parallel (capture vs classify)
- T015, T016, T017 can run in parallel (DB directory vs image format vs file save)

### User Story 2

- T018 and T020 can run in parallel (scheduler vs summary)

### User Story 3

- T023 can run in parallel with documentation updates in Phase 6

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1 → Phase 2
2. Implement US1 (T010–T014)
3. Validate manual record flow

### Incremental Delivery

1. Add US2 (T015–T019)
2. Add US3 (T020–T022)
3. Polish (T023–T025)
