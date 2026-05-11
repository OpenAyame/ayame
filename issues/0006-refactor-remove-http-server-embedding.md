# 0006-refactor-remove-http-server-embedding

- Created: 2026-05-11
- Priority: Low
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/change-remove-http-server-embedding

## 優先度

Low。`http.Server` の埋め込みにより `Serve()` や `ListenAndServeTLS()` が `Server` の公開 API として露出しているが、実際に誤用された事例はない。`Start()` のみが使用されており、埋め込みによる実害は発生していない。API の整理として時間があるときに対応すればよい。

## 概要

`Server` 構造体が `http.Server` を埋め込んでいるため、`Serve()`、`ListenAndServeTLS()`、`SetKeepAlivesEnabled()` などの `http.Server` のメソッドが `Server` を通じて外部に公開されている。利用者は `Server.Start()` メソッドのみを使用すべきであり、意図しないメソッドの公開は誤用のリスクになる。

## 問題

`server.go:24`:

```go
type Server struct {
    // ...
    http.Server  // 埋め込みにより全 http.Server メソッドが公開される
}
```

## 対応方針

埋め込みをやめ、`httpServer http.Server` という非公開フィールド名で保持する。`server.go` 内での `http.Server` のメソッド呼び出しを `s.httpServer.` に変更する。

### server.go の変更

```go
type Server struct {
    config          *Config
    signalingLogger *zerolog.Logger
    webhookLogger   *zerolog.Logger
    EchoPrometheus  *echo.Echo
    Metrics         *Metrics
    httpServer      http.Server  // 埋め込み → 非公開フィールド
}
```

`Start()` 内:

```go
if err := s.httpServer.ListenAndServe(); err != nil {  // s.ListenAndServe → s.httpServer
    // ...
}
if err := s.httpServer.Shutdown(shutdownCtx); err != nil {  // s.Shutdown → s.httpServer
    // ...
}
```

`NewServer()` 内:

```go
s := &Server{
    // ...
    httpServer: http.Server{  // Server: → httpServer:
        Addr:              url,
        ReadHeaderTimeout: readHeaderTimeout,
        Handler:           e,
    },
}
```

### テストファイルの変更

`ok_handler_test.go:16` および `signaling_handler_test.go:20`:

```go
s := &Server{
    httpServer: http.Server{  // Server: → httpServer:
        Handler: e,
    },
}
```

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [CHANGE] Server 構造体の http.Server 埋め込みを非公開フィールドに変更する
  - @voluntas
```

## 後方互換

`Server` の公開 API から `http.Server` のメソッドが除去される。`Start()` のみを使用しているコードには影響しない。

## 関連ファイル

- `server.go:15-25`（構造体定義）
- `server.go:27-73`（NewServer）
- `server.go:77-97`（Start）
- `ok_handler_test.go:15-18`（テスト）
- `signaling_handler_test.go:19-23`（テスト）
