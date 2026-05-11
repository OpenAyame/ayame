# 0023-doc-fix-ayame-ini-reference-in-use-md

- Created: 2026-05-11
- Completed: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/fix-ayame-ini-reference-in-use-md

## 概要

`docs/USE.md` で旧設定ファイル名 `ayame.ini` が使用されており、現在のデフォルト `config.ini` と矛盾している。

## 問題

`docs/USE.md:90` で `ayame.ini` に言及している。デフォルト設定ファイル名は 2023.1.0 で `config.ini` に変更されている（`CHANGES.md:157-158`）。

## 対応方針

`docs/USE.md:90` の `ayame.ini` を `config.ini` に修正する。

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] docs/USE.md の ayame.ini を config.ini に修正する
  - @voluntas
```

## 関連ファイル

- `docs/USE.md:90`
- `CHANGES.md:157-158`（参照用）
