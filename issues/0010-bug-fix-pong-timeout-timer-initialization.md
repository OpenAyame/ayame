# 0010-bug-fix-pong-timeout-timer-initialization

- Created: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/fix-pong-timeout-timer

## 概要

`connection.main()` で生成される `pongTimeoutTimer` が接続確立直後からカウントダウンを開始する。クライアントが register メッセージを送信する前にタイムアウトする可能性がある。

## 問題

`connection.go:155-156`:

```go
pongTimeoutTimer := time.NewTimer(time.Duration(c.config.WebSocketPongTimeoutSec) * time.Second)
pingTimer := time.NewTimer(time.Duration(c.config.WebSocketPingIntervalSec) * time.Second)
```

`pongTimeoutTimer` は生成と同時にカウントダウンを開始する。一方、クライアントからの最初の pong は register 処理完了後に ping が送られて初めて返ってくる。したがって pong が一度も受信されない状態でタイムアウトが発火しうる。

`openHandler` のような接続確立時のハンドラが存在しないため、クライアントは register 以前に pong を自発的に送ることはない。

## 影響

クライアントが register メッセージを送信する前に pong タイムアウトが発生し、接続が切断される。

## 対応方針

pong タイムアウトタイマーを無効化した状態で初期化し、最初の pong 受信時にリセットする。

### 方法: `time.NewTimer` の代わりに `time.AfterFunc` からのチャネルでの初期化

シンプルな方法として、初期タイマーは作成せず、pong 受信時に初めて `pongTimeoutTimer.Reset()` を呼ぶ。ただし現在のコードでは `pongTimeoutTimer` は `main` のローカル変数であり、`handleWsMessage` に引数で渡されている。

修正後の `main` のタイマー初期化:

```go
// pongTimeoutTimer は受信時まで無期限に待つ（初期は発火させない）
pongTimeoutTimer := time.NewTimer(0)
if !pongTimeoutTimer.Stop() {
    <-pongTimeoutTimer.C
}
```

`time.NewTimer(0).Stop()` + ドレインにより、タイマーは停止状態で初期化される。最初の register 完了後に ping が送信され、クライアントが pong を返した時点で `handleWsMessage` の pong ケースで `pongTimeoutTimer.Reset()` が呼ばれる。

あるいは、単に初期タイマーの有効時間を十分に長く設定する方法もある:

```go
pongTimeoutTimer := time.NewTimer(time.Duration(c.config.WebSocketPongTimeoutSec) * time.Second)
```

これで十分であり、register がタイムアウト内に完了すれば問題ない。ただし register + authn webhook の処理時間が `WebSocketPongTimeoutSec` を超えるとタイムアウトする。初期タイマーを停止する方法がより堅牢である。

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
- [FIX] pong タイムアウトタイマーが接続直後からカウントダウンする問題を修正する
  - @voluntas
```

## 後方互換

タイムアウトの開始タイミングが変わるが、意図した動作への修正であるため `[FIX]` 扱いとする。

## 関連ファイル

- `connection.go:155-156`（タイマー初期化）
- `connection.go:275-277`（pong 受信時のタイマーリセット）
