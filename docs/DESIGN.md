# Ayame 設計デザイン

## 設計

- ログの時間は UTC に固定する
- 切断が発生する場合のログレベルはエラー
- 切断しないけどログに残しておきたいときはワーニング
- Webhook 関連でエラーになった時はエラーをできるだけ詳細に出す
- クライアント側でできるだけ処理をする
- offer / answer / candidate は送られてきた JSON をそのまま転送する
  - forward の中身は connection, rawMessage を用意しておいてそれを書き込むだけにする
- unregister の成功は forward チャネルのクローズをトリガーとする
- チャネル利用時はインターフェースは使わない
- 認証ウェブフックは register 前に行う
  - signalingKey のチェックなどもここで行う
  - 認証が成功したら register 処理を走らせる
    - その際にすでに他の人がマッチしてたらキックされる
    - 認証が成功したとしても接続できる状態になるとは限らない
- 認証ウェブフックが ok になったら、register を試みる
  - 認証ウェブフックがなければそのまま処理に進む
- 終了は丁寧に終了する
- SDK サンプルがあるので、サンプルの提供はしない
- main, wsRecv の２つのプロセスが動く
  - wsRecv はとにかく WS でメッセージを受け取って main にわたすだけ
  - wsRecv は main が死んだ時用に ctx を共有しておく
- ping/pong は ping を 5 秒間隔でなげて 60 秒 pong が返ってこなかったら切断する
- シグナリングの送受信ログを取る
  - roomId / clientId / connectionId を突っ込む
  - JSON 形式でとる
  - 送られてきたメッセージをそのまま書き出す
- ログローテーションはすべてのログに共通にする
  - 個別には対応しない
- 認証ウェブフックのログを取る
  - 送信、受信ログを取る
- クライアント ID はオプション化する
- 接続の名前付けはコネクション ID (ULID) を利用する
- type: accept でコネクション ID を払い出す

## 利用ライブラリ

- WS は gorilla/websocket
  - https://github.com/gorilla/websocket
- ログは zerolog
  - https://github.com/rs/zerolog
- ログローテは lumberjack
  - https://github.com/natefinch/lumberjack
- コネクション ID は ulid
  - https://github.com/oklog/ulid

## 検討

- API の利用を想定する
  - 指定した roomId を切断する
- ConnectionID を UUIDv4 + Clockford Base32 に切り替える
  - ULID の意味があまりない
