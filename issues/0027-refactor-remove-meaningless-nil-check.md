# 0027-refactor-remove-meaningless-nil-check

- Created: 2026-05-11
- Priority: Low
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/refactor-remove-meaningless-nil-check

## 優先度

Low。`message == nil` チェックは `message := &message{}` で非 nil に初期化されているため常に false であり、到達不能コードである。存在しても実行時影響はなく、削除はコードクリーンアップの一環として行う。他の修正と同時に対応すればよい。

## 概要

`handleWsMessage` 内の `message == nil` チェックが常に `false` であり無意味。

## 問題

`connection.go:263-272`:

```go
message := &message{}
if err := json.Unmarshal(rawMessage, &message); err != nil {
    c.errLog().Err(err).Bytes("rawMessage", rawMessage).Msg("InvalidJSON")
    return errInvalidJSON
}

if message == nil {  // 常に false
    c.errLog().Bytes("rawMessage", rawMessage).Msg("UnexpectedJSON")
    return errUnexpectedJSON
}
```

`message` は `&message{}` で非 nil に初期化されており、`json.Unmarshal` がポインタを nil に変更することはない。この nil チェックは到達不能コードである。

## 対応方針

nil チェックブロック（269-272 行）を削除する。

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] handleWsMessage の到達不能 nil チェックを削除する
  - @voluntas
```

## 関連ファイル

- `connection.go:269-272`
