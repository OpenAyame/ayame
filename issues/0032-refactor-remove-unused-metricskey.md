# 0032-refactor-remove-unused-metricskey

- Created: 2026-05-11
- Priority: Low
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/refactor-remove-unused-metricskey

## 優先度

Low。`MetricsKey = "webhook_metrics"` はコードベース内で 1 度も参照されていない未使用定数である。コンパイルや実行に影響しないが、不要なコードの存在は混乱を招く。1 行削除で完了するため、他の修正と同時に対応すればよい。

## 概要

未使用定数 `MetricsKey` が定義されている。

## 問題

`metrics.go:14`:

```go
MetricsKey = "webhook_metrics"
```

コードベース内で `MetricsKey` は一切参照されていない。未使用コードの存在は混乱を招く。

## 対応方針

`metrics.go:14` の `MetricsKey` 定義を削除する。

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [UPDATE] 未使用定数 MetricsKey を削除する
  - @voluntas
```

## 関連ファイル

- `metrics.go:14`
