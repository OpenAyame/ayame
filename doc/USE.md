## Ayame を使ってみる

まずはこのリポジトリをクローンします。
ディレクトリ構成は以下のようになります。

```
$ ./nodejs/
.
├── public/
│   ├── index.html
│   ├── main.css
│   └── webrtc.js
├── src/
│   ├── server.ts
│   └── ws.ts
├── yarn.lock
└── package.json
```


## node, yarn のインストール

推奨バージョンは 2019/2/12 時点で以下のようになります。
```
node 11.9.0
yarn 1.13.0
```

### node

https://nodejs.org/ja/download/ から node (最新版 11.9.0 を推奨します) をインストールしてください。

### yarn

https://yarnpkg.com/lang/ja/docs/install から yarn をインストールしてください。(1.13.0)

Mac の場合 homebrew 経由でのインストールが可能です。

## 依存パッケージのインストール

yarn を利用します。

```
$ yarn install
```

## サーバーを起動する

依存パッケージのインストールに成功したら、以下のコマンドで Ayame サーバーを起動することができます。

```
$ yarn start:sample
```

起動したら、 http://localhost:3000 にアクセスすることでデモ画面にアクセスできます。

アクセス時に各ブラウザで「カメラ・マイクでのアクセス」権限を要求された場合は「許可する」を選択してください。

権限を確認できたら、「接続する」を選択してください。

別のタブ or ブラウザから同様にアクセスして、互いの画面が表示されたら接続成功です。

※ あくまで Peer 2 Peer なので、最大 2 クライアントまでの接続しかできません。

切断するときは「切断する」を選択してください。


## デモ無しでサーバを起動する

```
$ yarn start
```

このオプションで起動した場合、http://localhost:3000 にアクセスしても、デモ画面は表示されません。
あくまでシグナリングサーバの機能のみ用いたい場合に、こちらのオプションで起動します。

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

