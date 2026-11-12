# 0016-refactor-extract-webhook-common-logic

- Created: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/refactor-extract-webhook-common

## 優先度

Medium。`authn_webhook.go` と `disconnect_webhook.go` の 2 ファイルにわたって HTTP リクエスト送信・メトリクス記録・ステータスコードチェックのボイラープレートが重複している。修正時に両方を更新する必要があり保守コストが高いが、既存動作を変更しないリファクタリングである。0017 の error handling 統一と同時に対応するのが効率的。

## 概要

`authn_webhook.go` と `disconnect_webhook.go` で、HTTP リクエスト送信・メトリクス記録・ステータスコードチェック・ログ出力のボイラープレートが重複している。共通処理を抽出して保守性を向上させる。

## 問題

両関数とも以下の同一パターンを持つ:

1. `c.postRequest()` で HTTP POST
2. `defer resp.Body.Close()`
3. `c.webhookLog()` でリクエストログ出力
4. `url.Parse()` でホスト・パス抽出
5. メトリクス記録（IncWebhookReqCnt、ObserveWebhookReqDur、ObserveWebhookReqSz、ObserveWebhookResSz）
6. `io.ReadAll()` でレスポンスボディ読み取り
7. `httpResponse` 構造体生成
8. ステータスコード 200 チェック
9. `c.webhookLog()` でレスポンスログ出力

## 対応方針

共通の webhook 実行フローを `doWebhook()` ヘルパーに抽出する。レスポンス後の処理（JSON パースなど）は呼び出し側に委ねる。

```go
// doWebhook は webhook URL に POST し、メトリクス記録とステータスコードチェックを共通化する。
// レスポンスボディを []byte で返す。呼び出し側はこれを使って固有の後処理を行う。
func (c *connection) doWebhook(url, logName string, reqBody interface{}) ([]byte, *httpResponse, error) {
    start := time.Now()

    resp, err := c.postRequest(url, reqBody)
    if err != nil {
        return nil, nil, err
    }
    defer resp.Body.Close()

    c.webhookLog(logName+"Req", reqBody)

    u, err := url.Parse(url)
    if err != nil {
        return nil, nil, err
    }
    statusCode := fmt.Sprintf("%d", resp.StatusCode)
    m := c.metrics
    m.IncWebhookReqCnt(statusCode, "POST", u.Host, u.Path)
    m.ObserveWebhookReqDur(statusCode, "POST", u.Host, u.Path, time.Since(start).Seconds())
    m.ObserveWebhookReqSz(statusCode, "POST", u.Host, u.Path, resp.Request.ContentLength)
    m.ObserveWebhookResSz(statusCode, "POST", u.Host, u.Path, resp.ContentLength)

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, nil, err
    }

    httpResp := &httpResponse{
        Status: resp.Status,
        Proto:  resp.Proto,
        Header: resp.Header,
        Body:   string(body),
    }

    c.webhookLog(logName+"Resp", httpResp)

    if resp.StatusCode != 200 {
        return body, httpResp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return body, httpResp, nil
}
```

`authnWebhook()` と `disconnectWebhook()` からはボイラープレートを削除し、`doWebhook()` の呼び出しに置き換える。

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [UPDATE] authn_webhook と disconnect_webhook の共通処理を doWebhook に抽出する
  - @voluntas
```

## 後方互換

非公開ヘルパーの追加と既存メソッドの内部リファクタリングのみ。外部 API への影響はない。

## 関連ファイル

- `authn_webhook.go:29-107`（authnWebhook 関数）
- `disconnect_webhook.go:16-73`（disconnectWebhook 関数）
