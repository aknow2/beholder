# Feature Specification: 設定ベースの画像管理とストレージ簡素化

**Feature Branch**: `002-config-image-settings`  
**Created**: 2026-01-27  
**Status**: Draft  
**Input**: User description: "ローカルに画像のサイズや最大保存枚数をConfigから指定出来る。CategoryやSettingsのテーブルは不要です。Configを参照すれば分かるので、Eventに保存するカテゴリは文字列で保存すればOK。dbのデータを~/.beholderに保存する"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 設定ファイルで画像管理を制御 (Priority: P1)

ユーザーは`~/.beholder/config.yaml`で画像の保存サイズ、最大保存枚数、保存ポリシーを設定でき、アプリケーション再起動時に設定が反映される。

**Why this priority**: 現在はハードコードされた値（1280px、無制限保存）を使用しているため、ユーザーがディスク使用量や画像品質を制御できない。設定ファイルでの制御は最も基本的な要件。

**Independent Test**: config.yamlに画像設定（width: 800, max_files: 100）を追加し、`beholder record`を実行すると、指定サイズの画像が保存され、古い画像が削除されることを確認。

**Acceptance Scenarios**:

1. **Given** config.yamlにimage.max_width: 800を設定、**When** スクリーンショット取得、**Then** 画像の幅が800px以下で保存される
2. **Given** config.yamlにimage.max_files: 50を設定し51枚目を保存、**When** 保存処理実行、**Then** 最も古い画像が自動削除され50枚のみ残る
3. **Given** config.yamlにimage.save_images: falseを設定、**When** 記録実行、**Then** 分類結果のみ保存され画像ファイルは作成されない
4. **Given** config.yamlでimage.formatをjpegに設定、**When** スクリーンショット取得、**Then** JPEG形式で保存される

---

### User Story 2 - データベーススキーマの簡素化 (Priority: P2)

開発者はcategoriesテーブルとsettingsテーブルを削除し、Eventテーブルのcategory_idをcategory_name（文字列）に変更することで、データベースとConfigの二重管理を解消できる。

**Why this priority**: 現在のcategoriesテーブルとConfig.categoriesは同じ情報を重複管理しており、同期処理（UpsertCategories）が必要。Configを唯一の真実の源（Single Source of Truth）にすることで、コードが簡潔になりバグが減る。

**Independent Test**: 新規インストール後、`beholder record`を実行すると、category_nameとしてカテゴリの日本語名が直接保存され、`beholder events`と`beholder summary`が正常動作することを確認。

**Acceptance Scenarios**:

1. **Given** 新スキーマのデータベース、**When** 記録実行、**Then** Eventレコードにcategory_nameとして「実装」「調査」などの日本語名が直接保存される
2. **Given** categoriesテーブルとsettingsテーブルが存在しない、**When** アプリケーション起動、**Then** エラーなく起動し、Configからカテゴリ情報を読み込む
3. **Given** Configでカテゴリ名を変更（「実装」→「開発作業」）、**When** アプリケーション再起動、**Then** 新規記録時に「開発作業」として保存される（過去データは変更されない）

---

### User Story 3 - データベースを~/.beholderに保存 (Priority: P1)

データベースファイルは`~/.beholder/beholder.db`に保存され、画像と同じディレクトリで一元管理できる。

**Why this priority**: 現在はカレントディレクトリの`data/beholder.db`に保存されており、ユーザーのホームディレクトリ配下に統一することで、データの場所が明確になり管理しやすくなる。

**Independent Test**: 初回起動時に`~/.beholder/beholder.db`が自動作成され、`beholder record`でデータが正常に保存されることを確認。

**Acceptance Scenarios**:

1. **Given** `~/.beholder/`ディレクトリが存在しない、**When** アプリケーション初回起動、**Then** `~/.beholder/beholder.db`が自動作成される
2. **Given** データベースが`~/.beholder/beholder.db`に存在、**When** 記録実行、**Then** データが正常に保存される
3. **Given** 設定ファイルでDBパスをカスタマイズ可能、**When** storage.path: custom/path.dbを設定、**Then** 指定パスにDBが作成される

---

### Edge Cases

