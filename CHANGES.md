# リリースノート

- UPDATE
    - 下位互換がある変更
- ADD
    - 下位互換がある追加
- CHANGE
    - 下位互換のない変更
- FIX
    - バグ修正

## feature/rewrite

- [ADD] register メッセージで key と signalingKey のどちらかを指定できるようにする
    - signalingKey が優先される
    - 将来的に signalingKey のみになる
- [ADD] accept メッセージで isExistUser 以外に isExistClient を送るようにする
    - 将来的に isExistClient のみになる
- [ADD] 切断時にウェブフック通知を飛ばせるようにする
    - disconnect_webhook_url を設定
- [ADD] signaling.log を追加する
- [ADD] webhook.log を追加する
- [ADD] register メッセージで ayameClient / environment / libwebrtc の情報を追加する
    - 認証ウェブフック通知で含まれるようにする
- [CHANGE] コードベースを変更する
- [CHANGE] addr を listen_ipv4_address に変更する
- [CHANGE] port を listen_port_number に変更する
- [CHANGE] allow_origin 設定を削除する
- [CHANGE] ロガーを zerolog に変更する
- [CHANGE] ログローテーションを lumberjack に変更する
- [CHANGE] 同一クライアント ID での同一ルームへの接続をできなくする
- [CHANGE] サンプルを削除する
- [CHANGE] 登録済みのあとに WebSocket 切断した場合、 type: bye を送信するようにする
- [CHANGE] ウェブフックの戻り値のステータスコード 200 以外はエラーにする
- [FIX] サーバ側の切断の WS の終了処理を適切に行う
- [FIX] ウェブソケットの最大メッセージを 1MB に制限する