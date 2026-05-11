# 0011-bug-fix-metrics-negative-size-observation

- Created: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/fix-metrics-negative-size

## 概要

メトリクスの `ObserveWebhookResSz` および `ObserveWebhookReqSz` で、HTTP レスポンスの `ContentLength` が `-1`（サイズ不明）の場合に負の値が Prometheus ヒストグラムに記録される。

## 問題

`metrics.go:105-123`:

```go
func (m *Metrics) ObserveWebhookResSz(code, method, host, url string, sz int64) {
    // ...
    m.WebhookResSz.MetricCollector.(*prom.HistogramVec).With(labels).Observe(float64(sz))
}

func (m *Metrics) ObserveWebhookReqSz(code, method, host, url string, sz int64) {
    // ...
    m.WebhookReqSz.MetricCollector.(*prom.HistogramVec).With(labels).Observe(float64(sz))
}
```

`resp.Request.ContentLength` および `resp.ContentLength` は HTTP 仕様上、サイズ不明の場合に `-1` を返す。この `-1` がそのままヒストグラムに記録される。

呼び出し元は `disconnect_webhook.go:48-49` および `webhook.go` の authn webhook。

## 影響

Prometheus ヒストグラムに `-1` が記録され、メトリクスの正確性が損なわれる。アラート閾値が負の値に設定されている場合、意図しないトリガーが発生する。

## 対応方針

両関数で `sz < 0` の場合は観測をスキップする:

```go
func (m *Metrics) ObserveWebhookResSz(code, method, host, url string, sz int64) {
    if sz < 0 {
        return
    }
    labels := prom.Labels{
        "code":   code,
        "method": method,
        "host":   host,
        "url":    url,
    }
    m.WebhookResSz.MetricCollector.(*prom.HistogramVec).With(labels).Observe(float64(sz))
}

func (m *Metrics) ObserveWebhookReqSz(code, method, host, url string, sz int64) {
    if sz < 0 {
        return
    }
    labels := prom.Labels{
        "code":   code,
        "method": method,
        "host":   host,
        "url":    url,
    }
    m.WebhookReqSz.MetricCollector.(*prom.HistogramVec).With(labels).Observe(float64(sz))
}
```

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
- [FIX] メトリクスで ContentLength が -1 の場合に負の値が記録される問題を修正する
  - @voluntas
```

## 後方互換

サイズ不明のリクエスト/レスポンスがメトリクスから除外されるが、これらは元々意味のある値ではないため問題ない。

## 関連ファイル

- `metrics.go:105-123`
