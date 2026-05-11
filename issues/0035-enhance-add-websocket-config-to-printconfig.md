# 0035-enhance-add-websocket-config-to-printconfig

- Created: 2026-05-11
- Priority: Low
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/add-websocket-config-to-printconfig

## 概要

`PrintConfig()` で WebSocket 設定（`WebSocketReadTimeoutSec`、`WebSocketPongTimeoutSec`、`WebSocketPingIntervalSec`）がログ出力されない。

## 問題

`config.go:191-222` の `PrintConfig()` は各種設定をログ出力するが、以下の WebSocket 設定が含まれていない:

- `WebSocketReadTimeoutSec`
- `WebSocketPongTimeoutSec`
- `WebSocketPingIntervalSec`

## 対応方針

`PrintConfig()` に以下のログ出力を追加する:

```go
zlog.Info().Int32("websocket_read_timeout_sec", c.WebSocketReadTimeoutSec).Msg("AyameConf")
zlog.Info().Int32("websocket_pong_timeout_sec", c.WebSocketPongTimeoutSec).Msg("AyameConf")
zlog.Info().Int32("websocket_ping_interval_sec", c.WebSocketPingIntervalSec).Msg("AyameConf")
```

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [UPDATE] PrintConfig に WebSocket 設定のログ出力を追加する
  - @voluntas
```

## 関連ファイル

- `config.go:191-222`
