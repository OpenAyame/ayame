# 0005-refactor-move-global-channels-to-server

- Created: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/change-move-channels-to-server

## 優先度

Medium。パッケージレベルグローバルチャネルはテスト時の Server インスタンス分離を妨げているが、単一プロセスでの通常運用に支障はない。0003 のシャットダウン修正や 0039 の room テスト追加の前提となるリファクタリングであるため、依存関係を踏まえて優先的に対応する必要がある。

## 目的

`room.go` のパッケージレベルグローバルチャネルを `Server` 構造体のフィールドに移動する。これにより複数 Server インスタンスの起動とテスト時のチャネル注入を可能にする。

## 現状

`room.go:5-11`:

```go
var (
    registerChannel   = make(chan *register)
    unregisterChannel = make(chan *unregister)
    forwardChannel    = make(chan forward, 100)
)
```

これらのチャネルを以下の箇所が直接参照している:

- `room.go`: `StartMatchServer` の select ループ内（26, 51, 70 行目）
- `connection.go`: `register()` メソッド（127 行目）、`unregister()` メソッド（140 行目）、`forward()` メソッド（148 行目）
- `signaling_handler.go`: コネクションごとの `forwardChannel`（58 行目）← これは別物（`connection.forwardChannel`）

注意: グローバルの `forwardChannel` と `connection.forwardChannel` は別のチャネルである。グローバル版は room 間のメッセージ配信に使われ、接続ごとの版は room から各接続へのメッセージ配信に使われる。

## 設計方針

### チャネルの移動

3 つのグローバルチャネルを `Server` 構造体のフィールドに移動する。

名称は役割が明確になるよう、グローバル `forwardChannel` を `dispatchChannel` にリネームする。`connection.forwardChannel` はそのまま維持する。

```go
type Server struct {
    config          *Config
    signalingLogger *zerolog.Logger
    webhookLogger   *zerolog.Logger
    EchoPrometheus  *echo.Echo
    Metrics         *Metrics
    shutdownCh      chan struct{}
    // 以下を追加（グローバルから移動）
    registerChannel   chan *register      // room.go から移動
    unregisterChannel chan *unregister    // room.go から移動
    dispatchChannel   chan forward        // room.go の forwardChannel からリネーム
    http.Server
}
```

### connection へのチャネル参照注入

`connection` 構造体に Server のチャネルへの参照を持たせる。`Server` へのポインタを保持するのではなく、必要なチャネルのみを注入する（依存関係の明確化のため）。

```go
type connection struct {
    // ... 既存フィールド
    registerChannel   chan<- *register    // 追加
    unregisterChannel chan<- *unregister  // 追加
    dispatchChannel   chan<- forward      // 追加
    // forwardChannel は既存のまま（接続ごとのメッセージ受信チャネル）
}
```

### signalingHandler での注入

`signalingHandler` で Server のチャネルを connection に注入する:

```go
connection := connection{
    ID:               getULID(),
    wsConn:           wsConn,
    forwardChannel:   make(chan forward, 100),
    registerChannel:   s.registerChannel,
    unregisterChannel: s.unregisterChannel,
    dispatchChannel:   s.dispatchChannel,
    // ... その他フィールド
}
```

## 解決方法

### 1. Server 構造体にフィールド追加 (`server.go`)

```go
type Server struct {
    // ... 既存フィールド
    registerChannel   chan *register
    unregisterChannel chan *unregister
    dispatchChannel   chan forward
    http.Server
}
```

### 2. NewServer でチャネル初期化 (`server.go`)

```go
func NewServer(config *Config) (*Server, error) {
    // ... 既存コード
    s := &Server{
        config:            config,
        signalingLogger:   signalingLogger,
        webhookLogger:     webhookLogger,
        registerChannel:   make(chan *register),
        unregisterChannel: make(chan *unregister),
        dispatchChannel:   make(chan forward, 100),
        Server: http.Server{
            Addr:              url,
            ReadHeaderTimeout: readHeaderTimeout,
            Handler:           e,
        },
    }
    // ...
}
```

### 3. グローバルチャネル削除 (`room.go`)

`room.go:5-11` の `var` ブロックを削除する。

### 4. StartMatchServer を Server メソッドに (`room.go`)

既に Server メソッドである。グローバルチャネル参照を `s.registerChannel` / `s.unregisterChannel` / `s.dispatchChannel` に変更する。

修正例:

```go
case register := <-s.registerChannel:       // 変更
case unregister := <-s.unregisterChannel:   // 変更
case forward := <-s.dispatchChannel:        // 変更
```

ルーム内接続の `connection.forwardChannel` への close は変更不要（これは接続ごとのチャネルであり Server 移行後も `connection` 構造体に属する）。

### 5. connection 構造体にフィールド追加 (`connection.go`)

```go
type connection struct {
    // ... 既存フィールド
    registerChannel   chan<- *register
    unregisterChannel chan<- *unregister
    dispatchChannel   chan<- forward
    // forwardChannel chan forward は既存のまま
}
```

### 6. connection メソッドのグローバル参照を置換 (`connection.go`)

- `register()` (127 行目): `registerChannel <- ...` → `c.registerChannel <- ...`
- `unregister()` (140 行目): `unregisterChannel <- ...` → `c.unregisterChannel <- ...`
- `forward()` (148 行目): `forwardChannel <- ...` → `c.dispatchChannel <- ...`

### 7. signalingHandler で注入 (`signaling_handler.go`)

connection 生成時にチャネルを注入する（上記「設計方針」のコード参照）。

## 完了条件

- `room.go` のパッケージレベル `var` ブロックが削除されている
- すべてのチャネル参照が `Server` フィールドまたは `connection` フィールド経由になっている
- `go build ./...` が成功する
- 既存テストが通過する
- 複数 `Server` インスタンスの生成が可能である（構造体のフィールドであるため自動的に満たされる）

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [CHANGE] パッケージレベルグローバルチャネルを Server 構造体のフィールドに移動する
  - registerChannel / unregisterChannel / forwardChannel を Server のフィールドに移動
  - forwardChannel を dispatchChannel にリネーム
  - @voluntas
```

変更種別が `[CHANGE]` である理由: `connection` 構造体のフィールド追加により、外部から `connection` を直接生成しているコードに影響がある（ただし公開 API ではないため実質的には限定的）。

## 関連ファイル

- `room.go:5-11`（グローバルチャネル var ブロック削除）
- `room.go:18-84`（StartMatchServer のチャネル参照置換）
- `server.go:15-25`（Server 構造体フィールド追加）
- `server.go:27-73`（NewServer でのチャネル初期化）
- `connection.go:15-50`（connection 構造体フィールド追加）
- `connection.go:125-152`（register / unregister / forward メソッドの参照置換）
- `signaling_handler.go:52-67`（connection 生成時のチャネル注入）
