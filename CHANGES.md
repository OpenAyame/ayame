# リリースノート

- UPDATE
    - 下位互換がある変更
- ADD
    - 下位互換がある追加
- CHANGE
    - 下位互換のない変更
- FIX
    - バグ修正

## develop

## 2020.1.4

- [ADD] 起動時に INFO で Ayame のバージョンをログに書く仕組みを追加
- [ADD] GitHub Actions でリリースファイルをアップロードする仕組みを追加
    - 古い仕組みを整理
- [UPDATE] go.mod を Go 1.15 に上げる

## 2020.1.3

- [UPDATE] rs/zerolog を v1.19.0 に上げる
- [UPDATE] gorilla/websocket を v1.4.2 に上げる
- [UPDATE] yaml を v2.3.0 に上げる

## 2020.1.2

**昔にリリースミスが発覚したため、master を 19.08.0 まで戻してから再度リリースを行った**

## 2020.1.1

- [FIX] 受信したメッセージが null の場合に落ちるため、nil チェックを追加する
    - @kadoshita @Hexa

## 2020.1

- [ADD] register メッセージで key と signalingKey のどちらかを指定できるようにする
    - signalingKey が優先される
    - 将来的に signalingKey のみになる
    - @voluntas
- [ADD] accept メッセージで isExistUser 以外に isExistClient を送るようにする
    - 将来的に isExistClient のみになる
    - @voluntas
- [ADD] 切断時にウェブフック通知を飛ばせるようにする
    - disconnect_webhook_url を設定
    - @voluntas @Hexa
- [ADD] signaling.log を追加する
    - @voluntas @Hexa
- [ADD] webhook.log を追加する
    - @voluntas @Hexa
- [ADD] register メッセージで ayameClient / environment / libwebrtc の情報を追加する
    - 認証ウェブフック通知で含まれるようにする
    - @voluntas
- [ADD] type: accept 時に connectionId を払い出すようにする
    - @voluntas
- [CHANGE] コードベースを変更する
    - @voluntas @Hexa
- [CHANGE] addr を listen_ipv4_address に変更する
    - @voluntas
- [CHANGE] port を listen_port_number に変更する
    - @voluntas
- [CHANGE] allow_origin 設定を削除する
    - @voluntas
- [CHANGE] ロガーを zerolog に変更する
    - @voluntas
- [CHANGE] ログローテーションを lumberjack に変更する
    - @voluntas
- [CHANGE] サンプルを削除する
    - @voluntas
- [CHANGE] 登録済みのあとに WebSocket 切断した場合、 type: bye を送信するようにする
    - @voluntas @Hexa
- [CHANGE] ウェブフックの戻り値のステータスコード 200 以外はエラーにする
    - @voluntas @Hexa
- [CHANGE] ウェブフックの JSON のキーを snake_case から camelCase にする
    - @voluntas
- [CHANGE] clientId をオプション化する
    - @voluntas
- [FIX] サーバ側の切断の WS の終了処理を適切に行う
    - @voluntas @Hexa
- [FIX] ウェブソケットの最大メッセージを 1MB に制限する
    - @voluntas
- [FIX] ayame.log にターミナル用のカラーコードを含めないようにする
    - @Hexa
- [CHANGE] 指定したログレベルでの ayamne.log へのログ出力に対応する
    - @Hexa
