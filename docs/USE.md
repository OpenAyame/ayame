# Ayame を使ってみる

## セットアップ

まずはこのリポジトリをクローンします。

### Go のインストール

推奨バージョンは以下のようになります。

`go 1.27.2`

### パッケージのインストール

```bash
go install
```

### ビルドする

```bash
make
```

## 設定ファイルを生成する

```bash
make init
```

## サーバを起動する

> [!WARNING]  
> デフォルトではコンソールにログは出力されません。

```bash
./bin/ayame -C config.ini
```

## Ayame Web SDK サンプルを利用して動作確認をする

Ayame Web SDK のサンプルを利用することで動作を確認できます。

サンプルは <https://github.com/OpenAyame/ayame-web-sdk-examples> で公開されています。
実行には Node.js 20 以上と npm を利用します。

```bash
git clone https://github.com/OpenAyame/ayame-web-sdk-examples.git
cd ayame-web-sdk-examples
npm install
```

`.env.template` を `.env.local` にコピーし、起動した Ayame に合わせて設定します。

```bash
cp .env.template .env.local
```

```bash
VITE_AYAME_SIGNALING_URL=ws://127.0.0.1:3000/signaling
VITE_AYAME_ROOM_ID_PREFIX=ayame@
VITE_AYAME_ROOM_NAME=examples
VITE_AYAME_SIGNALING_KEY=
```

- `VITE_AYAME_SIGNALING_URL`: Ayame のシグナリング URL
- `VITE_AYAME_ROOM_ID_PREFIX` / `VITE_AYAME_ROOM_NAME`: ルーム ID はこの 2 つを連結した値
- `VITE_AYAME_SIGNALING_KEY`: 認証ウェブフックを設定していない場合には使われないため空

```bash
npm run dev
```

`npm run dev` を実行すると <http://localhost:5173/> で起動します。
<http://localhost:5173/sendrecv> をブラウザのタブで 2 つ開き、それぞれ Connect を押してください。
ルーム名を変更する場合は、2 つのタブで同じルーム名を指定します。

## コマンド

```bash
$ ./bin/ayame -V
WebRTC Signaling Server Ayame version 2026.2.2
```

```bash
$ ./bin/ayame -help
Usage of ./bin/ayame:
  -C string
    	設定ファイルへのパス (default "./config.ini")
  -V	バージョン
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

## 認証ウェブフックの `authn_webhook_url` オプションについて

`ayame.ini` にて `authn_webhook_url` を指定している場合、
ayame は client が `{"type": "register" }` メッセージを送信してきた際に `config.ini` に指定した `authn_webhook_url` に対して認証リクエストを JSON 形式で POST します。

また、 認証リクエストの返り値は JSON 形式で、以下のように想定されています。

- `"allowed"`: (boolean): 必須。認証の可否
- `"reason"`: (string): オプション。認証不可の際の理由 (`allowed` が false の場合のみ必須)
- `"iceServers"`: (array object): オプション。クライアントに peer connection で接続する iceServer 情報

`allowed` が false の場合 client の ayame への WebSocket 接続は切断されます。

### 認証ウェブフックリクエスト

- `"clientId"`: (string): 必須
- `"roomId"`: (string): 必須
- `"signalingkey"`(string): オプション
- `"authnMetadata"`: (object): オプション
  - register 時に `authnMetadata` をプロパティとして指定していると、その値がそのまま付与されます
- `"ayameClient"`(string): オプション
- `"environment"`(string): オプション
- `"libwebrtc"`(string): オプション

### 認証ウェブフックレスポンス

- `"allowed"`: (boolean): 必須。認証の可否
- `"reason"`: (string): 認証不可の際の理由 (`allowed` が false の場合のみ)
- `"authzMetadata"`(object): オプション
  - クライアントに対して任意に払い出せるメタデータ
  - client はこの値を読み込むことで、例えば username を認証サーバから送ったりということも可能になる

```json
{"allowed": true, "authzMetadata": {"username": "ayame", "owner": "true"}}
```

### ローカルで wss/https を試したい場合

[ngrok \- secure introspectable tunnels to localhost](https://ngrok.com/) の使用を推奨しています。

```bash
$ ngrok http 3000
ngrok by @xxxxx
Session Status online
Account        xxxxx
Forwarding     http://xxxxx.ngrok.io -> localhost:3000
Forwarding     https://xxxxx.ngrok.io -> localhost:3000
```
