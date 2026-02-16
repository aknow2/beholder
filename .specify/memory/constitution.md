<!--
Sync Impact Report:
Version: 1.0.1 → 1.0.2（パッチ：現行実装との整合とドキュメント反映）
Amended: 2026-01-31

Modified Principles:
- II. CLIファースト・インターフェース（README 反映）
- III. 設定駆動（変更なし）
- IV. ローカルファースト＆プライバシー（DB保存の表現を実装に整合）

Sections Added:
- なし

Sections Removed:
- なし

Templates Requiring Updates:
- ✅ .specify/templates/plan-template.md（参照の整理）
- ✅ .specify/templates/spec-template.md（変更不要）
- ✅ .specify/templates/tasks-template.md（変更不要）

Docs Requiring Updates:
- ✅ README.md（実装との差分修正）
- ✅ agents.md（実装との差分修正）

Follow-up TODOs:
- なし
-->

# Beholder 憲章

## コア原則

### I. モジュラーアーキテクチャ

**MUST**: コードは明確な境界を持つ内部パッケージに整理する。

- `internal/` 配下の各パッケージ（config, storage, classify, scheduler, summary, app）は単一責務
- パッケージ間のインターフェースは最小限で自己完結的
- 依存は内向きに流れる: app → services → storage/config
- `internal/` 間の循環依存は禁止

**Rationale**: 独立テスト、並行開発、関心の分離を可能にする。scheduler は storage に影響せず進化でき、classify は config を触らず差し替え可能。

### II. CLIファースト・インターフェース

**MUST**: すべての機能は CLI コマンドでテキスト I/O として提供する。

- コマンド形式: `beholder <command> [--flags]`
- 対応コマンド: `record`（`--oneshot` 付き）, `events`, `summary`, `reset`, `help`
- 標準出力に結果、標準エラーに失敗理由
- サマリー出力は text と Markdown を提供
- 終了コードは成功 0 / 失敗 1

**Rationale**: 自動化・スクリプト連携・他ツール統合を促進。CLI ファーストは明確なインターフェースとテスト可能な契約を強制する。

### III. 設定駆動

**MUST**: 実行時の挙動は YAML 設定ファイルで制御する。

- 単一の真実: `~/.beholder/config.yaml`
- 設定が存在しない場合は `internal/config/default.yaml` の埋め込みデフォルトで初期化
- ロード時にスキーマ検証（storage パス、copilot モデル、カテゴリ、画像制限）
- ストレージパス解決: 絶対パスはそのまま、相対パスは `~/.beholder/` 配下に解決
- ホットリロードは不要。変更は再起動で反映

**Rationale**: 方針（カテゴリ、間隔、モデル選択、画像保持）と実装（コード）を分離する。ユーザーはコード変更なしで挙動を調整できる。

### IV. ローカルファースト＆プライバシー

**MUST**: データはローカル保存を基本とし、保持制御を明示し、リモート利用は分類に限定する。

- DB はローカルに保存（SQLite）。デフォルトは `~/.beholder/beholder.db`（`storage.path` で変更可能）
- `image.save_images: true` の場合、画像は `~/.beholder/imgs` に保存
- `image.save_images: false` の場合、画像は一時ディレクトリに保存し分類後に削除
- `image.max_files` を設定（>0）した場合は FIFO で古い画像を削除
- 分類は GitHub Copilot SDK を使用しネットワークが必要。失敗時も `FAILED` でイベントを記録

**Rationale**: スクリーンショットはプライバシーに敏感。ローカル保存でユーザーの
制御を維持し、リモート分類は分離して失敗時もデータ欠落なく劣化運用する。

### V. 段階的強化

**MUST**: ユーザーストーリーの優先度順（P1 → P2 → P3）に実装する。

- P1（MVP）: CLI による手動記録と分類
- P2: スケジューリングと日次サマリー
- P3: 追加の設定と保持制御
- 各優先度は独立してデプロイ・テスト可能

**Rationale**: 価値を段階的に提供する。P1 で即時有用性を提供し、P2 で自動化、P3 で制御性を強化。ユーザーは段階的に採用できる。

## 技術標準

### 言語とツール

- **Go 1.24**: コアロジックは標準ライブラリ中心
- **依存関係**: 外部依存は最小化。十分なら stdlib を優先
- **DB**: SQLite（modernc.org/sqlite、pure Go / no cgo）
- **分類**: GitHub Copilot SDK（`github.com/github/copilot-sdk/go`）
- **スクリーンショット**: macOS の `screencapture` + `sips`
- **テスト**: 単体テストは `go test`

### コード品質

- **整形**: `gofmt` 準拠必須
- **エラー処理**: 明示的に返却し、文脈付きでラップ（`fmt.Errorf`）
- **ロギング**: エラーは stderr、情報は stdout
- **バリデーション**: config 読み込みと CLI 引数で入力検証

### パフォーマンス目標

- リサイズ後の画像は 3 MiB 以下（超過時はキャプチャ失敗）
- 画像幅は `image.max_width`（100–4096）で制約
- スケジューラは分単位（デフォルト 10）でリアルタイム性は要求しない

## 開発ワークフロー

### 機能実装

1. **Spec-First**: `specs/###-feature-name/spec.md` に仕様を書く
2. **Planning**: plan.md を生成し技術コンテキストと構成を明記
3. **Task Breakdown**: tasks.md をユーザーストーリー優先度で整理
4. **Incremental Delivery**: P1 → テスト → P2 → テスト → P3 → テスト
5. **Testing**: 新規パッケージは unit テスト、ワークフローは integration テスト

### ユーザーストーリー構成

- **Phase 1**: セットアップ（ディレクトリ、config テンプレ、.gitignore）
- **Phase 2**: 基盤（config ローダ、storage、app 初期化）
- **Phase 3+**: ユーザーストーリーごとに 1 フェーズ（独立テスト可能）
- **Final Phase**: 仕上げと横断的改善

### 憲章準拠

- **実装前**: 計画が原則に整合するか検証
- **実装中**: パッケージ境界と CLI 契約を維持
- **実装後**: 設定駆動とローカルファースト動作を検証

## ガバナンス

本憲章はアドホックな決定に優先する。原則違反は plan.md の Complexity Tracking で明示的に正当化する。

### 改定手順

1. 変更案と根拠、影響分析を提示
2. セマンティックバージョンで憲章を更新:
   - **MAJOR**: 後方互換性のない原則の削除/再定義
   - **MINOR**: 新規原則や重要なセクション追加
   - **PATCH**: 明確化、文言修正、非意味的な更新
3. テンプレート（plan, tasks, spec）と運用ガイダンス（agents.md）へ反映
4. Sync Impact Report（ファイル先頭の HTML コメント）へ記録

### 準拠レビュー

- すべての計画に Constitution Check セクションを含める
- タスク分解はユーザーストーリー優先度（原則 V）に整合
- コードレビューでパッケージ分離（原則 I）と CLI 契約（原則 II）を確認

### ランタイムガイダンス

開発コマンドや運用フローは [agents.md](../../agents.md) を参照。

**Version**: 1.0.2 | **Ratified**: 2026-01-26 | **Last Amended**: 2026-01-31
