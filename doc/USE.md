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
$ GO111MODULE=on go get -u && go build
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

