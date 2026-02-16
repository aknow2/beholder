# Beholder

スクリーンショットを取得して分類し、イベントとしてローカルに記録するCLIアプリケーションです。設定ファイルで画像サイズや保存枚数などを調整できます。

## 特徴
- CLIで記録・一覧・サマリー生成・スケジューラ起動
- 設定ファイル（~/.beholder/config.yaml）で挙動を制御
- 画像とDBはローカル（~/.beholder/）に保存

## インストール

### macOS

```bash
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.sh | sh
```

または、[GitHub Releases](https://github.com/aknow2/beholder/releases/latest)から手動でダウンロード：

```bash
# Intel Mac の場合
curl -LO https://github.com/aknow2/beholder/releases/download/v{VERSION}/beholder-v{VERSION}-darwin-amd64
chmod +x beholder-v{VERSION}-darwin-amd64
mv beholder-v{VERSION}-darwin-amd64 ~/.local/bin/beholder

# Apple Silicon (M1/M2) の場合
curl -LO https://github.com/aknow2/beholder/releases/download/v{VERSION}/beholder-v{VERSION}-darwin-arm64
chmod +x beholder-v{VERSION}-darwin-arm64
mv beholder-v{VERSION}-darwin-arm64 ~/.local/bin/beholder
```

### Windows

PowerShellで実行：

```powershell
Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/main/scripts/install.ps1 -UseBasicParsing | Invoke-Expression
```

または、[GitHub Releases](https://github.com/aknow2/beholder/releases/latest)から手動でダウンロードして `%USERPROFILE%\.beholder\bin` に配置します。

### インストール確認

```bash
beholder --version
```

### アンインストール

**macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/aknow2/beholder/main/scripts/uninstall.sh | sh
```

**Windows:**

```powershell
Invoke-WebRequest https://raw.githubusercontent.com/aknow2/beholder/main/scripts/uninstall.ps1 -UseBasicParsing | Invoke-Expression
```

アンインストール時にユーザーデータ（~/.beholder/）を削除するか保持するか選択できます。

---

## 開発者向けセットアップ

### 必須環境
- Go 1.24（go.mod準拠）
- Git

### セットアップ

```bash
go mod download
```

### ビルド

```bash
go build -o bin/beholder ./cmd/beholder
```

---

## クイックスタート

```bash
./bin/beholder record --oneshot
```

初回起動時に ~/.beholder/config.yaml が自動作成されます。

## コマンド

- `record` : スケジューラ起動（interval_minutes 間隔で記録を繰り返す）
  - `--oneshot` を付けると1回だけ記録して終了
- `events --date <YYYY-MM-DD>` : 指定日のイベント一覧
- `summary --date <YYYY-MM-DD> --format <text|markdown>` : 日次サマリー生成
- `reset --date <YYYY-MM-DD>` : 指定日のイベント削除（確認プロンプトあり）

例:

```bash
./bin/beholder events --date 2026-01-28
./bin/beholder summary --date 2026-01-28 --format markdown
./bin/beholder record
./bin/beholder reset --date 2026-01-28
```

## 設定

デフォルト設定は埋め込みの [internal/config/default.yaml](internal/config/default.yaml) です。設定ファイルが存在しない場合、自動生成されます。

主要項目:

```yaml
storage:
  path: ~/.beholder/beholder.db

scheduler:
  interval_minutes: 10

copilot:
  model: gpt-4.1

image:
  max_width: 1280
  max_files: 0
  save_images: true
  format: jpeg
```

- `storage.path` は相対パスの場合 ~/.beholder/ 基準で解決されます。
- `image.max_files` が0の場合は無制限です。
- `image.save_images: false` で画像ファイルを保存せず分類結果のみ記録します。

## 実行（go run）

```bash
go run ./cmd/beholder record
```

## テスト

```bash
go test ./... -v
```

## 開発ガイド

開発用コマンドや運用メモは [agents.md](agents.md) を参照してください。
