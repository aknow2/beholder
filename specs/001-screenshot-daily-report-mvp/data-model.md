# Data Model: スクリーンショット日報作成ツール

## Entity: Category
- **Purpose**: 作業分類の選択肢
- **Fields**:
  - `id` (string, PK)
  - `name` (string, required, 1..50)
  - `description` (string, optional, <= 200)
  - `examples` (string[], optional)
  - `color` (string, optional, hex)
  - `createdAt` (timestamp)
  - `updatedAt` (timestamp)
- **Validation**:
  - `name` は必須、ユニーク

## Entity: Event
- **Purpose**: 1回のスクショ取得・分類結果
- **Fields**:
  - `id` (string, PK)
  - `capturedAt` (timestamp, required)
  - `categoryId` (string, FK -> Category.id)
  - `confidence` (float, 0..1)
  - `status` (enum: `PENDING` | `FAILED` | `OK`)
  - `agentVersion` (string)
  - `screenshotHash` (string, optional)
  - `detectedApps` (string[], optional)
  - `detectedKeywords` (string[], optional)
  - `notes` (string, optional)
  - `createdAt` (timestamp)
- **Validation**:
  - `capturedAt` は必須
  - `confidence` は 0..1
- **State Transitions**:
  - `PENDING` -> `OK`
  - `PENDING` -> `FAILED`
  - `FAILED` -> `PENDING` (再分類)

## Entity: Screenshot (optional)
- **Purpose**: 画像保存が有効な場合のみ保存
- **Fields**:
  - `id` (string, PK)
  - `eventId` (string, FK -> Event.id)
  - `storageType` (enum: `LOCAL` | `S3` | `OTHER`)
  - `mimeType` (string)
  - `blob` (bytes, encrypted) OR `storageUrl` (string)
  - `createdAt` (timestamp)

## Entity: Settings
- **Purpose**: ユーザー設定
- **Fields**:
  - `id` (string, singleton)
  - `captureIntervalMinutes` (int, 1..60)
  - `saveImages` (bool)
  - `maskRegions` (array of rect, optional)
  - `categoryOrder` (string[], optional)
  - `createdAt` (timestamp)
  - `updatedAt` (timestamp)

## Entity: Summary (derived)
- **Purpose**: 日報出力
- **Fields**:
  - `date` (YYYY-MM-DD)
  - `totalsByCategory` (map categoryId -> minutes)
  - `timeline` (list of {timeRange, categoryId})
  - `format` (enum: `MARKDOWN` | `JSON` | `TEXT`)

## Relationships
- Category 1..* Event
- Event 0..1 Screenshot
- Settings 1..* Category (並び順) / 1..* Event (運用ポリシー)
