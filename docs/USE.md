# Ayame を使ってみる

## ビルドする

まずはこのリポジトリをクローンします。

### Go のインストール

推奨バージョンは以下のようになります。
```
go 1.14
```

### ビルドする

```
$ make
```

## 設定ファイルを生成する

```
$ make init
```

## サーバを起動する

ビルドに成功したら、以下のコマンドで Ayame サーバーを起動することができます。

```
$ ./ayame
```

## Ayame Web SDK サンプルを利用して動作確認をする

Ayame Web SDK のサンプルを利用することで動作を確認できます。

```
$ git clone git@github.com:OpenAyame/ayame-web-sdk-samples.git
$ cd ayame-web-sdk-samples
$ yarn install
```

main.js URL を `'ws://127.0.0.1:3000/signaling'` に変更

https://github.com/OpenAyame/ayame-web-sdk-samples/blob/master/main.js#L1

```
$ yarn serve
```

http://127.0.0.1:5000/sendrecv.html をブラウザタブで２つ開いて接続を押してみてください。


## コマンド

```
$ ./ayame version
WebRTC Signaling Server Ayame version 2020.1.2
```

```
$ ./ayame
2020-01-08 07:04:58.392536Z [INFO] AyameConf debug=true
2020-01-08 07:04:58.392685Z [INFO] AyameConf log_dir=.
2020-01-08 07:04:58.392714Z [INFO] AyameConf log_name=ayame.log
2020-01-08 07:04:58.392737Z [INFO] AyameConf log_level=debug
2020-01-08 07:04:58.392761Z [INFO] AyameConf signaling_log_name=signaling.log
2020-01-08 07:04:58.392781Z [INFO] AyameConf listen_ipv4_address=0.0.0.0
2020-01-08 07:04:58.392803Z [INFO] AyameConf listen_port_number=3000
2020-01-08 07:04:58.392829Z [INFO] AyameConf authn_webhook_url=
2020-01-08 07:04:58.392847Z [INFO] AyameConf disconnect_webhook_url=
2020-01-08 07:04:58.392868Z [INFO] AyameConf webhook_log_name=webhook.log
2020-01-08 07:04:58.392891Z [INFO] AyameConf webhook_request_timeout_sec=5
```

```
$ ./ayame -help
Usage of ./ayame:
  -c string
    	ayame の設定ファイルへのパス(yaml) (default "./ayame.yaml")
```

## `register` メッセージについて

クライアントは ayame への接続可否を問い合わせるために WebSocket に接続した際に、まず `"type": "register"` の JSON メッセージを WS で送信する必要があります。
register で送信できるプロパティは以下になります。

- `"type"`: (string): 必須。 `"register"` を指定する
- `"clientId"`: (string): 必須
- `"roomId"`: (string): 必須
- `"signalingkey"`(string): オプション
- `"authnMetadata"`(object): オプション
- `"ayameClient"`(string): オプション
- `"environment"`(string): オプション
- `"libwebrtc"`(string): オプション

## 認証ウェブフックの `auth_webhook_url` オプションについて

`ayame.yaml` にて `auth_webhook_url` を指定している場合、
ayame は client が `{"type": "register" }` メッセージを送信してきた際に `ayame.yaml` に指定した `auth_webhook_url` に対して認証リクエストを JSON 形式で POST します。

また、 認証リクエストの返り値は JSON 形式で、以下のように想定されています。

- `"allowed"`: (boolean): 必須。認証の可否
- `"reason"`: (string): オプション。認証不可の際の理由 (`allowed` が false の場合のみ必須)
- `"iceServers"`: (array object): オプション。クライアントに peer connection で接続する iceServer 情報

`allowed` が false の場合 client の ayame への WebSocket 接続は切断されます。

#### リクエスト

- `"clientId"`: (string): 必須
- `"roomId"`: (string): 必須
- `"signalingkey"`(string): オプション
- `"authnMetadata"`: (object): オプション
    - register 時に `authnMetadata` をプロパティとして指定していると、その値がそのまま付与されます
- `"ayameClient"`(string): オプション
- `"environment"`(string): オプション
- `"libwebrtc"`(string): オプション

#### レスポンス

- `"allowed"`: (boolean): 必須。認証の可否
- `"reason"`: (string): 認証不可の際の理由 (`allowed` が false の場合のみ)
- `"authzMetadata"`(object): オプション
    - クライアントに対して任意に払い出せるメタデータ
    -  client はこの値を読み込むことで、例えば username を認証サーバから送ったりということも可能になる

```
{"allowed": true, "authzMetadata": {"username": "ayame", "owner": "true"}}
```

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
