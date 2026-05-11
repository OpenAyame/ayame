# 0032-refactor-remove-unused-metricskey

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Low

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

## 作業ブランチ

`feature/refactor-remove-unused-metricskey`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [UPDATE] 未使用定数 MetricsKey を削除する
  - @voluntas
```

## 関連ファイル

- `metrics.go:14`
