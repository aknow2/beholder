# スクリーンショット日報作成ツール 仕様書（GitHub Copilot SDK）

## 1. 概要

PCの画面全体スクリーンショットを定期的に取得し、画像内容を **GitHub Copilot SDK** を用いたエージェントに解析させ、事前定義した「作業カテゴリ（選択肢）」のどれに該当するかを判定・保存する。保存された判定結果を集計して、1日の作業内容が分かるサマリー（日報）を生成する。

## 2. 目的

* 手動入力なし/最小入力で日報を自動生成する
* 「いつ」「何をしていたか」をカテゴリベースで追跡できる
* 後から検索・集計（週次/プロジェクト別）できる

## 3. 非目的（今回スコープ外）

* 画面内容の完全な文字起こし（OCR）や、詳細な操作ログの再現
* 監視/労務管理用途の強制運用（個人用/任意運用を想定）
* 画像そのものを長期保存すること（デフォルトは保存しない/最小化）

## 4. 想定ユーザー

* 個人開発者 / エンジニア / クリエイター
* 「その日何をしていたか忘れがち」な人

## 5. 用語

* **イベント（Event）**: 1回のスクショ取得→分類→保存までの単位
* **カテゴリ（Category）**: 事前定義された選択肢（例: 実装、調査、会議…）
* **サマリー（Summary）**: 1日分のイベントを集計し、日報として整形した結果
* **Copilot Agent**: Copilot SDKで動く画像分類担当のエージェント

---

## 6. ユースケース / ユーザーストーリー

1. ユーザーは常駐アプリを起動し、スクショ取得間隔（例: 5分）とカテゴリ一覧を設定する
2. アプリは定期的に画面全体のスクショを取得する
3. アプリはスクショをCopilot Agentに渡し、カテゴリのいずれかへ分類させる
4. 分類結果（カテゴリ、信頼度、推定アプリ/ウィンドウなど）をタイムスタンプ付きでDBに保存する
5. ユーザーが「Summary」を実行すると、当日の活動がカテゴリ別・時系列で表示され、日報テキストとして出力できる

---

## 7. 機能要件

### 7.1 スクリーンショット取得

* 画面全体（フルスクリーン）を取得できる
* 取得トリガー:

  * 自動: 設定間隔（例: 1/3/5/10/15分）
  * 手動: いま記録する（ワンクリック/ショートカット）
* 画像フォーマット:

  * デフォルト: PNG
  * オプション: JPEG（サイズ削減）
* 取得メタデータ:

  * capturedAt（UTC推奨）
  * displayCount（複数モニタの場合）
  * resolution（W×H）

### 7.2 Copilot Agentへの解析依頼

* 入力:

  * スクリーンショット画像（バイナリ or Base64）
  * カテゴリ一覧（固定/ユーザー設定）
  * 追加ヒント（任意）: プロジェクト名一覧、よく使うツール一覧、勤務時間帯
* 出力（必須）:

  * selectedCategoryId
  * confidence（0〜1）
  * rationale（短い理由。ユーザー表示は任意）
* 出力（任意）:

  * detectedApps（例: VS Code, Figma, Slack）
  * detectedKeywords（例: issue番号、リポジトリ名）

### 7.3 カテゴリ（選択肢）管理

* ユーザーがカテゴリを追加/編集/削除できる
* カテゴリには以下を持てる:

  * id
  * name（表示名）
  * description（判定の補助）
  * examples（例: 「PRレビュー」「docs読む」）
  * color（UI表示用、任意）

### 7.4 保存（DB）

* 各イベントごとに、分類結果をDBへ保存する
* 画像はデフォルトでは保存しない

  * 代替として「画像ハッシュ」「縮小サムネ（ぼかし）」「マスク済みサムネ」などを選べる
* 保存する推奨フィールド:

  * eventId
  * capturedAt
  * categoryId
  * confidence
  * agentVersion
  * screenshotHash（任意）
  * detectedApps（任意）
  * notes（ユーザーが後から補足できる）

### 7.5 Summary（日報生成）

* 対象日（YYYY-MM-DD）を指定して生成できる
* 出力形式:

  * 画面表示
  * Markdown（デフォルト）
  * JSON（エクスポート用）
* サマリー内容（最低限）:

  * カテゴリ別の合計時間（イベント間隔×件数）
  * 時系列のタイムライン（いつ何カテゴリだったか）
  * 重要っぽい変化点（カテゴリの切替が多い時間帯など）
* 追加（任意）:

  * 「今日やったこと（箇条書き）」を自然言語で整形（Copilot Agentで再要約）

---

## 8. 非機能要件

### 8.1 プライバシー / セキュリティ

* デフォルトは **画像を永続保存しない**
* 画像を外部に送る場合は、ユーザーに明示（初回同意）
* マスキング（任意機能）:

  * 画面の特定領域（座標指定）を黒塗り
  * あるアプリ検出時は送信しない（例: パスワード管理、銀行）
* ローカル保存する場合:

  * DB暗号化（OSキーチェーン利用推奨）
  * 画像は暗号化して保存（オプション）

### 8.2 性能

