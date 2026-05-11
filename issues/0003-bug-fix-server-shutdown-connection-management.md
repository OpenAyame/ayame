# 0003-bug-fix-server-shutdown-connection-management

- Created: 2026-05-11
- Priority: High
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/fix-server-shutdown

## 概要

サーバーシャットダウン時に以下の複数の問題が連鎖し、graceful shutdown が完了しない:

1. 各 WebSocket コネクションの context が `context.Background()` を親としており、サーバーのシャットダウン通知が伝播しない
2. `connection.main()` が `context.Context` を受け取っておらず、select ループ内でシャットダウンを検知できない
3. `connection.wsRecv()` の `ReadMessage()` がブロッキング呼び出しであり、シャットダウンを検知しても即座に抜けられない
4. `server.Start()` が `<-ctx.Done()` 後にキャンセル済みの context で `Shutdown()` を呼んでおり、http.Server の graceful drain が機能していない
5. `StartMatchServer` が先に終了すると `unregisterChannel` への送信が永久ブロックしデッドロックする
6. `ListenAndServe()` のゴルーチンがバッファなしチャネルへの送信でブロックしリークする

## 再現手順

1. サーバー起動後、1 つ以上の WebSocket コネクションを確立する
2. SIGTERM 等でサーバーを停止する
3. `s.Start(ctx)` の `<-ctx.Done()` が発火する
4. defer 内の `s.Shutdown(ctx)` がキャンセル済み context のため即時リターンする（http.Server のドレイン待ちなし）
5. WebSocket コネクションはサーバー停止を検知できず、`ReadMessage()` がクライアント切断を待ち続ける
6. 結果としてプロセスがハングする

## 影響

- シグナルハンドリングによる graceful shutdown が完了せず、プロセスが残留する
- コンテナオーケストレーション（Kubernetes 等）で強制終了（SIGKILL）されるまでプロセスが停止しない
- クライアントに切断通知（Close フレーム、disconnect webhook）が送られない

## 対応方針

以下の問題を 1 つずつ修正する。各修正は独立して検証可能であり、すべて適用することで graceful shutdown が完了する。

### 1. Server に shutdownCh を追加

`Server` 構造体に `shutdownCh` チャネルフィールドを追加する。`context.Context` を構造体に保持しない（Go のコンテキスト設計指針に反するため）。

`NewServer()` でチャネルを初期化し、`Start()` の defer で close する。

```go
type Server struct {
    config          *Config
    signalingLogger *zerolog.Logger
    webhookLogger   *zerolog.Logger
    EchoPrometheus  *echo.Echo
    Metrics         *Metrics
    shutdownCh      chan struct{}  // 追加
    http.Server
}
```

`NewServer()`:

```go
s := &Server{
    // ...
    shutdownCh: make(chan struct{}),
}
```

`Start()`:

```go
func (s *Server) Start(ctx context.Context) error {
    defer close(s.shutdownCh)  // シャットダウン時に close
    // ... 既存の処理
}
```

`signalingHandler` は `Server` のメソッドであるため `s.shutdownCh` にアクセス可能。`signalingHandler` が `ListenAndServe` 開始後にのみ呼ばれることを前提とする。

### 2. signalingHandler でシャットダウンを通知

`signaling_handler.go:70-71`:

修正前:
```go
ctx := context.Background()
ctx, cancel := context.WithCancel(ctx)
```

修正後:
```go
ctx, cancel := context.WithCancel(context.Background())
go func() {
    select {
    case <-s.shutdownCh:
        cancel()
    case <-ctx.Done():
    }
}()
```

元の `context.Background()` + `context.WithCancel` のパターンを維持しつつ、`s.shutdownCh` の close を監視して `cancel` を呼ぶ goroutine を追加する。

### 3. connection.main にシャットダウン監視を追加

`connection.main` のシグネチャを変更し、親 context を受け取る:

修正前:
```go
func (c *connection) main(cancel context.CancelFunc, messageChannel chan []byte) {
```

修正後:
```go
func (c *connection) main(ctx context.Context, cancel context.CancelFunc, messageChannel chan []byte) {
```

select ループに `ctx.Done()` ケースを追加:

```go
loop:
    for {
        select {
        case <-ctx.Done():
            // サーバーシャットダウン
            break loop
        case <-pingTimer.C:
            // ... 既存コード
        // ... 他のケース
        }
    }
```

`signalingHandler` の呼び出しも変更:

```go
go connection.main(ctx, cancel, messageChannel)
```

### 4. ReadMessage ブロッキングからの早期脱出

