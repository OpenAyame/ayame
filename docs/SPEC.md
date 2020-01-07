# Ayame 要求仕様書

## 設定ファイル

**ayame.conf**

```
log_dir: .
log_name: ayame.log
log_level: info

signaling_log_name: signaling.log
debug: true
log_dir: .
log_name: ayame.log
log_level: debug

signaling_log_name: signaling.log
webhook_log_name: webhook.log

listen_ipv4_address: 127.0.0.1
listen_port_number: 3000

# authn_webhook_url: http://127.0.0.1:3001/authn_webhook_url
# disconnect_webhook_url: http://127.0.0.1:3001/disconnect_webhook_url
# webhook_request_timeout: 5

# allow_origin: "*.example.com"
```

- allow_origin: は "" で囲う必要あり

## シグナリング

- 基本的には AppRTC へ準拠する
- 必要があれば拡張する

- 登録
    - type: register
        - roomId
            - 必須
        - clientId
            - オプション
                - 要検討
        - authnMetadata
            - オプション
            - any
        - signalingKey
            - オプション
            - 互換性のため key も許可する
                - 2020.2 で互換性はなくす
            - signalingKey と key 両方飛んできたら signalingKey を優先する
- 切断
    - type: bye
        - 1:1 のどちらかが切断したら飛ばす
- 認証成功
    - 拡張
    - type: accept
        - authzMetadata: interface
        - iceServers:
            - stun/turn の払い出し
        - isExistClient: bool
            - isExistUser も 2020.2 までは飛ばす

- 認証拒否
    - 拡張
    - type: reject
        - reason
- ピンポン
    - 拡張
    - type: ping

## ウェブフック

- URL の設定は yaml に設定可能にする

### 認証サーバへ飛ばす情報

- clientId
    - 必須
    - string
- roomId
    - 必須
    - string
- authnMetadata
    - オプション
    - any
- signalingKey
    - オプション
    - string
- ayameClient
    - オプション
    - string
- environment
    - オプション
    - string
- libwebrtc
    - オプション
    - string

### 認証サーバから払い出す情報

- allowed
    - 必須
    - boolean
- reason
    - オプション
    - allowed が false のときのみ必須となる
    - string
- authzMetadata
    - オプション
    - クライアントまで届く
    - any
- iceServers:
    - オプション
    - 構成は考える

### シグナリング切断時に飛ばす情報

- roomId
    - 必須
    - string
- clientId
    - 必須
    - string