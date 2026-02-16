---
description: "Task list for feature implementation"
---

# Tasks: 設定ベースの画像管理とストレージ簡素化

**Input**: Design documents from `/specs/002-config-image-settings/`
**Prerequisites**: plan.md (required), spec.md (required)

**Tests**: 明示的なテスト要件は指定されていないため、テストタスクは含めない

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- All tasks include exact file paths

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: プロジェクト準備とバックアップ

- [X] T001 既存のcategories.goとsettings.goをバックアップ（削除前）
- [X] T002 既存のmigrate.goのスキーマ定義を確認・記録

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 全ユーザーストーリーで共通して必要な基盤変更

**⚠️ CRITICAL**: このフェーズ完了まで、ユーザーストーリーの作業を開始しない

- [X] T003 ImageConfig構造体をinternal/config/config.goに追加（MaxWidth, MaxFiles, SaveImages, Format）
- [X] T004 Config構造体にImage ImageConfig `yaml:"image"`フィールドを追加 in internal/config/config.go
- [X] T005 [P] 画像設定のバリデーションをinternal/config/validate.goに追加（max_width: 100-4096, max_files: >=0, format: jpeg|png）
  - エラーメッセージ形式: "image.max_width must be between 100 and 4096, got: %d", "image.format must be 'jpeg' or 'png', got: %s"
- [X] T006 [P] default.yamlにimageセクションを追加（デフォルト値: max_width=1280, max_files=0, save_images=true, format=jpeg）
- [X] T007 Event構造体のCategoryIDフィールドをCategoryNameに変更 in internal/storage/models.go
- [X] T008 Migrate()関数でcategories/settingsテーブル作成を削除、eventsテーブルのcategory_idをcategory_nameに変更 in internal/storage/migrate.go

**Checkpoint**: 基盤準備完了 - ユーザーストーリーの実装を並行開始可能

---

## Phase 3: User Story 3 - データベースを~/.beholderに保存 (Priority: P1) 🎯

**Goal**: データベースファイルを~/.beholder/beholder.dbに保存し、画像と同じディレクトリで一元管理

**Independent Test**: 初回起動時に`~/.beholder/beholder.db`が自動作成され、`beholder record`でデータが正常に保存される

### Implementation for User Story 3

- [X] T009 [US3] default.yamlのstorage.pathを~/.beholder/beholder.dbに変更 in internal/config/default.yaml
- [X] T010 [US3] Open()関数でstorage.pathが相対パスの場合~/.beholder/基準で解決するロジックを追加 in internal/storage/db.go
- [X] T011 [US3] 絶対パスの場合はそのまま使用する分岐処理を追加 in internal/storage/db.go

**Checkpoint**: User Story 3完了 - DBが~/.beholderに保存される

---

## Phase 4: User Story 1 - 設定ファイルで画像管理を制御 (Priority: P1) 🎯 MVP

**Goal**: config.yamlで画像の保存サイズ、最大保存枚数、保存ポリシーを設定でき、再起動時に反映される

**Independent Test**: config.yamlに画像設定（width: 800, max_files: 100）を追加し、`beholder record`を実行すると指定サイズの画像が保存され、古い画像が削除される

### Implementation for User Story 1

- [X] T012 [P] [US1] captureFullScreenPNG()でConfig.Image.MaxWidthを使用してsipsコマンドの-Zオプションを動的に設定 in internal/app/capture.go
- [X] T013 [P] [US1] Config.Image.Formatに基づいて画像フォーマット（jpeg/png）を選択する分岐処理を追加 in internal/app/capture.go
- [X] T014 [US1] Config.Image.SaveImagesがfalseの場合、画像ファイル保存をスキップする処理を追加 in internal/app/capture.go
  - 注: max_files=0は「無制限保存」を意味し、削除処理をスキップ（save_images=falseとは異なる動作）
- [X] T015 [US1] 画像保存後、cleanupOldImages()関数を実装（Config.Image.MaxFilesに基づきFIFO削除） in internal/app/capture.go
- [X] T016 [US1] cleanupOldImages()内でimgディレクトリのファイルをタイムスタンプソート、古い順に削除 in internal/app/capture.go
- [X] T017 [US1] 画像削除エラー時のログ出力とgraceful継続処理を追加 in internal/app/capture.go
  - ログレベル: log.Printf（警告レベル）を使用、削除権限エラーでも処理継続

