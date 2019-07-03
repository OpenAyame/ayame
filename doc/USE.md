## Ayame を使ってみる

まずはこのリポジトリをクローンします。
ディレクトリ構成は以下のようになります。

```
$ ./
.
├── sample/
│   ├── index.html
│   ├── main.css
│   └── webrtc.js
├── .doc/
│   ├── USE.md
│   └── DETAIL.md
├── go.mod
├── go.sum
├── ws_handler.go
├── client.go
├── hub.go
└── main.go
```


## Go のインストール

推奨バージョンは以下のようになります。
```
go 1.12
```

## ビルドする

```
$ go build
```

`make` でもビルド出来ます。

```
$ make
```

## サーバを起動する

ビルドに成功したら、以下のコマンドで Ayame サーバーを起動することができます。

```
$ ./ayame
```

起動したら、 http://localhost:3000 にアクセスすることでデモ画面にアクセスできます。

アクセス時に各ブラウザで「カメラ・マイクでのアクセス」権限を要求された場合は「許可する」を選択してください。

権限を確認できたら、「接続する」を選択してください。

別のタブ or ブラウザから同様にアクセスして、互いの画面が表示されたら接続成功です。

※ あくまで Peer 2 Peer なので、最大 2 クライアントまでの接続しかできません。

切断するときは「切断する」を選択してください。

## コマンド


```
$ ./ayame version
WebRTC Signaling Server Ayame version 19.02.1⏎
```

```
$ ./ayame -c ./config.yaml
time="2019-06-10T00:23:16+09:00" level=info msg="Setup log finished."
time="2019-06-10T00:23:16+09:00" level=info msg="WebRTC Signaling Server Ayame. version=19.02.1"
time="2019-06-10T00:23:16+09:00" level=info msg="running on http://localhost:3000 (Press Ctrl+C quit)"
```

```
$ ./ayame -help
Usage of ./ayame:
  -c string
    	ayame の設定ファイルへのパス(yaml) (default "./config.yaml")
```

## `over_ws_ping_pong` オプションについて

- `config.yaml` にて `over_ws_ping_pong: true` に設定した場合、 ayame はクライアントに対して(WebSocket の ping frame の代わりに) ** 9 ** 秒おきに JSON 形式で `{"type": "ping"}` メッセージを送信します。
- これに対してクライアントは ** 10 ** 秒以内に JSON 形式で `{"type": "pong"}` を返すことで ping-pong を実現します。

クライアント(javascript) のサンプルコードを以下に示します。

```javascript
ws = new WebSocket(signalingUrl);
ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      console.log(message.type)
      switch(message.type){
        case 'ping': {
          console.log('Received Ping, Send Pong.');
          ws.send(JSON.stringify({
            "type": "pong"
          }))
          break;
        }
        ...
```


## `use_auth_webhook` オプションについて

`config.yaml` にて `use_auth_webhook: true` に設定した場合、 ayame は client が {"type": "register" } メッセージを送信してきた際に `config.yaml` に指定した `auth_webhook_url` に対して認証リクエストをJSON 形式で POST します。


このとき、{"type": "register" } のメッセージに

- `"metadata"`(string)
- `"key"`(string)

を含めていると、そのデータを ayame はそのまま指定した `auth_webhook_url` に JSON 形式で送信します。


また、 auth webhook の返り値は JSON 形式で、以下のように想定されています。

- `allowed`: boolean。認証の可否
- `reason`: string。認証不可の際の理由

`allowed` が false の場合 client の ayame への WebSocket 接続は切断されます。


### ローカルで wss/https を試したい場合

[ngrok \- secure introspectable tunnels to localhost](https://ngrok.com/) の使用を推奨しています。

```
$ ngrok http 3000
ngrok by @xxxxx
Session Status online
Account        xxxxx
Forwarding     http://xxxxx.ngrok.io -> localhost:3000
Forwarding     https://xxxxx.ngrok.io -> localhost:3000
```

