# 0022-doc-fix-signalingkey-case-in-use-md

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Medium

## 概要

`docs/USE.md` で JSON キー `signalingKey` が `"signalingkey"`（全小文字）と誤って記載されている。

## 問題

`docs/USE.md:82,105` に `"signalingkey"` と記載されているが、実際の JSON タグは `"signalingKey"`（camelCase）である。

| 参照 | 値 |
|-----|-----|
| `docs/USE.md:82` | `"signalingkey"` (誤) |
| `ws_messages.go` | `json:"signalingKey"` (正) |
| `authn_webhook.go:15` | `json:"signalingKey"` (正) |

## 対応方針

`docs/USE.md:82` および `docs/USE.md:105` の `"signalingkey"` を `"signalingKey"` に修正する。

## 作業ブランチ

`feature/fix-signalingkey-case-in-use-md`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] docs/USE.md の signalingkey を signalingKey に修正する
  - @voluntas
```

## 関連ファイル

- `docs/USE.md:82,105`
- `ws_messages.go:13`（参照用）
