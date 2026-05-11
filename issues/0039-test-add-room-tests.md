# 0039-test-add-room-tests

- Created: 2026-05-11
- Completed: 2026-05-11
- Priority: High
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/add-room-tests

## 概要

`room.go`（84 行）の `StartMatchServer` が未テスト。マッチングロジックの状態遷移が検証されていない。

**依存**: 0005-refactor-move-global-channels-to-server（チャネルの Server 移動が完了していること）

## 問題

以下のマッチング状態遷移が未テスト:

- 登録 → 1 人目（部屋作成、`one` / `registerResultCreated`）
- 登録 → 2 人目（相手待機中に接続、`two` / `registerResultPaired`）
- 登録 → 満杯（`full` / `registerResultFull`）
- 解除 → room 削除

現在の実装ではグローバルチャネル（`registerChannel`、`unregisterChannel`、`forwardChannel`）を使用しているため、テスト間の分離が困難。テスト実行時のデータレースのリスクもある。

## 依存関係

この issue は 0005（パッケージレベルグローバルチャネルの Server 移動）に依存する。0005 の完了後、チャネルが注入可能になりテストが可能になる。

## 対応方針

0005 の完了後、以下の手順でテストを追加する:

1. 各テストで独立した `Server` インスタンスを作成（チャネルが分離される）
2. `StartMatchServer` を goroutine で起動
3. `registerChannel` に register を送り、結果を検証
4. `unregisterChannel` に unregister を送り、room 削除を検証
5. `go test -race` で並行安全性を検証

## テストケース

- `TestStartMatchServerRegisterFirst`: 1 人目の登録で `one` が返る
- `TestStartMatchServerRegisterSecond`: 2 人目の登録で `two` が返り、両者の forwardChannel が close される
- `TestStartMatchServerRegisterFull`: 3 人目の登録で `full` が返る
- `TestStartMatchServerUnregister`: 解除後に部屋が削除される
- `TestStartMatchServerUnregisterNotRegistered`: 登録されていないコネクションの解除が安全に処理される

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [ADD] room.go の単体テストを追加する
  - @voluntas
```

## 関連ファイル

- `room.go`
- `room_test.go`（新規作成）
