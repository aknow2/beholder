# Quickstart: スクリーンショット日報作成ツール

## Prerequisites
- Go 1.24
- macOS の Screen Recording 権限（画面キャプチャ用）
- Copilot SDK の認証情報（詳細は実装時に定義）

## Setup
1. 依存関係の取得
   - `go mod download`
2. 初期設定
   - 収集間隔（例: 5分）とカテゴリ一覧を登録

## Run (MVP)
- アプリ起動後、手動で1件記録を作成
- 自動収集をONにして一定間隔でイベントを作成
- 日報を指定日で生成（Markdown/JSON）

## Notes
- 画像はデフォルトで保存しない
- macOS では Screen Recording 権限が必要
