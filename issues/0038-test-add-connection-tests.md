# 0038-test-add-connection-tests

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: High

## 概要

`connection.go`（447 行）に対応するテストファイルが存在せず、コアロジックのテストカバレッジがゼロ。

## 問題

以下の重要ロジックが未テスト:

| 関数 | 行 | 内容 |
|-----|-----|------|
| `handleWsMessage` | 262-430 | メッセージパース・ルーティングの中核 |
| `main` | 154-227 | コネクション状態マシン |
| `wsRecv` | 229-259 | WebSocket 受信ループ |
| `register` | 125-136 | room 登録 |
| `unregister` | 138-144 | room 解除 |
| `sendAcceptMessage` | 80-95 | accept メッセージ送信 |
| `sendRejectMessage` | 97-107 | reject メッセージ送信 |
| `sendByeMessage` | 109-118 | bye メッセージ送信 |

特に `handleWsMessage` のエラーパス（Invalid JSON、Missing RoomID、Authn webhook 失敗、Room full、Registration incomplete 等）が未テスト。

## 対応方針

`connection_test.go` を新規作成し、`handleWsMessage` を優先的にテーブル駆動テストでカバーする。

authnWebhook はインターフェース抽出または `connection` struct への direct field injection によりモック化する。

最低限カバーすべきテストケース:
- valid JSON → parse success
- invalid JSON → `errInvalidJSON`
- null JSON → `errUnexpectedJSON`
- missing roomID → `errMissingRoomID`
- room full → reject message
- authn webhook allowed → accept message
- authn webhook denied → reject message

## 作業ブランチ

`feature/add-connection-tests`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [ADD] connection.go の単体テストを追加する
  - @voluntas
```

## 関連ファイル

- `connection.go`
- `connection_test.go`（新規作成）
