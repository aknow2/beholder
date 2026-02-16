# Research: スクリーンショット日報作成ツール

## Decision 1: 画面キャプチャ実装
- **Decision**: `github.com/kbinani/screenshot` を採用
- **Rationale**: マルチディスプレイ対応・APIが簡潔・実績が多い。macOSでのフルスクリーン取得に適合。
- **Alternatives considered**:
  - `github.com/go-vgo/robotgo`: 依存が重くオーバースペック
  - `github.com/vova616/screenshot`: メンテナンスが弱い
  - `screencapture` コマンド: OS依存が強くエラー制御が難しい

## Decision 2: ローカルDB
- **Decision**: `modernc.org/sqlite` を採用（pure Go）
- **Rationale**: CGO不要でビルド・配布が容易。Dockerやクロスビルドとの相性が良い。
- **Alternatives considered**:
  - `github.com/mattn/go-sqlite3`: 性能と互換性は高いがCGO必須
  - JSONファイル保存: 検索/集計や整合性が弱い

## Decision 3: 実行管理はCLI/設定ファイル
- **Decision**: GUI/トレイUIを採用せず、CLIと設定ファイルのみで実行・管理する
- **Rationale**: 依存を最小化し、自動化・運用が容易。ヘッドレス運用にも対応。
- **Alternatives considered**:
  - トレイ/メニューバーUI: OS依存とテストコストが高い
  - フルGUI（Fyne/Wails）: MVPには過剰で配布コストが高い

## Decision 4: Copilot分類の統合
- **Decision**: `github.com/github/copilot-sdk/go` を利用
- **Rationale**: 仕様要件に一致し、SDKの公式サンプルがある。
- **Alternatives considered**:
  - 独自API実装: 運用コストが高い
  - 外部LLMの直接呼び出し: 権限とプロンプト設計が必要

## Decision 5: 失敗時の扱い
- **Decision**: `PENDING/FAILED/OK` の状態をイベントに保持し再分類可能にする
- **Rationale**: オフライン・失敗時の再実行要件に一致。
- **Alternatives considered**:
  - 失敗イベントを破棄: 日報の欠落につながる

## Open Questions → Resolved
- **技術スタックが不明**: 既存リポジトリはGo 1.24を使用。Goのデスクトップ常駐アプリとして進める。
- **ターゲットOS**: macOSを主対象（開発環境がmacOS）。他OSは将来対応とする。