* 1イベントあたりの処理時間目標:

  * スクショ取得: 1秒以内
  * エージェント分類: 5〜15秒以内（ネットワーク依存）
* オフライン時:

  * スクショは取得する/しないを設定可能
  * 取得する場合は「未分類キュー」に積み、オンライン復帰で分類

### 8.3 信頼性

* Copilot API失敗時は:

  * リトライ（指数バックオフ、最大N回）
  * 失敗イベントとして保存（status=FAILED）し後で再実行できる

---

## 9. 画面 / UI（最小）

### 9.1 常駐UI（Tray / Menu bar）

* 状態表示: 収集中 / 停止中 / オフライン / エラー
* 操作:

  * 今スクショして分類
  * 収集開始/停止
  * 今日のSummaryを表示

### 9.2 設定画面

* 取得間隔
* カテゴリ編集
* 送信/保存ポリシー（画像保存ON/OFF、マスク設定）
* APIキー/認証設定（Copilot SDK利用に必要な情報）

### 9.3 Summary画面

* 日付選択
* タイムライン（カテゴリの帯）
* カテゴリ別集計（合計）
* Markdown出力（コピー/保存）

---

## 10. Copilot Agent 仕様

### 10.1 Agent責務

* スクショ画像を見て、カテゴリ一覧のうち最も適切なものを選ぶ
* 迷う場合は confidence を下げる
* 追加情報（検出アプリ/キーワード）は任意

### 10.2 入出力スキーマ（例）

#### Input

* image: bytes | base64
* categories: { id: string; name: string; description?: string; examples?: string[] }[]
* context?: { localTime: string; timezone: string; userHints?: string[] }

#### Output

* selectedCategoryId: string
* confidence: number
* rationale: string
* detectedApps?: string[]
* detectedKeywords?: string[]

### 10.3 推奨プロンプト（草案）

  * 「あなたは作業スクリーンショットをカテゴリ分類するアシスタント。必ず categories の id から1つ選べ。わからない場合は最も近いものを選び confidence を低くする。」


---

## 11. データモデル（DB）

### 11.1 テーブル: categories

* id (PK)
* name
* description (nullable)
* examples (json, nullable)
* color (nullable)
* createdAt
* updatedAt

### 11.2 テーブル: events

* id (PK)
* capturedAt (indexed)
* categoryId (FK -> categories.id)
* confidence
* status (OK | FAILED | PENDING)
* agentVersion
* screenshotHash (nullable)
* detectedApps (json, nullable)
* detectedKeywords (json, nullable)
* notes (text, nullable)
* createdAt

### 11.3 テーブル（任意）: screenshots

※「画像を保存する」設定時のみ利用

* id (PK)
* eventId (FK)
* storageType (LOCAL | S3 | OTHER)
* mimeType
* blob (encrypted) or storageUrl
* createdAt

---

## 12. 処理フロー

### 12.1 自動収集フロー

0. アプリを起動する
1. Schedulerが起動（interval）
2. Screen Captureでフルスクリーン取得
3. Preprocess（リサイズ/圧縮/マスク）
4. Copilot Agentへ分類依頼
5. 結果をDBに保存
6. アプリが終了したら自動で収集を止める

### 12.2 Summary生成フロー

1. 対象日の events をDBから取得
2. interval推定（設定値 or capturedAt差分の中央値）
3. カテゴリ別に件数×interval で合計時間算出
4. タイムライン整形
5. Markdown/JSONとして出力
6. （任意）Copilot Agentで自然言語の「今日やったこと」を生成

---

## 13. エラーハンドリング
* スクショ取得失敗: OS権限不足 / 画面録画権限
  * 対応: 権限付与ガイドを表示
* Copilot失敗:
  * status=FAILEDで保存し、後から「再分類」できる
* DB書き込み失敗:
  * ローカルキューに積み、復旧後フラッシュ
---

## 14. ログ / 計測

* event処理時間（capture / preprocess / agent / db）
* agent失敗率
* confidence分布（低信頼が多いカテゴリは選択肢を改善する手がかり）

---

## 15. テスト要件

* 単体:
  * カテゴリ管理CRUD
  * Summary集計ロジック（境界: 0件、1件、間隔ブレ）
* 結合:
  * スクショ取得→agent→保存の一連
* E2E:
  * 1日分のダミーイベント投入→Summary生成→期待フォーマット

---

## 16. 受け入れ基準（MVP）

* 任意の間隔でフルスクリーンを取得できる
* Copilot Agentがカテゴリを1つ返し、DBに保存される
* 指定日のSummaryで、カテゴリ別合計とタイムラインが出る
* 画像を保存しない設定で動作する

---

## 17. MVP開発マイルストーン（例）

1. スクショ取得（手動）→分類→保存
2. スクショ取得（自動）＋失敗リトライ
3. カテゴリ定義（固定JSON）＋ローカルDB保存
4. Summary（Markdown出力）
5. 設定UI（カテゴリ編集・間隔変更・保存ポリシー）

---

## 18. 参考 
Copilot SDKのサンプルコード
https://github.com/github/copilot-sdk
https://github.com/github/copilot-sdk/blob/main/cookbook/go/README.md
