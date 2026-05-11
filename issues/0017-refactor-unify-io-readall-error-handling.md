# 0017-refactor-unify-io-readall-error-handling

- Created: 2026-05-11
- Completed: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/refactor-unify-io-readall-error

## 概要

`io.ReadAll` エラー発生時のエラー返却方法が `authn_webhook.go` と `disconnect_webhook.go` で一貫していない。

## 問題

| ファイル | 行 | エラー返却 |
|---------|-----|-----------|
| `authn_webhook.go` | 75 | `return nil, err`（生のエラー） |
| `disconnect_webhook.go` | 54 | `return errDisconnectWebhookResponse`（センチネルエラー） |

同一の `io.ReadAll` 失敗という状況に対して、一方は生の `io` エラーを返し、他方は専用のセンチネルエラーを返している。呼び出し側でエラー種別の判定が統一できない。

## 対応方針

両方で専用のセンチネルエラーを使用する。新たに `errWebhookReadBody` を定義する。

### authn_webhook.go:75

修正前:
```go
return nil, err
```

修正後:
```go
return nil, errWebhookReadBody
```

### disconnect_webhook.go:54

修正前:
```go
return errDisconnectWebhookResponse
```

修正後:
```go
return errWebhookReadBody
```

### errors.go に追加

```go
var errWebhookReadBody = errors.New("failed to read webhook response body")
```

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [UPDATE] io.ReadAll エラー時の返却を errWebhookReadBody に統一する
  - @voluntas
```

## 後方互換

エラー種別の変更。呼び出し側が `errWebhookReadBody` で判定できるようになる。

## 関連ファイル

- `authn_webhook.go:75`
- `disconnect_webhook.go:54`
- `errors.go`（errWebhookReadBody 追加）
