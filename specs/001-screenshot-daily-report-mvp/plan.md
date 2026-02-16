# Implementation Plan: スクリーンショット日報作成ツール

**Branch**: `001-screenshot-daily-report` | **Date**: 2026-01-24 | **Spec**: [specs/001-screenshot-daily-report/spec.md](specs/001-screenshot-daily-report/spec.md)
**Input**: Feature specification from `/specs/001-screenshot-daily-report/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

スクリーンショットを定期/手動で取得し、Copilot SDKでカテゴリ分類した結果を保存し、指定日のサマリー（日報）を生成する。技術的にはGo 1.24を中心に、画面キャプチャ・CLI・ローカルSQLite保存・分類ワーカーを組み合わせる。

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.24  
**Primary Dependencies**: Copilot SDK (go), modernc.org/sqlite, macOS screencapture/sips  
**Storage**: ローカルSQLite（ファイル）、画像ファイル（~/.beholder/imgs/）  
**Testing**: Go標準 `go test`（unit/integration）  
**Target Platform**: macOS 13+（主対象）、Windows/Linuxは将来対応  
**Project Type**: single（CLI/デーモン運用）  
**Performance Goals**: スクショ取得 <1秒、分類 5–15秒、日報生成 <2分/1,000件  
**Constraints**: 画像は~/.beholder/imgs/に保存、オフライン時は保留、OS権限が必要、Copilot SDK Attachments機能を利用してトークン制限を回避  
**Scale/Scope**: 個人利用、最大1,000件/日

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Constitutionはテンプレート状態のため、追加制約は未定義。現時点では**No additional gates**として進める。
- Phase 1後の再確認: 追加ゲートなし（変更なし）。

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
cmd/
└── beholder/           # エントリポイント（CLI/デーモン起動）

internal/
├── capture/            # スクショ取得・前処理
├── classify/           # Copilot SDK 連携
├── scheduler/          # 定期実行
├── storage/            # SQLite 永続化
├── summary/            # 日報集計
└── config/             # 設定ファイル読み込み

configs/
└── default.yaml

tests/
├── integration/
└── unit/
```

**Structure Decision**: single構成でGoデスクトップ常駐アプリとして実装する。上記のディレクトリ構成を採用。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