`main` の defer で `c.wsConn.Close()` を呼び出す。これにより `wsRecv` のブロッキング `ReadMessage()` がエラーを返し、ループを抜ける。

`main` の defer に以下の行を **先頭に** 追加（`cancel()` より先）:

```go
defer func() {
    // ReadMessage のブロッキングを強制解除する
    c.wsConn.Close()
    // ... 既存コード
    timerStop(pongTimeoutTimer)
    timerStop(pingTimer)
    cancel()
    // ...
}()
```

**注意**: `wsConn.Close()` を呼んだ後でも、`sendCloseMessage()` はエラーを返す（既に閉じられた接続のため）が、安全に無視できる。`sendCloseMessage()` が `return` で抜け、その後の defer が実行される流れとなる。

### 5. Start の Shutdown 用 context を修正

`server.go:86-89` で、`<-ctx.Done()` 通過後にキャンセル済み context を `Shutdown()` に渡している問題を修正:

修正前:
```go
defer func() {
    if err := s.Shutdown(ctx); err != nil {
        zlog.Error().Err(err).Send()
    }
}()
```

修正後（独立したタイムアウト付き context を使用）:
```go
defer func() {
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()
    if err := s.Shutdown(shutdownCtx); err != nil {
        zlog.Error().Err(err).Send()
    }
}()
```

### 6. ListenAndServe ゴルーチンリーク修正

`server.go:78-84` でチャネルがバッファなしのため `ListenAndServe` 終了後にブロックする問題を修正:

修正前:
```go
ch := make(chan error)
go func() {
    defer close(ch)
    if err := s.ListenAndServe(); err != nil {
        ch <- err
    }
}()
```

修正後（バッファ付きチャネルを使用）:
```go
ch := make(chan error, 1)
go func() {
    defer close(ch)
    if err := s.ListenAndServe(); err != nil {
        ch <- err
    }
}()
```

### 7. StartMatchServer シャットダウン時の全 room クリーンアップ

`StartMatchServer` が `ctx.Done()` で終了する前に、全 room の全 connection の `forwardChannel` を close する。これにより `main` の `forwardChannel` 受信ケースが `!ok` を検出し、ループを抜けられる。

また、`unregisterChannel` をバッファ付きチャネルに変更することで、`StartMatchServer` 終了後の `unregister()` 呼び出しがデッドロックしないようにする。

`room.go:8`:

修正前:
```go
unregisterChannel = make(chan *unregister)
```

修正後:
```go
unregisterChannel = make(chan *unregister, 1024)
```

```go
// StartMatchServer 内、<-ctx.Done() の前に追加
defer func() {
    for _, r := range m {
        for connID, conn := range r.connections {
            close(conn.forwardChannel)
            delete(r.connections, connID)
        }
    }
}()
```

## テスト戦略

### 統合テスト

- サーバー起動 → WebSocket 接続確立 → `cancel()` 呼び出し → 全 goroutine が 10 秒以内に終了すること
- `errgroup` を使用して複数 goroutine の終了を確認
- `StartMatchServer` が先に終了した場合のデッドロック不在を確認

### 単体テスト

- `connection.main` が ctx.Done() で select ループを抜けること
- `connection.main` が defer で wsConn.Close() を呼ぶこと
- シャットダウン時に wsRecv の ReadMessage がエラーを返しループを抜けること

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
- [FIX] サーバーシャットダウン時に WebSocket コネクションが強制終了されない問題を修正する
  - connection.main に context 監視を追加し、シャットダウン時に即座にループを抜けるようにする
  - server.Start の Shutdown 呼び出しに独立したタイムアウト付き context を使用する
  - StartMatchServer 終了時に全 room の forwardChannel を close しデッドロックを防止する
  - @voluntas
```

## 後方互換

- `Server` 構造体に `shutdownCh` フィールドを追加（非公開フィールド）
- `connection.main` のシグネチャ変更（非公開関数）
- 公開 API への影響なし

## 関連ファイル

- `signaling_handler.go:70-71`（context.Background 使用箇所）
- `server.go:77-97`（Start 関数、Shutdown 呼び出し、ListenAndServe ゴルーチン）
- `server.go:15-25`（Server 構造体定義）
- `connection.go:154-227`（main 関数、select ループ）
- `connection.go:229-259`（wsRecv 関数、ReadMessage ブロッキング）
- `connection.go:138-144`（unregister 関数、unregisterChannel 送信）
- `room.go`（StartMatchServer、unregisterChannel、rooms 管理）
