# 0037-fix-add-skip-reason-to-test

- Created: 2026-05-11
- Priority: Low
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/fix-add-skip-reason-to-test

## 概要

`signaling_handler_test.go` の `t.Skip("")` が理由なしでスキップされている。

## 問題

`signaling_handler_test.go:13`:

```go
func TestSignalingHandler(t *testing.T) {
    t.Skip("")
}
```

スキップ理由が空文字列のため、なぜこのテストがスキップされているのか、いつ再有効化すべきかが判断できない。

## 対応方針

スキップ理由を明記する:

```go
t.Skip("WebSocket upgrade requires real TCP connection; use integration test")
```

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] signaling_handler_test.go の t.Skip に理由を追加する
  - @voluntas
```

## 関連ファイル

- `signaling_handler_test.go:13`
