# agents.md

## 概要
このリポジトリ（beholder）で開発を始めるための基本操作をまとめたガイドです。ローカルGoでのビルド、テスト、開発フローなどを記載しています。

## 必須環境
- Go 1.24（`go.mod` に準拠）
- Git

## セットアップ
1. リポジトリをクローン

```bash
git clone https://github.com/aknow2/beholder.git
cd beholder
```

2. 依存関係のダウンロード

```bash
go mod download
```

## ビルド
- バイナリをビルド

```bash
go build -o bin/beholder ./cmd/beholder
```

## 設定

デフォルトでは `~/.beholder/config.yaml` が自動作成されます。初回実行時に以下のファイル・ディレクトリが作成されます：

- `~/.beholder/config.yaml` - 設定ファイル
- `~/.beholder/beholder.db` - SQLiteデータベース
- `~/.beholder/imgs/` - スクリーンショット保存ディレクトリ

### 画像管理設定（image セクション）

```yaml
image:
  max_width: 1280      # 画像の最大幅（ピクセル、100-4096の範囲）
  max_files: 0         # 最大保存枚数（0=無制限）
  save_images: true    # 画像保存の有効/無効
  format: jpeg         # 画像フォーマット（jpeg または png）
```

- **max_width**: スクリーンショットを指定幅にリサイズ（アスペクト比維持）
- **max_files**: 0より大きい値を設定すると、古い画像から自動削除（FIFO）
- **save_images**: `false`にすると画像ファイルを保存せず分類結果のみDB記録
- **format**: JPEG（容量小）またはPNG（可逆圧縮）を選択可能

## 実行
- `go run` で直接実行

```bash
go run ./cmd/beholder record --oneshot
```

- ビルド済みバイナリを実行

```bash
./bin/beholder record --oneshot
```

## テスト
- テストを実行

```bash
go test ./... -v
```

## フォーマットとリンティング
- フォーマット

```bash
gofmt -w .
```

- 簡易チェック

```bash
go vet ./...
```

（プロジェクトで `golangci-lint` 等を導入する場合はそのコマンドを追加してください）

## デバッグ
- `dlv` を使ったデバッグ

```bash
dlv debug ./cmd/beholder
```

- 簡易的には `go run` にログを追加して確認します。

## インストーラーテスト

### ローカルテスト（手動ビルド）

```bash
# クロスコンパイル（macOS Intel）
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(git describe --tags --always)" -o dist/beholder-darwin-amd64 ./cmd/beholder

# クロスコンパイル（macOS Apple Silicon）
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$(git describe --tags --always)" -o dist/beholder-darwin-arm64 ./cmd/beholder

# クロスコンパイル（Windows）
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$(git describe --tags --always)" -o dist/beholder-windows-amd64.exe ./cmd/beholder

# バージョン確認（Linux上で実行不可のため、ビルドのみ確認）
ls -lh dist/
```

### GitHub Actions経由のリリース

```bash
# テストリリースを作成
git tag -a v0.1.0-test -m "Test release"
git push origin v0.1.0-test

# ワークフロー確認
# https://github.com/aknow2/beholder/actions

# リリース確認
# https://github.com/aknow2/beholder/releases
```

### インストールスクリプトのテスト

**注意**: 実際のGitHub Releaseが存在する必要があります。

```bash
# macOS（実行には macOS 環境が必要）
./scripts/install.sh

# Windows（実行には Windows 環境が必要）
# PowerShell で実行:
# .\scripts\install.ps1

# アンインストールテスト
./scripts/uninstall.sh  # macOS
# .\scripts\uninstall.ps1  # Windows
```

### インストールスクリプトのドライラン

実際にインストールせずにスクリプトの動作を確認：

```bash
# プラットフォーム検出のみテスト
bash -c 'OS=$(uname -s | tr "[:upper:]" "[:lower:]"); ARCH=$(uname -m); case "$ARCH" in x86_64) ARCH="amd64";; arm64|aarch64) ARCH="arm64";; esac; echo "Platform: $OS-$ARCH"'

# GitHub API 確認
curl -fsSL https://api.github.com/repos/aknow2/beholder/releases/latest | jq -r '.tag_name'
```

## よく使うコマンド
- `go run ./cmd/beholder` : アプリケーション実行
- `go build -o bin/beholder ./cmd/beholder` : バイナリビルド
- `go test ./... -v` : テスト実行
- `go mod tidy` : 依存関係の整理

---
このファイルに追加してほしいコマンドやワークフローがあれば教えてください。
