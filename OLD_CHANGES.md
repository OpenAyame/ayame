# 古いリリースノート

- UPDATE
    - 下位互換がある変更
- ADD
    - 下位互換がある追加
- CHANGE
    - 下位互換のない変更
- FIX
    - バグ修正

## develop

- [ADD] CI の go を 1.13 に上げる
- [UPDATE] @kdxu 推奨の go version を 1.13 にあげる
- [UPDATE] @kdxu 先に入室しているユーザーがいる場合 isExistUser をaccept時にtrue にして返す
- [FIX] @kdxu roomId か clientId が空文字列の場合は reject するように修正する
- [FIX] @kdxu -c での config ファイル指定が効いていなかったのを修正する

## 19.08.0

2019-08-16

- [UPDATE] `/ws` エンドポイントと同様のものを `/signaling` エンドポイントとして追加する
- [UPDATE] ayame register 時に key も送信できるようにする
- [UPDATE] auth webhook の返り値に iceServers があれば返却するようにする

## 19.07.1
- [CHANGE] @kdxu サンプルを ayame-web-sdk を用いたものに置き換える

## 19.07.0

- [UPDATE] @kdxu -overWsPingPong オプションで over WS の ping-pong にも対応できるようにした
- [FIX] @kdxu サンプルを unified plan に対応する
- [ADD] @kdxu ayame 起動時に少し説明を出す
- [ADD] @kdxu `ayame version` でバージョンを表示するようにする
- [ADD] @kdxu 認証ウェブフック機能を追加する
- [ADD] @kdxu 多段認証ウェブフック機能を追加する
- [CHANGE] @kdxu 設定を `config.yaml` に切り分けるよう変更する


## 19.02.1

- [FIX] @kdxu uuid を使わず、client_id で持ち回すよう修正する

## 19.02.0

**ファーストリリース**

- [ADD] ルーム機能の追加
- [ADD] type: accept/reject の追加
- [ADD} 3 人以上はキックする機能の追加
