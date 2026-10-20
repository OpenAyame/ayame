# 0008-bug-fix-webhook-log-client-id-value

- Created: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/fix-webhook-log-client-id

## 優先度

Medium。`webhookLog` が接続 ID（ULID）を `clientId` フィールドに記録している一方、`signalingLog` はクライアント ID を記録している。さらにフィールドキーも `clientID` と `clientId` で揺れている。ログ集計時に別フィールドとして扱われ、接続 ID とクライアント ID が混同される。分析の正確性を損なうため優先的に対応する。

## 概要

`webhookLog` で `clientId` フィールドに接続 ID（`c.ID`、ULID）を記録している。同じフィールド名で他のログはクライアント ID（`c.clientID`）を記録しており、値の意味が一貫していない。

## 問題

| ファイル | 行 | キー | 値 |
|---------|-----|------|-----|
| `webhook.go` | 48 | `clientId` | `c.ID`（接続 ID / ULID） |
| `connection_log.go` | 19 | `clientId` | `c.clientID`（クライアント ID） |
| `connection_log.go` | 29 | `clientId` | `c.clientID`（クライアント ID） |
| `connection_log.go` | 38 | `clientID` | `c.clientID`（クライアント ID） |
| `connection_log.go` | 45 | `clientID` | `c.clientID`（クライアント ID） |

2 つの問題がある:
1. `webhook.go:48` が `c.ID` を記録している（他は `c.clientID`）
2. `connection_log.go:38,45` のキーが `clientID` になっている（他は `clientId`）

## 影響

ログ集計時に `clientId` と `clientID` が別フィールドとして扱われる。また接続 ID とクライアント ID が同じフィールド名で混在し、分析を誤る可能性がある。

## 対応方針

フィールドキーを `clientId` に統一し、値も `c.clientID` に統一する。

接続 ID（ULID）が必要な場合は別フィールド `connectionId` で記録する方針とするが、`webhookLog` には `roomId` と `clientId` で十分なため `c.ID` の記録は不要。

### webhook.go:48

修正前:
```go
Str("clientId", c.ID).
```

修正後:
```go
Str("clientId", c.clientID).
```

### connection_log.go:38,45

修正前:
```go
Str("clientID", c.clientID).
```

修正後:
```go
Str("clientId", c.clientID).
```

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
- [FIX] webhookLog の clientId フィールドに接続 ID ではなくクライアント ID を記録するように修正する
  - connection_log.go の clientID キーを clientId に統一する
  - @voluntas
```

## 後方互換

ログフィールドの値が変更される。既存のログ集計クエリに影響する可能性があるが、バグ修正として `[FIX]` 扱いとする。

## 関連ファイル

- `webhook.go:48`
- `connection_log.go:19,29,38,45`
