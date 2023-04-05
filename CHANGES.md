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

- [ADD] standalone モードを追加する
    - standalone モード時に、type: connected を受信した場合は WebSocket を切断する
    - standalone モードでは WebSocket を切断するため、PING を送信しないようにする
    - @Hexa
- [ADD] ヘルスチェック用の URL を追加する
    - @Hexa
- [CHANGE] 設定ファイルを yaml から ini に変更する
    - @Hexa
- [FIX] webhook log が出力されるように修正する
    - @Hexa
- [UPADTE] handler に echo を使用するように変更する
    - @Hexa
- [ADD] Prometheus に対応する
    - @Hexa


## 2022.2.0

- [CHANGE] ログに github.com/rs/zerolog/log を利用するように変更する
    - @voluntas
- [CHANGE] lumberjack を shiguredo/lumberjack/v3 に変更する
    - @voluntas
- [CHANGE] websocket を shiguredo/websocket v1.6.0 に変更する
    - @voluntas

## 2022.1.3

- [FIX] リリースバイナリが生成されないのを修正する
    - @voluntas

## 2022.1.2

- [FIX] リリースバイナリが生成されないのを修正する
    - @voluntas

## 2022.1.1

- [CHANGE] gox を利用しないリリース方式に変更する
    - @voluntas
- [UPDATE] ビルドテストを Go 1.18 以上にする
    - @voluntas

## 2022.1.0

- [UPDATE] actions/checkout@v3 に上げる
    - @voluntas
- [UPDATE] GitHub Actions の Go を 1.18 に上げる
    - @voluntas
- [UPDATE] toml を v1.0.0 に上げる
    - @voluntas
- [UPDATE] rs/zerolog を v1.26.1 に上げる
    - @voluntas

## 2021.2.1

- [UPDATE] GitHub Actions の Go を 1.17.1 に上げる
    - @voluntas
- [UPDATE] rs/zerolog を v1.25.0 に上げる
    - @voluntas
- [UPDATE] yaml.v3 に戻す
    - @voluntas

## 2021.2

- [UPDATE] go.mod を Go 1.17 に上げる
    - @voluntas
- [UPDATE] GitHub Actions の Go を 1.17 に上げる
    - @voluntas
- [UPDATE] rs/zerolog を v1.21.0 に上げる
    - @voluntas
- [UPDATE] rs/zerolog を v1.24.0 に上げる
    - @voluntas
- [CHANGE] "github.com/goccy/go-yam" に変更する
    - @voluntas

## 2021.1

- [ADD] GitHub Actions の Go を 1.16 に上げる
    - @voluntas
- [UPDATE] go.mod を Go 1.16 に上げる
    - @voluntas

## 2020.1.5

- [UPDATE] rs/zerolog を v1.20.0 に上げる
    - @voluntas
- [UPDATE] yaml を v2.4.0 に上げる
    - @voluntas

## 2020.1.4

- [ADD] 起動時に INFO で Ayame のバージョンをログに書く仕組みを追加
    - @voluntas
- [ADD] GitHub Actions でリリースファイルをアップロードする仕組みを追加
    - 古い仕組みを整理
    - @voluntas
- [UPDATE] go.mod を Go 1.15 に上げる
    - @voluntas

## 2020.1.3

- [UPDATE] rs/zerolog を v1.19.0 に上げる
    - @voluntas
- [UPDATE] gorilla/websocket を v1.4.2 に上げる
    - @voluntas
- [UPDATE] yaml を v2.3.0 に上げる
    - @voluntas

## 2020.1.2

**昔にリリースミスが発覚したため、master を 19.08.0 まで戻してから再度リリースを行った**

## 2020.1.1

- [FIX] 受信したメッセージが null の場合に落ちるため、nil チェックを追加する
    - @kadoshita @Hexa

## 2020.1

**全て 1 から書き直している**

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
