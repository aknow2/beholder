# Feature Specification: スクリーンショット日報作成ツール

**Feature Branch**: `001-screenshot-daily-report`  
**Created**: 2026-01-24  
**Status**: Draft  
**Input**: User description: "スクリーンショット日報作成ツールを開発したい（spec_summary参照） MVPとして実装したいので画面UI機能はまずは実装しない事とする"

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

### User Story 1 - CLIで手動記録する (Priority: P1)

ユーザーはCLIコマンドで「今の作業」を1回だけ記録でき、分類結果が保存される。

**Why this priority**: 最小の価値（記録が残ること）を最短で提供できるため。

**Independent Test**: CLIコマンドを1回実行するだけで1件の記録が作成され、一覧で確認できる。

**Acceptance Scenarios**:

1. **Given** 記録対象のカテゴリが設定済み、**When** ユーザーがCLIで「今記録する」を実行、**Then** 1件のイベントがカテゴリ付きで保存される
2. **Given** 分類が失敗した、**When** 保存が完了、**Then** 失敗状態として記録され再実行できる

---

### User Story 2 - 定期収集から日報を作る (Priority: P2)

ユーザーは一定間隔で自動記録を行い、指定日のサマリー（日報）を生成できる。

**Why this priority**: 連続的な記録と日報の作成が本機能の中心価値であるため。

**Independent Test**: 設定ファイルで収集間隔を設定し、CLIで収集を開始して一定数のイベントを作成し、日報出力を確認できる。

**Acceptance Scenarios**:

1. **Given** 自動収集が有効、**When** 設定した間隔が経過、**Then** 新しいイベントが作成される
2. **Given** 対象日のイベントが存在する、**When** 日報生成を実行、**Then** カテゴリ別合計と時系列が出力される

---

### User Story 3 - 設定ファイルでカテゴリと保存方針を管理する (Priority: P3)

ユーザーは設定ファイルでカテゴリの追加・編集・削除と、保存方針（画像を保存しない等）を変更できる。

**Why this priority**: 精度とプライバシーはユーザーごとに異なるため、後追いで調整できる必要がある。

**Independent Test**: 設定ファイル変更後にCLIを再起動し、変更が反映されることを確認できる。

**Acceptance Scenarios**:

1. **Given** 設定ファイルにカテゴリを追加、**When** 再起動、**Then** 次回以降の分類候補に含まれる
2. **Given** 画像保存が無効、**When** イベントが保存される、**Then** 画像データは永続保存されない

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

- 収集中にネットワークが不安定な場合、未分類イベントが保留される
- 1日分のイベントが0件の場合、空のサマリーが生成される
- 連続するスクリーンショットがほぼ同一の場合、同一カテゴリが連続して記録される
- 収集停止中に手動記録を行った場合でもイベントが保存される

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: システムは画面全体のスクリーンショットを取得できなければならない
- **FR-002**: システムは手動トリガーで1件の記録を作成できなければならない
- **FR-003**: システムは設定した間隔で自動的に記録を作成できなければならない
- **FR-004**: システムはカテゴリ一覧に基づき記録を1つのカテゴリに分類しなければならない
- **FR-005**: システムは各イベントにタイムスタンプと分類結果を保存しなければならない
- **FR-006**: システムは指定日のサマリーを生成し、カテゴリ別合計と時系列を出力しなければならない
- **FR-007**: ユーザーは設定ファイルでカテゴリを追加・編集・削除できなければならない
- **FR-008**: ユーザーは設定ファイルで画像を保存しない方針を選択でき、選択時は画像を永続保存しない
- **FR-009**: 分類失敗時は失敗状態として保存し、再実行できなければならない
- **FR-010**: オフライン時に記録を行う場合、未分類として保留し再開後に分類できなければならない
- **FR-011**: ユーザーは日報出力形式（Markdown・JSON・テキスト）を選択できなければならない
- **FR-012**: システムはCLIと設定ファイルのみで実行・管理できなければならない
- **FR-013**: システムはスクリーンショット画像を指定ディレクトリに保存し、パスをイベントに記録しなければならない
- **FR-014**: システムはデータベースファイル用ディレクトリが存在しない場合、自動的に作成しなければならない
- **FR-015**: システムは分類エラーが発生した場合、詳細なエラーログを出力しなければならない

### Key Entities *(include if feature involves data)*

- **Event**: 1回の記録単位。時刻、カテゴリ、分類信頼度、状態、補足情報を持つ
- **Category**: 作業分類の選択肢。名称、説明、例、表示色を持つ
- **Summary**: 1日分の集計結果。カテゴリ別合計、時系列、出力形式を持つ
- **Settings**: 取得間隔、保存方針、カテゴリ一覧を含む（設定ファイルで管理）

### Assumptions

- 画像はデフォルトで永続保存しない
- 収集間隔の初期値は5分を想定する
- サマリーの対象日はユーザーが選択する
- GUIは提供せず、CLIと設定ファイルのみで運用する

### Scope Boundaries

- 画面内容の完全な文字起こしや操作ログの再現は対象外
- 強制的な監視用途ではなく、個人の任意利用を前提とする
- 画像の長期保存はデフォルトでは行わない

### Dependencies

- 画面キャプチャのOS権限が必要
- 分類処理に外部サービスを利用する場合はネットワーク接続が必要

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: ユーザーは手動記録を30秒以内に完了できる
- **SC-002**: 1日分（最大1,000件）のイベントからサマリーを2分以内に生成できる
- **SC-003**: 主要なユーザータスク（記録開始・停止・日報生成）の初回完了率が90%以上
- **SC-004**: 日報のカテゴリ合計とタイムラインが95%以上のユーザーに有用と評価される
