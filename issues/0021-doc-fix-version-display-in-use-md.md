# 0021-doc-fix-version-display-in-use-md

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Low

## 概要

`docs/USE.md` のバージョン表示例が古い。

## 問題

`docs/USE.md:64` に `version 2025.2.0` と記載されている。現在のバージョンは `VERSION` ファイルの値 `2026.1.2` である。

## 対応方針

`docs/USE.md:64` の `version 2025.2.0` を `version 2026.1.2` に更新する。

## 作業ブランチ

`feature/fix-version-display-in-use-md`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] docs/USE.md のバージョン表示例を 2026.1.2 に更新する
  - @voluntas
```

## 関連ファイル

- `docs/USE.md:64`
- `VERSION:1`（参照用）
