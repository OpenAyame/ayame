# 0026-doc-fix-upadte-typo-in-changes-md

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Low

## 概要

`CHANGES.md` に `[UPADTE]` というタイプミスがある。

## 問題

`CHANGES.md:171`:
```
- [UPADTE] handler に echo を使用するように変更する
```

`UPADTE` は `UPDATE` のタイプミス（`D` と `A` の順序が逆）。

## 対応方針

`CHANGES.md:171` の `[UPADTE]` を `[UPDATE]` に修正する。

## 作業ブランチ

`feature/fix-upadte-typo-in-changes-md`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] CHANGES.md 2023.1.0 の [UPADTE] タイプミスを [UPDATE] に修正する
  - @voluntas
```

## 関連ファイル

- `CHANGES.md:171`
