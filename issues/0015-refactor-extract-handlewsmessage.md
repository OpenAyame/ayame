# 0015-refactor-extract-handlewsmessage

- Created: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/refactor-extract-handlewsmessage

## 優先度

Medium。`handleWsMessage` が 169 行の monolithic 関数になっており、特に `case "register":` ブロックが約 100 行を占める。テスト容易性と可読性を著しく損なっているが、既存動作を変更しない内部リファクタリングである。0015 の分割は後続のテスト追加（0038）の前提となるため、テスト戦略全体の中で優先的に対応する。

## 概要

`handleWsMessage` 関数（169 行）が肥大化しており、JSON パース・認証・登録・転送・ログが 1 関数に混在している。これを責務ごとに分割し、テスト可能にする。

## 問題

`connection.go:262-430`:

```go
func (c *connection) handleWsMessage(rawMessage []byte, pongTimeoutTimer *time.Timer) error {
```

この関数は以下の責務をすべて持っている:
- JSON パースと型判定
- pong 処理
- register 処理（約 100 行）: メッセージパース、バリデーション、認証 webhook、accept/reject 送信
- シグナリングメッセージ（pong 以外）の転送
- connected メッセージの処理

単一責任原則に違反しており、register のテストだけを書くことができない。

## 対応方針

以下の関数に分割する:

### `handleRegister()`

register メッセージのパース・バリデーション・フィールド設定・認証 webhook 呼び出し・accept/reject 送信を担当。

```go
func (c *connection) handleRegister(rawMessage []byte) error
```

戻り値で register が完了したかどうかを示し、呼び出し側でシグナリングログ出力を行う。

### `handleSignalingMessage()`

pong 以外の type 付きシグナリングメッセージ（offer、answer、candidate 等）の転送を担当。

```go
func (c *connection) handleSignalingMessage(rawMessage []byte) error
```

### `handleWsMessage()` （ディスパッチャ）

分割後はディスパッチャに徹する:

```go
func (c *connection) handleWsMessage(rawMessage []byte, pongTimeoutTimer *time.Timer) error {
    message := &message{}
    if err := json.Unmarshal(rawMessage, &message); err != nil {
        c.errLog().Err(err).Bytes("rawMessage", rawMessage).Msg("InvalidJSON")
        return errInvalidJSON
    }
    if message == nil {
        c.errLog().Bytes("rawMessage", rawMessage).Msg("UnexpectedJSON")
        return errUnexpectedJSON
    }

    switch message.Type {
    case "pong":
        timerStop(pongTimeoutTimer)
        pongTimeoutTimer.Reset(time.Duration(c.config.WebSocketPongTimeoutSec) * time.Second)
        return nil
    case "register":
        return c.handleRegister(rawMessage)
    default:
        return c.handleSignalingMessage(rawMessage)
    }
}
```

## テスト戦略

抽出した各関数の単体テストを追加する:

- `TestHandleRegister`: register メッセージのパースとバリデーション
- `TestHandleSignalingMessage`: シグナリングメッセージの転送

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [UPDATE] handleWsMessage を handleRegister / handleSignalingMessage に分割する
  - @voluntas
```

## 後方互換

既存の挙動は変更されない。非公開メソッドの分割リファクタリングであるため、外部 API への影響はない。

## 関連ファイル

- `connection.go:262-430`（handleWsMessage 関数）
