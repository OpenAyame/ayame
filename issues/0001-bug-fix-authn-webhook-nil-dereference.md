# 0001-bug-fix-authn-webhook-nil-dereference

- Created: 2026-05-11
- Priority: High
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/fix-authn-webhook-nil-dereference

## 概要

`authnWebhook()` 関数内で、認証 webhook レスポンスの `Allowed` フィールド（`*bool`）が nil の場合に panic する。`Reason` は nil チェックされているが、`Allowed` はチェックされていない。

## 再現手順

1. `config.ini` に `authn_webhook_url = http://127.0.0.1:9999/auth` を設定する
2. 認証 webhook サーバーが `allowed` フィールドを含まない JSON レスポンス（`{}` または `{"reason":"error"}`）を返すようにする
3. クライアントが WebSocket 接続し `{"type":"register","roomId":"test","clientId":"test"}` を送信する
4. `authn_webhook.go:101` または `authn_webhook.go:103` で `*authnWebhookResponse.Allowed` が nil デリファレンスとなり panic する

## 問題

`authn_webhook.go:100-104` で、`Reason` は nil チェックされているが `Allowed` はチェックされていない。

```go
// authn_webhook.go:100-104
if authnWebhookResponse.Reason == nil {
    m.IncAuthnWebhookCnt(statusCode, "POST", u.Host, u.Path, *authnWebhookResponse.Allowed, "") // panic
} else {
    m.IncAuthnWebhookCnt(statusCode, "POST", u.Host, u.Path, *authnWebhookResponse.Allowed, *authnWebhookResponse.Reason) // panic
}
```

`Allowed` は `*bool` 型で、JSON に `allowed` フィールドが含まれていなければ nil になる。

### コードパス

| パス | Allowed の状態 | 現在の挙動 |
|------|---------------|-----------|
| JSON に `allowed` あり | non-nil | 安全 |
| JSON に `allowed` なし | **nil** | **panic** |

## 影響

- 認証 webhook のレスポンスに `allowed` が含まれていない場合、プロセス全体が panic する
- panic 発生時、`IncWebhookReqCnt` 等の先行メトリクスは既に記録されているが `IncAuthnWebhookCnt` は記録されず、メトリクスの部分欠損が発生する
- 後方互換: 修正により panic は解消される。`Allowed` がない場合に reject が返るため、誤った認証レスポンスを返す webhook サーバーとの組み合わせで挙動が変わる可能性がある

## 設計判断

**採用案**: `authnWebhook()` 内で `Allowed == nil` をチェックし、nil なら `return nil, errAuthnWebhookResponse` とする。

- `allowed` は認証 webhook レスポンスの必須フィールドであり、欠損は不正なレスポンスとみなすのが妥当
- メトリクス（`IncAuthnWebhookCnt`）は記録しない。`allowed bool` がシグネチャのため nil を表現できず、不正確な値で記録するより未記録の方が誤解がない
- この設計は `AuthnWebhookURL==""` 時の早期リターンパス（`authn_webhook.go:31`、`IncAuthnWebhookCnt` 未呼び出し）と挙動が整合する
- `connection.go:334-341` の既存 nil チェックは `authnWebhook()` が nil `Allowed` を返さなくなるため到達不能になる。削除する

## 対応方針

1. **`authn_webhook.go` の `if authnWebhookResponse.Reason == nil`（100 行目）の直前に nil チェックを追加**:

   ```go
   if authnWebhookResponse.Allowed == nil {
       c.errLog().Caller().Msg("AuthnWebhookAllowedMissing")
       return nil, errAuthnWebhookResponse
   }
   ```

   注: `AuthnWebhookAllowedMissing` は新規ログメッセージ。

2. **`connection.go:333-341` を削除する**（コメント行 `// 認証サーバの戻り値がおかしい場合は全部 Error にする` を含む丸ごと 9 行）。`authnWebhook()` が nil `Allowed` を返さないため到達不能。

3. **CHANGES.md の `## develop` セクションに以下を追記する**:

   ```
   - [FIX] 認証 webhook レスポンスの allowed フィールドが nil の場合に panic するのを修正する
   ```

### 修正後バリエーション

| シナリオ | 修正前 | 修正後 |
|---------|--------|--------|
| `allowed` なし | panic | `errAuthnWebhookResponse`、ログ `AuthnWebhookAllowedMissing`、呼び出し元で reject（`"InternalServerError"`） |

`allowed` が存在する正常系の挙動は変更なし。

修正後のログ出力フロー: `authn_webhook.go` で `AuthnWebhookAllowedMissing` → `connection.go:325` で `AuthnWebhookError`（共通エラーハンドリング）。2 段階のログ出力になる。

## テスト戦略

### `authn_webhook_test.go` を新規作成

テスト用 `connection` 構造体の構築:

- `Config` に `AuthnWebhookURL` と `WebhookRequestTimeoutSec` を設定
- `Metrics` は `NewMetrics()` で生成する（`MetricCollector` が nil のため、本テストではメトリクス呼び出しのアサーションは行わない）
- `postRequest` に `http.Client{Timeout: ...}` が使われるが、モックサーバーへのリクエストのためタイムアウトは問題にならない
- `errLog()` / `webhookLog()` はグローバルロガー（zerolog）に依存するが、テストではログ出力を許容する

モック webhook サーバー: `httptest.NewServer` で立ち上げる。

| テストケース | レスポンス JSON | HTTP ステータス | 期待結果 |
|------------|----------------|----------------|---------|
| allowed なし | `{}` | 200 | `errAuthnWebhookResponse`、panic しない |
| allowed なし reason のみ | `{"reason":"error"}` | 200 | `errAuthnWebhookResponse`、panic しない |
| allowed=true | `{"allowed":true}` | 200 | エラーなし、`Allowed=true` |

allowed なしの 2 ケースは、Reason の有無にかかわらず `Allowed == nil` の早期リターンに到達することを確認する。

### 補足: メトリクス検証について

`IncAuthnWebhookCnt` は `prometheus.CounterVec` に依存し、単体テスト環境では `MetricCollector` が nil のままになる。メトリクス記録の正確性は Prometheus 統合テストか手動確認に委ね、本テストでは panic 防止とエラー返却の検証に集中する。

## 関連ファイル

- `authn_webhook.go:100-104`（修正箇所）
- `connection.go:333-341`（削除箇所）
- `errors.go:18`（`errAuthnWebhookResponse`）
- `CHANGES.md`（`[FIX]` エントリ追加）