**Checkpoint**: User Story 1完了 - 画像管理が設定ファイルで制御可能

---

## Phase 5: User Story 2 - データベーススキーマの簡素化 (Priority: P2)

**Goal**: categories/settingsテーブルを削除し、category_idをcategory_nameに変更してConfigを唯一の真実の源とする

**Independent Test**: 新規インストール後、`beholder record`を実行するとcategory_nameとしてカテゴリの日本語名が直接保存され、`beholder events`と`beholder summary`が正常動作する

### Implementation for User Story 2

- [X] T018 [P] [US2] InsertEvent()でCategoryIDの代わりにCategoryNameを保存する処理に変更 in internal/storage/events.go
- [X] T019 [P] [US2] ListEventsByDate()のクエリをcategory_idからcategory_nameに変更 in internal/storage/events.go
- [X] T020 [US2] RecordOnce()で分類結果のcategory IDをConfigから対応するNameに変換してEventに設定 in internal/app/record_once.go
- [X] T021 [US2] Generate()関数でcategoryMapパラメータを削除し、event.CategoryNameを直接使用 in internal/summary/generate.go
- [X] T022 [US2] summaryCmd()でcategoryMap生成を削除、Generate()呼び出しを更新 in cmd/beholder/cli.go
- [X] T023 [US2] UpsertCategories()関数をinternal/app/app.goから削除
- [X] T024 [US2] internal/storage/categories.goファイルを削除
- [X] T025 [US2] internal/storage/settings.goファイルを削除

**Checkpoint**: User Story 2完了 - スキーマ簡素化完了、コード削減達成

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 最終調整とドキュメント更新

- [ ] T026 [P] 設定検証のユニットテストを追加 in internal/config/validate_test.go（ImageConfig範囲チェック）
- [ ] T027 [P] storage層のユニットテストを追加 in internal/storage/storage_test.go（新スキーマでのInsert/List）
- [ ] T028 [P] 画像削除のintegrationテストを追加 in tests/integration/image_cleanup_test.go（max_files動作確認）
- [X] T029 READMEまたはagents.mdに新しい設定項目（imageセクション）を記載
- [X] T030 全コマンド（record, events, summary, start）の動作確認とエンドツーエンドテスト

---

## Dependencies & Execution Order

### Critical Path

1. **Phase 1 (Setup)** → **Phase 2 (Foundational)** → **Phase 3, 4, 5** (parallel) → **Phase 6 (Polish)**

### Between Phases

- **Phase 2 MUST complete before Phase 3, 4, 5**: Config/Model構造変更が全ての後続作業のベース
- **Phase 3, 4, 5 can run in parallel**: 異なるファイル、独立した機能

### Within Each User Story

- User Story 3: T009 → T010 → T011（順次）
- User Story 1: T012, T013 parallel → T014 → T015 → T016 → T017（部分並列）
- User Story 2: T018, T019 parallel → T020 → T021 → T022 → T023 → T024, T025 parallel（部分並列）

### Parallel Opportunities

- **Foundational**: T005, T006 can run in parallel with T003, T004
- **User Story 1**: T012, T013 can run in parallel
- **User Story 2**: T018, T019 can run in parallel; T024, T025 can run in parallel
- **Polish**: T026, T027, T028 can all run in parallel

---

## Implementation Strategy

### MVP First (User Story 3 + User Story 1)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 3 (DB path to ~/.beholder)
4. Complete Phase 4: User Story 1 (Image management config)
5. **STOP and VALIDATE**: Test US3 + US1 independently
6. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 3 → Test independently → Deploy/Demo
3. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
4. Add User Story 2 → Test independently → Deploy/Demo (Code simplification)
5. Polish phase → Final refinements

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 3 (DB path)
   - Developer B: User Story 1 (Image management)
   - Developer C: User Story 2 (Schema simplification)
3. Stories complete and integrate independently

---

## Notes

- **No migration needed**: まだリリース前のため、既存データ変換は不要
- **Code deletion**: categories.go, settings.go削除により約100行以上削減
- **Single Source of Truth**: Configが全ての設定情報を保持、DBテーブルとの二重管理を解消
- **Backward compatibility**: Not required（未リリースのため後方互換性不要）
