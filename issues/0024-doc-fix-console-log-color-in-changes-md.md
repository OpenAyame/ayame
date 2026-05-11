# 0024-doc-fix-console-log-color-in-changes-md

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Medium

## 概要

`CHANGES.md` で追加されたと記載されている `console_log_color` 設定が、実際には実装されていない。

## 問題

`CHANGES.md:104`:
```
- [ADD] コンソールログのカラー有効化を指定できるようにする
  - 設定ファイルに `console_log_color` を追加
  - デフォルト true
```

しかし、以下の事実と矛盾する:
- `config.go` に `ConsoleLogColor` フィールドが存在しない
- `logger.go:56` で `NoColor: false` がハードコードされている
- `config_example.ini` に `console_log_color` の記載がない

つまり `console_log_color` の設定機能自体がこれまで一度も実装されていない。

## 対応方針

`CHANGES.md` の 2025.2.0 セクションから `console_log_color` の項目を削除し、当該バージョンに実装されていない事実を履歴に反映する。

## 作業ブランチ

`feature/fix-console-log-color-in-changes-md`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] CHANGES.md 2025.2.0 から未実装の console_log_color の記載を削除する
  - @voluntas
```

## 後方互換

ドキュメント修正のみ。必要に応じて別 issue で `console_log_color` の実装を行う。

## 関連ファイル

- `CHANGES.md:104-106`
- `config.go`（参照用、フィールド不在確認）
- `logger.go:56`（参照用、ハードコード確認）
