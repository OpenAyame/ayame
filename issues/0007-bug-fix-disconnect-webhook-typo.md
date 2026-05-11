# 0007-bug-fix-disconnect-webhook-typo

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Low

## 概要

`disconnect_webhook.go` のログメッセージにタイプミスがある。

## 問題

`disconnect_webhook.go:53`:

```go
c.errLog().Bytes("body", body).Err(err).Caller().Msg("DiconnectWebhookResponseError")
```

`DiconnectWebhookResponseError` は `DisconnectWebhookResponseError` の誤り（`s` が 1 文字欠落）。

## 影響

ログ検索時に正しいメッセージ `DisconnectWebhookResponseError` でヒットせず、トラブルシューティングの妨げになる。

## 対応方針

`disconnect_webhook.go:53` の `Diconnect` を `Disconnect` に修正する。

修正後:
```go
c.errLog().Bytes("body", body).Err(err).Caller().Msg("DisconnectWebhookResponseError")
```

## 作業ブランチ

`feature/fix-disconnect-webhook-typo`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] disconnect_webhook.go のログメッセージのタイプミスを修正する
  - @voluntas
```

## 後方互換

ログメッセージの文字列変更のみ。API や挙動への影響はない。

## 関連ファイル

- `disconnect_webhook.go:53`
