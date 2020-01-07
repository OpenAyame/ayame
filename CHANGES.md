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
    - signalingKey が優先されます
    - 将来的に signalingKey のみになります
- [ADD] accept メッセージで isExistUser 以外に isExistClient を送るようにする
    - 将来的に isExistClient のみになります
- [ADD] signaling.log を追加する
- [ADD] webhook.log を追加する
- [CHANGE] addr を listen_ipv4_address に変更する
- [CHANGE] port を listen_port_number に変更する
- [CHANGE] allow_origin 設定を削除する
- [CHANGE] ロガーを zerolog に変更する
- [CHANGE] ログローテーションを lumberjack に変更する
- [CHANGE] サンプルを削除する
- [FIX] サーバ側の切断の WS の終了処理を適切に行う