- **画像サイズが0または負の値**: config.yamlでmax_widthに無効な値を設定した場合、バリデーションエラーを返す
- **max_filesが0**: 画像保存が無効化される（save_images: falseと同じ動作）
- **画像ディレクトリの削除権限がない**: 古い画像削除時にエラーが発生した場合、ログに記録するが処理は継続
- **既存の画像ファイルが手動削除されている**: max_files計算時にファイル数の不整合があっても処理を継続
- **~/.beholderディレクトリの作成権限がない**: 適切なエラーメッセージを返してアプリケーションを終了

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Configに`image`セクションを追加し、以下の設定をサポートする:
  - `max_width`: 画像の最大幅（ピクセル、デフォルト: 1280）
  - `max_files`: 最大保存枚数（デフォルト: 0=無制限）
  - `save_images`: 画像保存の有効/無効（デフォルト: true）
  - `format`: 画像フォーマット（"jpeg"または"png"、デフォルト: "jpeg"）

- **FR-002**: 画像保存時に`image.max_width`を超える場合、アスペクト比を維持して縮小する

- **FR-003**: `image.max_files`が設定されている場合、保存枚数が上限を超えたら最も古いファイルから削除する（FIFO方式）

- **FR-004**: `image.save_images`がfalseの場合、画像ファイルを保存せず分類結果のみDBに記録する

- **FR-005**: Eventテーブルの`category_id`カラムを`category_name`（TEXT型）に変更し、カテゴリの日本語名を直接保存する

- **FR-006**: `categories`テーブルと`settings`テーブルを削除する

- **FR-007**: Config検証時に`image`セクションの値が有効範囲内であることを確認する:
  - `max_width`: 100～4096の範囲
  - `max_files`: 0以上の整数
  - `format`: "jpeg"または"png"のみ

- **FR-008**: 画像削除処理はファイルシステムのみに作用し、DBレコードには影響しない（Event.screenshot_hashは残る）

- **FR-009**: Summaryコマンドはcategory_nameを使用してカテゴリ別集計を実行する（Configとの照合は不要）

- **FR-010**: デフォルトのデータベースパスを`~/.beholder/beholder.db`に変更する

- **FR-011**: storage.pathが相対パスの場合、`~/.beholder/`を基準に解決する（例: `path: data.db` → `~/.beholder/data.db`）
  - チルダ（~）で始まるパスはホームディレクトリに展開後、絶対パスとして扱う

- **FR-012**: storage.pathが絶対パスの場合、そのまま使用する

### Key Entities

- **ImageConfig**: 画像管理設定（max_width, max_files, save_images, format）
- **Event**: イベント記録（category_nameフィールドに変更、category_id削除）

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: ユーザーは設定ファイル編集のみで画像サイズと保存枚数を変更でき、コード変更不要（設定変更→再起動で反映）

- **SC-002**: データベースサイズが削減される（categoriesテーブル削除により約10-20行のレコード削減）

- **SC-003**: コードの行数が削減される（UpsertCategories関数、categories.go、settings.go削除により約100行以上削減）

- **SC-004**: 設定検証が追加され、無効な画像設定でアプリケーション起動時にエラーが返される（ユーザーに即座にフィードバック）

- **SC-005**: すべてのデータが`~/.beholder/`ディレクトリに集約され、データの場所が明確になる

## Assumptions

- まだリリース前のため、既存データのマイグレーションは不要（新規スキーマのみ実装）
- 画像削除は非同期処理不要（記録時に同期的に削除してよい）
- 画像ファイル名のタイムスタンプベースのソートで古さを判定できる（ファイル作成日時は使用しない）
  - ファイル名からタイムスタンプが解析不可能な場合、警告ログを出力してスキップ（削除対象外）
- 複数プロセスからの同時記録は想定しない（単一ユーザー・単一プロセス利用）
  - ファイルロックは実装しない（単一プロセス前提のため並行削除の安全性確保は不要）

## Out of Scope

- 既存データのマイグレーション（まだリリースしていないため不要）
- 画像圧縮品質の設定（JPEGのquality設定は今回対象外、デフォルト品質を使用）
- 画像の暗号化や署名（セキュリティ強化は別フィーチャー）
- カテゴリ名の履歴管理（過去のEventレコードのcategory_nameは更新しない）
- 画像ファイルのアーカイブ機能（削除ではなく別ディレクトリに移動する機能は対象外）
- 動的なカテゴリ追加UI（設定ファイル手動編集のみサポート）

