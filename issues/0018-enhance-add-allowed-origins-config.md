# 0018-enhance-add-allowed-origins-config

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: High

## 概要

WebSocket アップグレーダーの `CheckOrigin` が常に `true` を返しており、任意の Origin からの WebSocket 接続を許可している。本番環境ではクロスオリジン WebSocket 接続（Cross-Site WebSocket Hijacking）のリスクがある。

## 問題

`signaling_handler.go:20-28`:

```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024 * 4,
    WriteBufferSize: 1024 * 4,
    CheckOrigin: func(r *http.Request) bool {
        return true  // 全オリジンを許可
    },
}
```

`upgrader` はパッケージレベル変数であるため、設定値に応じて動的に `CheckOrigin` を切り替えられない。

## 対応方針

`allowed_origins` 設定を追加し、設定に基づいて Origin を検証する。

### config.go

デフォルトは空スライスとし、空の場合は全オリジンを拒否する（セキュアデフォルト）:

```go
type Config struct {
    // ... 既存フィールド
    AllowedOrigins []string `ini:"allowed_origins"`  // 追加
}
```

### signaling_handler.go

`upgrader` をパッケージレベル変数から、`signalingHandler` 内で毎回生成する方式に変更する。または `CheckOrigin` のみ動的に設定する。

```go
func (s *Server) signalingHandler(c echo.Context) error {
    // ...
    upgrader := websocket.Upgrader{
        ReadBufferSize:  1024 * 4,
        WriteBufferSize: 1024 * 4,
        CheckOrigin: func(r *http.Request) bool {
            origin := r.Header.Get("Origin")
            if origin == "" {
                return true  // 同一オリジンからのリクエスト（Origin なし）
            }
            for _, allowed := range s.config.AllowedOrigins {
                if allowed == origin {
                    return true
                }
            }
            return false
        },
    }
    // ...
}
```

### 設定例

```ini
allowed_origins = https://example.com, https://app.example.com
```

## 作業ブランチ

`feature/add-allowed-origins`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
- [ADD] WebSocket 接続時の Origin チェック用に allowed_origins 設定を追加する
  - デフォルトは空（全オリジン拒否）
  - @voluntas
```

## 後方互換

**破壊的変更**: 従来は全オリジンが許可されていたが、修正後は `allowed_origins` を明示的に設定しないと WebSocket 接続が拒否される。CHANGES.md の種別は `[CHANGE]` 相当だが、セキュリティ修正のため `[FIX]` 扱いとする判断もある。

## 関連ファイル

- `signaling_handler.go:20-28`（upgrader 定義）
- `config.go`（Config 構造体と SetDefaultsConfig）
