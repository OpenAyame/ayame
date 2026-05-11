# 0019-doc-fix-go-version-in-use-md

- Created: 2026-05-11
- Completed: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/fix-go-version-in-use-md

## 概要

`docs/USE.md` に記載されている推奨 Go バージョンが `go.mod` と不一致。

## 問題

`docs/USE.md:11` の推奨 Go バージョン:
```
`go 1.24`
```

`go.mod:3` の実際のバージョン:
```
go 1.26.3
```

## 対応方針

`docs/USE.md:11` の `go 1.24` を `go 1.26` に修正する。マイナーバージョンまでで十分であり、パッチバージョンは不要（Go 1.26.x であればどのパッチでも動作するため）。

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] docs/USE.md の推奨 Go バージョンを go 1.26 に修正する
  - @voluntas
```

## 後方互換

ドキュメントの修正のみ。

## 関連ファイル

- `docs/USE.md:11`
- `go.mod:3`（参照用）
