# Ayame 技術概要

この文書は Ayame のシグナリングサーバ及びデモアプリケーションがどのように動作しているかを説明することを目的とします。

## 互換性

このシグナリングサーバは [webrtc/apprtc: The video chat demo app based on WebRTC](https://github.com/webrtc/apprtc) と互換性があります。

## 設計

- ログの時間は UTC に固定する
- 切断が発生する場合のログレベルはエラー
- 切断しないけどログに残しておきたいときはワーニング
- Webhook 関連でエラーになった時はエラーをできるだけ詳細に出す
- クライアント側でできるだけ処理をする
- offer / answer / candidate は forward する
    - forward メッセージを ayame.go に渡してスルーしていく
- register の戻りは result (int) で受け取る
- unregister は戻り値はなし、送って受信したら終了
- forward の中身は client, raw_message を用意しておいてそれを書き込むだけにする
- チャネル利用時はインターフェースはできるだけ使わない
- クライアントとは ping/pong を独自に行う
- API の利用を想定する
    - 指定した RoomID の切断 API を用意する
- 認証ウェブフックは signaling_handler 部分で行う
    - signalingKey のチェックなども個々で行う
    - 認証が成功したら register 処理を走らせる
        - その際にすでに他の人がマッチしてたらキックされる
        - 認証が成功したとしても接続できる状態になるとは限らない
- 認証ウェブフックが ok になったら、register を試みる
    - signalingKey がない場合は直接 register を試みる
- 終了は丁寧に終了する
    - どうするか要検討
    - 作ってみてから考える
- SDK サンプルがあるので、サンプルの提供はしない
- main, wsRecv の２つのプロセスがぐるぐる動く
    - wsRecv はとにかくバイナリを受け取って main にわたすだけ
    - wsRecv は main が死んだ時用に ctx を共有しておく
- ping/pong は ping を 5 秒間隔でなげて 60 秒 pong が返ってこなかったら切断する
    - Sora 方式を採用する
- シグナリングの送受信ログを取る
    - roomId / clientId を突っ込む
    - JSON 形式でとる
    - 送られてきたメッセージをそのまま書き出す
- ログローテーションはすべてのログに共通にする
    - 個別には対応しない
- 認証ウェブフックのログを取る
    - 送信、受信ログを取る

## 利用ライブラリ

- WS は gorilla/websocket
- ログは zerolog
- ログローテは lumberjack

## 検討

- クライアント ID はオプションにする
    - 送ってこなかったら動的生成
