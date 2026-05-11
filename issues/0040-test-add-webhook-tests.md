# 0040-test-add-webhook-tests

- Created: 2026-05-11
- Priority: High
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/add-webhook-tests

## 優先度

High。認証 webhook（`authn_webhook.go`）と切断 webhook（`disconnect_webhook.go`）は外部の認証サーバーとの連携ポイントであり、不具合が発生した場合の影響範囲が大きい。特に 0001 の nil デリファレンス panic のような致命的バグがテストなしで見過ごされている。`httptest.NewServer` を用いたモックテストにより、allowed true/false、エラーレスポンス、JSON パースエラー等のケースを自動検証できるようにする必要がある。

## 概要

Webhook 系（`authn_webhook.go`、`disconnect_webhook.go`、`webhook.go`）のテストが存在しない。

## 問題

以下の関数が未テスト:

| ファイル | 行 | 関数 |
|---------|-----|------|
| `authn_webhook.go` | 29-107 | `authnWebhook()` |
| `disconnect_webhook.go` | 16-73 | `disconnectWebhook()` |
| `webhook.go` | 18-43 | `postRequest()` |

**authnWebhook の未カバーケース**:
- URL 未設定時の自動成功（allowed=true）
- HTTP エラー時の `errAuthnWebhook` 返却
- ステータスコード 200 以外の拒否
- JSON パースエラー時の `errAuthnWebhookResponse` 返却
- allowed=true の正常系
- allowed=false の拒否系

**disconnectWebhook の未カバーケース**:
- URL 未設定時の早期リターン
- HTTP エラー時の `errDisconnectWebhook` 返却
- ステータスコード 200 以外のエラー
- 正常系の 200 返却

## 対応方針

`net/http/httptest.NewServer` でモック Webhook サーバーを立て、各レスポンスパターンをテストする。

### postRequest のテスト

```go
func TestPostRequest(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
        assert.Equal(t, "POST", r.Method)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"allowed": true}`))
    }))
    defer ts.Close()

    c := &connection{config: &Config{}}
    resp, err := c.postRequest(ts.URL, map[string]string{"roomId": "test"})
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

### authnWebhook のテスト

```go
func TestAuthnWebhookSuccess(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"allowed": true}`))
    }))
    defer ts.Close()

    c := &connection{
        config:         &Config{AuthnWebhookURL: ts.URL},
        signalingKey:   strPtr("test-key"),
        errLog:         func() *zerolog.Event { return zlog.Err(nil) },
        webhookLogger:  &zlog.Logger,
        webhookLog:     func(n string, v interface{}) {},
        metrics:        &Metrics{},
    }
    // postRequest の実装が必要
    resp, err := c.authnWebhook()
    assert.NoError(t, err)
    assert.True(t, *resp.Allowed)
}
```

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [ADD] webhook 系の単体テストを追加する
  - @voluntas
```

## 関連ファイル

- `authn_webhook.go`
- `disconnect_webhook.go`
- `webhook.go`
- `authn_webhook_test.go`（新規作成）
- `disconnect_webhook_test.go`（新規作成）
