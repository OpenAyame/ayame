---
name: ayame
description: WebRTC Signaling Server Ayame (OpenAyame/ayame) を利用するためのリファレンス。Ayame の導入・起動・設定ファイル (config.ini) の全項目、シグナリングプロトコル (register / accept / reject / offer / answer / candidate / bye / ping / pong / message / connected) とそのシーケンス、standalone モード、認証ウェブフック・切断ウェブフックの仕様、ログ (ayame.jsonl / signaling.jsonl / webhook.jsonl)、Prometheus メトリクス、リバースプロキシ配下での運用に関する質問で使用する。
---

# WebRTC Signaling Server Ayame

- **対象バージョン**: 2026.2.0
- **リポジトリ**: https://github.com/OpenAyame/ayame
- **リリースバイナリ**: https://github.com/OpenAyame/ayame/releases
- **仕様書**: https://github.com/OpenAyame/ayame-spec
- **Web SDK**: https://github.com/OpenAyame/ayame-web-sdk (npm: `@open-ayame/ayame-web-sdk`)
- **Web SDK サンプル**: https://github.com/OpenAyame/ayame-web-sdk-examples
- **ホスティング版**: [Ayame Labo](https://ayame-labo.shiguredo.app/) (時雨堂が無料提供。STUN/TURN とシグナリングキー認証込み)
- **言語**: Go (バージョンは `go.mod` に従う。2026.2.0 は Go 1.27)
- **ライセンス**: Apache License 2.0

WebRTC の P2P 接続を確立するためだけのシグナリングサーバ。クライアント同士の offer / answer / ICE candidate を WebSocket 経由で *そのまま* 相手に中継する。
メディアは一切扱わない (SFU / MCU ではない)。

質問・相談は時雨堂の Discord (https://discord.gg/shiguredo) のみ。Discord で議論していない PR / issue には対応しない。

## 前提と制限

- P2P 専用。1 ルーム最大 2 名。3 人目の `register` は `reject` (`reason: "full"`) になる
- エンドポイントは WebSocket の `/signaling` と、ヘルスチェック用の `GET /.ok` (200 / 空ボディ) の 2 つだけ
- TLS 終端は持たない。`wss://` が必要なら nginx / Caddy / ngrok 等のリバースプロキシを前段に置く
- STUN / TURN は持たない。必要なら認証ウェブフックのレスポンス `iceServers` でクライアントに払い出す
- WebSocket の Origin チェックはしない (全許可)。接続制御は認証ウェブフックで行う
- WebSocket 1 メッセージの上限は 1 MB
- サンプルコードは同梱しない (Web SDK とそのサンプルが一次資料)

## 導入

### リリースバイナリ

linux / darwin の amd64 / arm64 向けに gzip 圧縮した単一バイナリが配布される。Windows 向けはない。

```bash
VERSION=2026.2.0   # リリースページのタグ名
curl -LO https://github.com/OpenAyame/ayame/releases/download/${VERSION}/ayame_linux_amd64.gz
gunzip ayame_linux_amd64.gz
chmod +x ayame_linux_amd64
```

### ソースからビルド

```bash
git clone https://github.com/OpenAyame/ayame.git
cd ayame
make            # bin/ayame を生成
make init       # config_example.ini を config.ini にコピー (既存なら上書きしない)
```

Makefile の主なターゲット:

| ターゲット | 内容 |
| --- | --- |
| `ayame` (デフォルト) | `bin/ayame` をビルド |
| `linux` / `darwin` | amd64 向けクロスビルド (`bin/ayame-linux` / `bin/ayame-darwin`) |
| `init` | `config_example.ini` → `config.ini` |
| `test` | `go test -race -v` |
| `lint` / `fmt` | golangci-lint |

### 起動

```bash
./bin/ayame -C config.ini
```

| フラグ | 説明 |
| --- | --- |
| `-C <path>` | 設定ファイル (ini) のパス (default: `./config.ini`)。**大文字の C**。`-c` は受け付けない |
| `-V` | バージョンを表示して終了 |

起動すると 3 つの goroutine が動く: シグナリング HTTP サーバ (`listen_ipv4_address:listen_port_number`)、Prometheus 用 HTTP サーバ (`listen_prometheus_ipv4_address:listen_prometheus_port_number`)、ルーム管理 (マッチング) ループ。

動作確認:

```bash
curl -i http://127.0.0.1:3000/.ok        # HTTP/1.1 200 OK
curl http://127.0.0.1:4000/metrics       # Prometheus メトリクス
```

注意点:

- **デフォルトではコンソールに何も出ない**。ログはカレントディレクトリの `ayame.jsonl` / `signaling.jsonl` / `webhook.jsonl` に書かれる。標準出力に出したいなら `log_stdout = true`、開発中なら `debug = true` + `debug_console_log = true`
- `log_dir` は**存在するディレクトリ**でなければ起動に失敗する (自動作成しない)
- 設定ファイルに**未知のキーがあると起動に失敗する** (厳密マッピング)。キー名の大文字小文字は区別しない
- `authn_webhook_url` / `disconnect_webhook_url` は URL として妥当でないと起動に失敗する

## 設定ファイル (config.ini)

ini 形式。セクションなし、`key = value`。全項目省略可能で、省略時はデフォルト値になる。`config_example.ini` に全項目がコメントアウト付きで書かれている。

### ログ

| キー | デフォルト | 説明 |
| --- | --- | --- |
| `log_dir` | `.` | ログファイルの出力先ディレクトリ (存在必須) |
| `log_name` | `ayame.jsonl` | サーバログのファイル名 |
| `log_stdout` | `false` | `true` で全ログを標準出力に出す (ファイルには書かない) |
| `log_rotate_max_size` | `200` | ローテーションするサイズ (MB)。全ログ共通 |
| `log_rotate_max_backups` | `7` | 保持する世代数 |
| `log_rotate_max_age` | `30` | 保持日数 |
| `log_rotate_compress` | `false` | ローテーション後のファイルを gzip 圧縮する |
| `log_message_key_name` | `message` | ログの JSON でメッセージを入れるキー名 |
| `log_timestamp_key_name` | `time` | ログの JSON でタイムスタンプを入れるキー名 |
| `signaling_log_name` | `signaling.jsonl` | シグナリングログのファイル名 |
| `signaling_log_filters` | 全て | シグナリングログに出す `type` をカンマ区切りで指定。指定可能: `register,offer,answer,candidate,connected,message` |
| `webhook_log_name` | `webhook.jsonl` | ウェブフックログのファイル名 |

### デバッグ

| キー | デフォルト | 説明 |
| --- | --- | --- |
| `debug` | `false` | ログレベルを debug にする。シグナリングログに `rawMessage` (受信 JSON そのもの) が含まれるようになる |
| `debug_console_log` | `false` | `debug = true` のときのみ有効。ログをファイルではなく標準出力に色付きで出す (呼び出し元ファイル名付き) |
| `debug_console_log_json` | `false` | `debug_console_log` の出力を JSON にする |

### 機能

| キー | デフォルト | 説明 |
| --- | --- | --- |
| `type_message` | `false` | `true` で `{"type": "message", ...}` の中継を許可する。`false` のときに `message` を送ると切断される |

### 待ち受け

| キー | デフォルト | 説明 |
| --- | --- | --- |
| `listen_ipv4_address` | `0.0.0.0` | シグナリングの待ち受けアドレス |
| `listen_port_number` | `3000` | シグナリングの待ち受けポート |
| `listen_prometheus_ipv4_address` | `0.0.0.0` | Prometheus `/metrics` の待ち受けアドレス |
| `listen_prometheus_port_number` | `4000` | Prometheus `/metrics` の待ち受けポート |

### WebSocket

| キー | デフォルト | 説明 |
| --- | --- | --- |
| `websocket_read_timeout_sec` | `90` | この秒数クライアントから何も受信しないと切断する |
| `websocket_ping_interval_sec` | `5` | `{"type": "ping"}` を送る間隔 |
| `websocket_pong_timeout_sec` | `60` | この秒数 `{"type": "pong"}` が返ってこないと切断する |

### ウェブフック

| キー | デフォルト | 説明 |
| --- | --- | --- |
| `authn_webhook_url` | (なし) | 認証ウェブフックの URL。未設定なら全員許可 |
| `disconnect_webhook_url` | (なし) | 切断ウェブフックの URL。未設定なら何もしない |
| `webhook_request_timeout_sec` | `5` | ウェブフック HTTP リクエストのタイムアウト |
| `copy_websocket_header_names` | (なし) | WebSocket ハンドシェイク (HTTP Upgrade) リクエストのヘッダーのうち、ここに列挙した名前のものをウェブフックリクエストの HTTP ヘッダーにコピーする。カンマ区切り、大文字小文字を区別しない。例: `X-Forwarded-For, X-Real-IP` |

### 最小構成例

```ini
log_dir = /var/log/ayame
listen_ipv4_address = 127.0.0.1
listen_port_number = 3000
listen_prometheus_ipv4_address = 127.0.0.1
listen_prometheus_port_number = 4000

copy_websocket_header_names = X-Forwarded-For
authn_webhook_url = http://127.0.0.1:3001/ayame/webhook/authn
disconnect_webhook_url = http://127.0.0.1:3001/ayame/webhook/disconnect
```

## シグナリングプロトコル

- エンドポイント: `ws://<host>:<listen_port_number>/signaling`
- WebSocket のテキストフレームで JSON をやり取りする。全メッセージは `type` (string) を持つ
- `type` が未知、JSON として壊れている、`register` 前に `offer` 等を送る、などのプロトコル違反は **`reject` を返さずに WebSocket を閉じる** (Close frame 1000)
- Ayame が自分の都合で接続を終わらせるときは常に Close frame (1000) を送ってから閉じる

### シーケンス

```text
client1                    Ayame                     client2          authn server
  |-- WS 接続 ------------->|                           |                   |
  |-- register ------------>|-- POST authn ------------------------------->|
  |                         |<-- {"allowed": true, ...} -------------------|
  |<-- accept --------------|   (isExistClient: false → 相手を待つ)        |
  |                         |<-- WS 接続 ---------------|                   |
  |                         |<-- register --------------|                   |
  |                         |-- POST authn ------------------------------->|
  |                         |<-- {"allowed": true, ...} -------------------|
  |                         |-- accept ---------------->| (isExistClient: true → offer を作る)
  |<-- offer ---------------|<-- offer -----------------|
  |-- answer -------------->|-- answer ---------------->|
  |<-> candidate <--------->|<-> candidate <----------->|   (Trickle ICE、以後も随時)
  |   ... P2P 確立。以後 WS は ping/pong と bye の通知にだけ使われる ...  |
  |                         |<-- WS 切断 ---------------|
  |<-- bye -----------------|                           |-- (切断 webhook)  |
  |<-- Close (1000) --------|                           |                   |
  |-- (切断 webhook)        |                           |                   |
```

要点:

- 先に `register` した側が `isExistClient: false` を受け取り、**後から来た側 (`isExistClient: true`) が offer を作る**
- `offer` / `answer` / `candidate` の中身を Ayame は検証も加工もしない。受け取った JSON のバイト列をそのまま相手に書き出す
- どちらかが切断するとルームは削除され、残った側に `bye` が届いてから WebSocket が閉じられる。再接続は両者とも `register` からやり直し
- ping/pong は WebSocket の制御フレームではなく **JSON メッセージ**。クライアントは `{"type": "ping"}` を受け取ったら `{"type": "pong"}` を返す必要がある。pong タイマーは WebSocket 接続直後から動くので、pong を返さないクライアントは (register の有無に関係なく) `websocket_pong_timeout_sec` 秒で切断される

### クライアント → Ayame

#### `register`

接続後、最初に 1 回だけ送る。2 回目を送ると切断される。

| プロパティ | 型 | 必須 | 説明 |
| --- | --- | --- | --- |
| `type` | string | 必須 | `"register"` |
| `roomId` | string | 必須 | ルーム ID。欠けていると切断 |
| `clientId` | string | 任意 | 省略時は Ayame が払い出した `connectionId` と同じ値になる |
| `signalingKey` | string | 任意 | 認証ウェブフックにそのまま渡す。後方互換で `key` も受け付け、両方あれば `signalingKey` を優先 (Web SDK は `key` で送る) |
| `authnMetadata` | 任意の JSON | 任意 | 認証ウェブフックにそのまま渡す |
| `ayameClient` / `environment` / `libwebrtc` | string | 任意 | クライアント情報。認証ウェブフックにそのまま渡す |
| `standalone` | boolean | 任意 | `true` で standalone モード (後述) |

```json
{"type": "register", "roomId": "room1", "clientId": "alice", "signalingKey": "xxx", "authnMetadata": {"user": "alice"}}
```

#### `offer` / `answer` / `candidate`

`accept` 受信後にのみ送れる。中身はそのまま相手に転送される。慣例 (仕様書 / Web SDK) の形式:

```json
{"type": "offer", "sdp": "v=0\r\no=- ..."}
{"type": "answer", "sdp": "v=0\r\no=- ..."}
{"type": "candidate", "ice": {"candidate": "candidate:...", "sdpMid": "0", "sdpMLineIndex": 0}}
```

#### `pong`

`ping` への応答。`{"type": "pong"}`。

#### `message`

`type_message = true` のときだけ許可される任意メッセージ。`accept` 後にのみ送れ、そのまま相手に転送される。シグナリング経路を使った軽量なアプリメッセージ用。

```json
{"type": "message", "data": {"anything": "you want"}}
```

#### `connected`

standalone モード用。PeerConnection が接続済みになったことを Ayame に伝える。非 standalone のクライアントが送った場合はログに残るだけで無視される。

### Ayame → クライアント

#### `accept`

`register` が受理された。

| プロパティ | 型 | 説明 |
| --- | --- | --- |
| `type` | string | `"accept"` |
| `connectionId` | string | Ayame が払い出した接続 ID (ULID)。ログとウェブフックにも同じ値が出る |
| `isExistClient` | boolean | `true` なら相手が既にいるので **自分が offer を作る**。`false` なら相手の offer を待つ |
| `isExistUser` | boolean | `isExistClient` と同じ値 (後方互換のため残っている) |
| `authzMetadata` | 任意の JSON | 認証ウェブフックが返した場合のみ付く |
| `iceServers` | array | 認証ウェブフックが返した場合のみ付く。`RTCPeerConnection` の `iceServers` にそのまま渡せる形式 |

```json
{"type": "accept", "connectionId": "01J...", "isExistClient": true, "isExistUser": true,
 "iceServers": [{"urls": ["turn:turn.example.com:3478"], "username": "u", "credential": "p"}],
 "authzMetadata": {"role": "owner"}}
```

#### `reject`

`register` を拒否した。この後 WebSocket は閉じられる。クライアントは PeerConnection を破棄する。

| `reason` | 意味 |
| --- | --- |
| `"full"` | ルームに既に 2 名いる |
| `"InternalServerError"` | 認証ウェブフックの呼び出し失敗、タイムアウト、200 以外、レスポンス不正 (`allowed` 欠落、`allowed: false` なのに `reason` なし) |
| それ以外 | 認証ウェブフックが `allowed: false` と共に返した `reason` の文字列 |

#### `ping`

`{"type": "ping"}`。`websocket_ping_interval_sec` ごとに送られる。`pong` で返す。

#### `bye`

`{"type": "bye"}`。相手が切断した。受け取ったら PeerConnection を閉じ、リモートの映像要素等を破棄する。この後 Close frame (1000) が届く。

#### 中継メッセージ

相手が送った `offer` / `answer` / `candidate` / `message` が、送られてきたバイト列のまま届く。

### standalone モード

`register` に `standalone: true` を付けると、シグナリングを「接続確立のときだけ」使い、P2P 確立後は WebSocket を手放す運用になる。

- Ayame は `ping` を送らず、pong タイムアウトで切断せず、相手切断時の `bye` も送らない
- クライアントは PeerConnection の `connectionState` が `connected` になったら `{"type": "connected"}` を送り、自分で WebSocket を閉じる (Web SDK の `standalone: true` はこの動きをする)
- Ayame は `connected` を受け取るとそのセッションを終了する (Close frame を送る)。ルームは削除され、相手側のシグナリングセッションも終了する。相手が非 standalone だと相手には `bye` が届いてしまうので、**ルームの両クライアントで standalone の有無を揃える**
- `bye` が来ないので、相手の切断は PeerConnection の状態 (`disconnected` / `failed` / `closed`) で検知する

### 最小クライアント実装例 (ブラウザ)

Web SDK を使わずに直接プロトコルを話す場合の骨格。実運用では Web SDK を使うことを推奨する。

```js
const ws = new WebSocket("wss://ayame.example.com/signaling");
let pc;

ws.onopen = () => {
  ws.send(JSON.stringify({ type: "register", roomId: "room1", clientId: "alice", signalingKey: "xxx" }));
};

ws.onmessage = async (ev) => {
  const msg = JSON.parse(ev.data);
  switch (msg.type) {
    case "ping":
      ws.send(JSON.stringify({ type: "pong" }));
      break;
    case "accept":
      pc = new RTCPeerConnection({ iceServers: msg.iceServers ?? [] });
      pc.onicecandidate = (e) => {
        if (e.candidate) ws.send(JSON.stringify({ type: "candidate", ice: e.candidate.toJSON() }));
      };
      pc.ontrack = (e) => { /* e.streams[0] を表示 */ };
      // 送信するトラックがあればここで pc.addTrack(...)
      if (msg.isExistClient) {
        await pc.setLocalDescription(await pc.createOffer());
        ws.send(JSON.stringify({ type: "offer", sdp: pc.localDescription.sdp }));
      }
      break;
    case "offer":
      await pc.setRemoteDescription({ type: "offer", sdp: msg.sdp });
      await pc.setLocalDescription(await pc.createAnswer());
      ws.send(JSON.stringify({ type: "answer", sdp: pc.localDescription.sdp }));
      break;
    case "answer":
      await pc.setRemoteDescription({ type: "answer", sdp: msg.sdp });
      break;
    case "candidate":
      await pc.addIceCandidate(msg.ice);
      break;
    case "reject":
    case "bye":
      pc?.close();
      ws.close();
      break;
  }
};
```

## ウェブフック

共通仕様:

- `POST`、`Content-Type: application/json`、ボディは JSON
- `copy_websocket_header_names` で指定したヘッダーが、WebSocket ハンドシェイク時の値のまま付与される
- タイムアウトは `webhook_request_timeout_sec` (default 5 秒)
- **HTTP ステータス 200 以外はエラー扱い** (201 や 204 も不可)
- 認証ヘッダー等の仕組みは持たない。ウェブフックのエンドポイントはプライベートネットワークや認証付きリバースプロキシで保護する
- 送受信の内容は `webhook.jsonl` に記録される

### 認証ウェブフック (`authn_webhook_url`)

`register` を受け取るたびに、ルームへ登録する **前** に呼ばれる。未設定なら全員許可。シグナリングキーの検証、ルームごとの入室制御、TURN クレデンシャルの払い出しはここで行う。

リクエスト:

```json
{
  "roomId": "room1",
  "clientId": "alice",
  "connectionId": "01J5X8K2Q6W7R3N9P1T4V0Y2ZB",
  "signalingKey": "xxx",
  "authnMetadata": {"user": "alice"},
  "ayameClient": "...",
  "libwebrtc": "...",
  "environment": "..."
}
```

`roomId` / `clientId` / `connectionId` は常にある。それ以外は `register` で送られてきた場合のみ含まれる。

レスポンス (200):

| プロパティ | 型 | 必須 | 説明 |
| --- | --- | --- | --- |
| `allowed` | boolean | 必須 | 入室可否 |
| `reason` | string | `allowed: false` のとき必須 | `reject` の `reason` としてクライアントに届く |
| `authzMetadata` | 任意の JSON | 任意 | `accept` の `authzMetadata` としてクライアントに届く |
| `iceServers` | array | 任意 | `accept` の `iceServers` としてクライアントに届く。要素は `{"urls": [string, ...], "username"?: string, "credential"?: string}`。**`urls` は文字列 1 つでも配列にする** (文字列のままだとレスポンス不正扱いで `reject` になる) |

```json
{"allowed": true, "authzMetadata": {"role": "owner"}, "iceServers": [{"urls": ["stun:stun.example.com:3478"]}]}
```

```json
{"allowed": false, "reason": "InvalidSignalingKey"}
```

Ayame 側の挙動:

| ウェブフックの結果 | クライアントへ | 備考 |
| --- | --- | --- |
| `allowed: true` | `accept` | ただしルームが満員なら `reject` (`full`)。認証 OK でも入室できるとは限らない |
| `allowed: false` + `reason` | `reject` (その `reason`) | |
| `allowed: false` で `reason` なし | `reject` (`InternalServerError`) | レスポンス不正扱い |
| `allowed` がない / JSON でない | `reject` (`InternalServerError`) | |
| 200 以外 / 接続失敗 / タイムアウト | `reject` (`InternalServerError`) | `ayame.jsonl` にエラーの詳細が出る |

Python (標準ライブラリのみ) の最小実装例:

```python
import json
from http.server import BaseHTTPRequestHandler, HTTPServer

class Handler(BaseHTTPRequestHandler):
    def do_POST(self):
        req = json.loads(self.rfile.read(int(self.headers["Content-Length"])))
        if self.path == "/ayame/webhook/authn":
            if req.get("signalingKey") == "xxx":
                resp = {"allowed": True, "iceServers": [{"urls": ["stun:stun.l.google.com:19302"]}]}
            else:
                resp = {"allowed": False, "reason": "InvalidSignalingKey"}
        else:  # /ayame/webhook/disconnect
            resp = {}
        body = json.dumps(resp).encode()
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(body)

HTTPServer(("127.0.0.1", 3001), Handler).serve_forever()
```

### 切断ウェブフック (`disconnect_webhook_url`)

WebSocket 接続が終了した後に呼ばれる。**`register` 前に切れた接続や `reject` された接続でも呼ばれる** (`register` 前なら `roomId` / `clientId` は空文字列)。未設定なら何もしない。

リクエスト:

```json
{"roomId": "room1", "clientId": "alice", "connectionId": "01J5X8K2Q6W7R3N9P1T4V0Y2ZB"}
```

レスポンスは 200 であればよく、ボディは使われない (ログには残る)。失敗してもクライアントには影響せず、`ayame.jsonl` にエラーが出るだけ。

## ログ

3 種類のログを JSON Lines で出す。全てタイムスタンプは UTC (RFC 3339 ナノ秒)、`domain` フィールドでどのログかわかる。ローテーション設定 (`log_rotate_*`) は 3 つに共通。`log_stdout = true` のときは 3 つとも標準出力に混ざって出るので `domain` で分ける。

| ファイル (デフォルト) | `domain` | 内容 |
| --- | --- | --- |
| `ayame.jsonl` | `ayame` | サーバログ。起動時の設定値 (`"message": "AyameConf"` が 1 項目 1 行)、接続ごとのエラー、`debug` 時の内部状態遷移 |
| `signaling.jsonl` | `signaling` | クライアントから受信したシグナリングメッセージの記録 (`register` / `offer` / `answer` / `candidate` / `connected` / `message`)。`signaling_log_filters` で絞れる |
| `webhook.jsonl` | `webhook` | ウェブフックの送受信 (`authnReq` / `authnResp` / `disconnectReq` / `disconnectResp`) |

フィールドの注意点 (ログを集計するときに引っかかりやすい):

- `ayame.jsonl` のエラーログは `roomId` / `clientID` / `connectionId` (clientID の ID が大文字)
- `signaling.jsonl` は `roomId` / `clientId` / `connectionId` / `type`。`debug = true` のときだけ `rawMessage` (受信 JSON の文字列) が付く
- `webhook.jsonl` の `clientId` には **connectionId の値** が入る。`copyHeaders` にコピーしたヘッダーが入る
- `message` / `time` のキー名は `log_message_key_name` / `log_timestamp_key_name` で変えられる (ログ基盤の予約語と衝突するとき用)

`ayame.jsonl` に出る主なエラーメッセージ (`message` の値):

| メッセージ | 意味 |
| --- | --- |
| `InvalidJSON` / `UnexpectedJSON` | 受信メッセージが JSON として不正 |
| `InvalidMessageType` | 未知の `type`、または `type_message = false` で `message` を受信 |
| `MissingRoomID` | `register` に `roomId` がない |
| `RegistrationIncomplete` | `accept` 前に `offer` / `answer` / `candidate` / `message` / `connected` を受信 |
| `InternalServer` | 2 回目の `register` |
| `RoomFilled` | 満員で `reject` した |
| `PongTimeout` | pong が返ってこなかった |
| `AuthnWebhookError` / `AuthnWebhookResponseError` / `AuthnWebhookUnexpectedStatusCode` | 認証ウェブフックの失敗 (接続エラー / レスポンス不正 / 200 以外) |
| `DisconnectWebhookError` / `DisconnectWebhookUnexpectedStatusCode` | 切断ウェブフックの失敗 |
| `FailedToSendMsg` / `FailedWriteMessage` | クライアントへの書き込み失敗 (ほぼ切断済み) |

## Prometheus メトリクス

`http://<listen_prometheus_ipv4_address>:<listen_prometheus_port_number>/metrics` で公開する。デフォルトは `0.0.0.0:4000` なので、外部に晒さないようバインドアドレスかファイアウォールで制限する。

| メトリクス | 型 | ラベル | 説明 |
| --- | --- | --- | --- |
| `ayame_requests_total` | counter | `code`, `method`, `host`, `url` | シグナリング HTTP サーバへのリクエスト数 (`/signaling` の Upgrade、`/.ok` を含む) |
| `ayame_request_duration_seconds` | histogram | 同上 | 同リクエストの処理時間 |
| `ayame_request_size_bytes` / `ayame_response_size_bytes` | histogram | 同上 | 同リクエスト / レスポンスのサイズ |
| `ayame_webhook_requests_total` | counter | `code`, `method`, `host`, `url` | ウェブフック呼び出し回数 (認証 / 切断の両方。`url` のパスで区別) |
| `ayame_webhook_request_duration_seconds` | histogram | 同上 | ウェブフックの応答時間 |
| `ayame_webhook_request_message_size_bytes` / `ayame_webhook_response_message_size_bytes` | histogram | 同上 | ウェブフックのリクエスト / レスポンスサイズ |
| `ayame_authn_webhook_responses_total` | counter | `code`, `method`, `host`, `url`, `allowed`, `reason` | 認証ウェブフックの結果別カウント (許可 / 拒否理由の集計に使う) |

Go ランタイムの標準メトリクス (`go_*` / `process_*`) も同時に出る。

## リバースプロキシ配下での運用

ブラウザの `getUserMedia` は https 必須なので、実運用では `wss://` がほぼ必須。Ayame は TLS を持たないので前段で終端する。

nginx の例:

```nginx
location /signaling {
    proxy_pass http://127.0.0.1:3000;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Real-IP $remote_addr;
}
```

- Ayame は 5 秒ごとに `ping` を書くので、プロキシの read timeout がデフォルト (nginx は 60 秒) でも接続は維持される
- クライアントの実 IP を認証ウェブフックで使いたいときは、プロキシで `X-Forwarded-For` / `X-Real-IP` を付けて `copy_websocket_header_names = X-Forwarded-For, X-Real-IP` を設定する。ウェブフック側はリクエストヘッダーから読む (ボディには入らない)
- `/.ok` をロードバランサのヘルスチェックに使える
- ローカルで手早く wss を試すなら `ngrok http 3000`

## Ayame Labo を使う場合

自前で立てずに済ませたいときは Ayame Labo (https://ayame-labo.shiguredo.app/) が使える。Ayame と同じプロトコルで、STUN / TURN とシグナリングキー認証が組み込まれている。

- シグナリング URL: `wss://ayame-labo.shiguredo.app/signaling`
- ルーム ID は `{GitHub ログイン名}@{任意の文字列}`
- シグナリングキーはダッシュボードで取得し、`register` の `signalingKey` (Web SDK では `signalingKey` オプション) に渡す

## ハマりどころ

- **ログが出ない**: デフォルトはファイル出力。`log_stdout = true` か `debug = true` + `debug_console_log = true`
- **起動直後に落ちる**: `log_dir` が存在しない、設定ファイルに未知のキーがある、ウェブフック URL が不正、のいずれかが多い
- **`-c config.ini` が効かない**: フラグは大文字 `-C`
- **3 人目が入れない**: 仕様。1 ルーム 2 名まで。複数人なら WebRTC SFU Sora 等を使う
- **どちらが offer を出すかわからない**: `accept` の `isExistClient` が `true` の側 (後から入った側) が出す
- **60 秒で切られる**: `{"type": "ping"}` に `{"type": "pong"}` を返していない。WebSocket 制御フレームの ping/pong では代用できない
- **NAT 越えできない**: Ayame は TURN を持たない。coturn 等を用意して認証ウェブフックの `iceServers` で配る。Ayame Labo なら組み込み済み
- **誰でも接続できてしまう**: Origin チェックは常に許可。認証ウェブフックで `signalingKey` 等を検証する
- **`message` を送ると切断される**: `type_message = true` が必要
- **切断ウェブフックの `roomId` が空**: `register` 前に切れた接続。`connectionId` で `ayame.jsonl` と突き合わせる
- **Web SDK から `signalingKey` が届かない**: Web SDK は後方互換の `key` で送る。Ayame は `key` を `signalingKey` として扱い、認証ウェブフックには `signalingKey` で渡すので、ウェブフック側は `signalingKey` を見ればよい

## 関連ドキュメント

リポジトリ内:

- `README.md` — 概要・方針・利用例
- `docs/USE.md` — セットアップ手順、`register` と認証ウェブフックの説明
- `docs/DESIGN.md` — 設計メモ (ログ方針、ping/pong、goroutine 構成)
- `config_example.ini` — 全設定項目とデフォルト値
- `CHANGES.md` — リリースノート

外部:

- [ayame-spec](https://github.com/OpenAyame/ayame-spec) — シグナリング仕様とシーケンス図
- [ayame-web-sdk](https://github.com/OpenAyame/ayame-web-sdk) — ブラウザ向け SDK、[API ドキュメント](https://openayame.github.io/ayame-web-sdk/typedoc/index.html)、[オンライン DevTools](https://openayame.github.io/ayame-web-sdk/devtools/index.html)
- [ayame-web-sdk-examples](https://github.com/OpenAyame/ayame-web-sdk-examples) — sendrecv / sendonly / recvonly / datachannel のサンプル
- [Ayame Labo](https://ayame-labo.shiguredo.app/) — ホスティング版
