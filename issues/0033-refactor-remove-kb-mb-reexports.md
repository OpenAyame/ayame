# 0033-refactor-remove-kb-mb-reexports

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Low

## 概要

`metrics.go` で `prometheus.KB` と `prometheus.MB` が不要に再エクスポートされている。

## 問題

`metrics.go:11-12`:

```go
const (
    KB = prometheus.KB
    MB = prometheus.MB
)
```

使用箇所 `metrics.go:19-20`:
```go
Buckets: prometheus.LinearBuckets(1*KB, 1*KB, 5),
Buckets: prometheus.LinearBuckets(1*MB, 1*MB, 5),
```

## 対応方針

再エクスポートを削除し、使用箇所で直接 `prometheus.KB` / `prometheus.MB` を使用する:

```go
Buckets: prometheus.LinearBuckets(1*prometheus.KB, 1*prometheus.KB, 5),
Buckets: prometheus.LinearBuckets(1*prometheus.MB, 1*prometheus.MB, 5),
```

`const` ブロック（11-13 行）を削除する。

## 作業ブランチ

`feature/refactor-remove-kb-mb-reexports`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [UPDATE] metrics.go の不要な KB/MB 再エクスポートを削除する
  - @voluntas
```

## 関連ファイル

- `metrics.go:11-12`（定数定義削除）
- `metrics.go:19-20`（使用箇所修正）
