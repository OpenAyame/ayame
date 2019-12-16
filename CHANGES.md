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

- [UPDATE] 推奨の go version を 1.13 にあげる
    - @kdxu 
- [UPDATE] 先に入室しているユーザーがいる場合 isExistUser をaccept時にtrue にして返す
    - @kdxu 
- [UPDATE] ドキュメントを最新のシグナリングの内容にする
    - @Hexa
- [ADD] CI の go を 1.13 に上げる
    - @kdxu 
- [CHANGE] デフォルトで websocket での ping/pong を有効にする
    - @Hexa
- [CHANGE] 多段ウェブフック認証を削除する
    - @Hexa
- [CHANGE] シグナリングの close を bye に変更する
    - AppRTC に揃える
    - @Hexa
- [CHANGE] webhook レスポンスの authWebhookUrl を削除する
    - @Hexa
- [CHANGE] webhook リクエストの key を signalingKey に変更する
    - @Hexa
- [CHANGE] webhook のリクエストに clientId, authnMetadata を追加する
    - @Hexa
- [CHANGE] 多段ウェブフック認証時のリクエストから host を削除する
    - @Hexa
- [CHANGE] 認証ウェブフック機能時の JSON を lowerCamelCase で統一する
    - @Hexa
- [FIX] WS 切断時に Bye を送っていなかったのを修正する
    - @Hexa
- [FIX] roomId か clientId が空文字列の場合は reject するように修正する
    - @kdxu 
- [FIX] -c での config ファイル指定が効いていなかったのを修正する
    - @kdxu 
- [FIX] AllowOrigin のサンプルを "" で囲う
    - @Hexa
- [FIX] Upgrade の前に CheckOrigin を定義して Origin チェックを有効にする
    - @Hexa
- [FIX] origin ヘッダがない場合は allow_origin の値にかかわらず通す
    - @Hexa
- [FIX] ドキュメントの特殊文字を半角スペースに置き換える
    - @Hexa
- [FIX] log_dir にディレクトリがない場合は終了する
    - @Hexa
- [FIX] webhook 時のレスポンスに allowed がない場合は認証を拒否する
    - @Hexa
- [FIX] 認証失敗時に reason が含まれていない場合は INTERNAL-ERROR を返す

## 19.08.0

- [UPDATE] `/ws` エンドポイントと同様のものを `/signaling` エンドポイントとして追加する
    - @kdxu 
- [UPDATE] ayame register 時に key も送信できるようにする
    - @kdxu 
- [UPDATE] auth webhook の返り値に iceServers があれば返却するようにする
    - @kdxu 

## 19.07.1

- [CHANGE] サンプルを ayame-web-sdk を用いたものに置き換える
    - @kdxu 

## 19.07.0

- [UPDATE] -overWsPingPong オプションで over WS の ping-pong にも対応できるようにした
    - @kdxu 
- [FIX] サンプルを unified plan に対応する
    - @kdxu 
- [ADD] ayame 起動時に少し説明を出す
    - @kdxu 
- [ADD] `ayame version` でバージョンを表示するようにする
    - @kdxu 
- [ADD] 認証ウェブフック機能を追加する
    - @kdxu 
- [ADD] 多段認証ウェブフック機能を追加する
    - @kdxu 
- [CHANGE] 設定を `config.yaml` に切り分けるよう変更する
    - @kdxu 

## 19.02.1

- [FIX] uuid を使わず、client_id で持ち回すよう修正する
    - @kdxu 

## 19.02.0

**ファーストリリース**

- [ADD] AppRTC 互換のシグナリングサーバの追加
    - @kdxu 
- [ADD] ルーム機能の追加
    - @kdxu 
- [ADD] type: accept/reject の追加
    - @kdxu 
- [ADD] 3 人以上はキックする機能の追加
    - @kdxu 
