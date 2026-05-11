# 0009-bug-fix-signaling-log-filters-empty-slice

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Medium

## 概要

`SignalingLogFilters` が空スライス `[]string{}` の場合、デフォルトフィルターが適用されず、すべてのシグナリングログがフィルタリングされて何も出力されない。

## 問題

`config.go:149-151`:

```go
if config.SignalingLogFilters == nil {
    config.SignalingLogFilters = defaultSignalingLogFilters
}
```

`nil` チェックのみで空スライスを考慮していない。ini ファイルで `signaling_log_filters =` と空で設定した場合、空スライス `[]string{}` になる可能性がある。この場合、フィルター一覧が空のため、すべてのシグナリングログが出力されなくなる。

## 再現手順

1. 設定ファイルで `signaling_log_filters =` と空で指定する
2. サーバーを起動する
3. シグナリングログが一切出力されないことを確認する

## 影響

設定ファイルで空のフィルターを指定すると、デバッグに必要なシグナリングログが失われる。

## 対応方針

`nil` チェックに加えて `len` による空チェックを追加する:

修正後:
```go
if len(config.SignalingLogFilters) == 0 {
    config.SignalingLogFilters = defaultSignalingLogFilters
}
```

`len(nil) == 0` は `true` のため、`nil` チェックと空スライスチェックの両方を 1 行でカバーできる。

## 作業ブランチ

`feature/fix-signaling-log-filters-empty`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
- [FIX] SignalingLogFilters が空スライスの場合にデフォルトフィルターが適用されない問題を修正する
  - @voluntas
```

## 後方互換

挙動が変わるため `[FIX]` 扱い。空スライス設定時に意図的に全フィルタリングを期待していた場合は影響があるが、そのようなユースケースは想定されていない。

## 関連ファイル

- `config.go:149-151`
