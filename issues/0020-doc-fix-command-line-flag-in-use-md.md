# 0020-doc-fix-command-line-flag-in-use-md

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Medium

## 概要

`docs/USE.md` のコマンドラインフラグの記載が実際の実装と異なる。

## 問題

`docs/USE.md:70`:
```
  -c string
```

`cmd/ayame/main.go:21`:
```go
configFilePath := flag.String("C", "./config.ini", "設定ファイルへのパス")
```

フラグは `-C`（大文字）が正しい。`-c`（小文字）では設定ファイルを指定できない。

## 対応方針

`docs/USE.md:70` の `-c string` を `-C string` に修正する。

## 作業ブランチ

`feature/fix-command-line-flag-in-use-md`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] docs/USE.md のコマンドラインフラグを -C に修正する
  - @voluntas
```

## 後方互換

ドキュメントの修正のみ。

## 関連ファイル

- `docs/USE.md:70`
- `cmd/ayame/main.go:21`（参照用）
