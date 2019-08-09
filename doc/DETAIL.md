## Ayame 技術概要

この文書は Ayame のシグナリングサーバ及びデモアプリケーションがどのように動作しているかを説明することを目的とします。

### 互換性

このシグナリングサーバは [webrtc/apprtc: The video chat demo app based on WebRTC](https://github.com/webrtc/apprtc) と互換性があります。

### シグナリングについて

Ayame Server の `ws://localhost:3000/signaling` がクライアントからの WebSocket 接続を確立し、管理するエンドポイントとなります。

このエンドポイントに WebSocket で接続すると、Ayame Server は接続したクライアントを保持します。

このシグナリングサーバは 1 対 1 専用のため、3 つ以上のクライアントの接続要求は拒否します。

Ayame は WebSocket で接続しているクライアントのうちどれかからデータが来ると、送信元のクライアント以外の接続済みのクライアントにデータを*そのまま* WebSocket で送信します。これらはすべて非同期で行われます。

これが「シグナリング」です。

### 接続確立までのシーケンス図

Ayame Server が互いのSDP 交換や peer connection の接続をシグナリングによってやり取りします。

SDP とは WebRTC の接続に必要な peer connection の 内部情報です。 

- [RFC 4566 \- SDP: Session Description Protocol](https://tools.ietf.org/html/rfc4566)
- [Annotated Example SDP for WebRTC](https://tools.ietf.org/html/draft-ietf-rtcweb-sdp-11)

 ```

  +-------------+     +-------------------+    +-------------+
  |   browser1  |     |   Ayame Server    |    |   browser2  |
  +-----+-------+     +--------+----------+    +------+------+
        |                      |                      |
    ----------------------WebSocket 接続確立----------------
        +--------------------->|                      |
        |      websocket 接続  | <--------------------+
        |                      |     websocket 接続   |
    -----------------Peer-Connection の初期化---------------
        |                      |                      |
        | getUserMedia()       |                      | getUserMedia() 
        | localStream の取得   |                      | localStream の取得 
        | peer = new PeerConnection()                 | peer = new PeerConnection()
    -----------------クライアント情報の登録----------------------
        +--------------------->|                      | room の id と client のid を登録する
        |   ws message         |                      | 2 人以下で入室可能であれば ayame は accept を返却
        |   {type: register,   |                      | それ以外の場合 reject を返却
        |   roomId: roomId,  |                      | TURN などのメタデータも将来的にここで交換する
        |   client: clientId} |                      |
        |<---------------------+                      |
        |  {type: accept }     |<---------------------|
        |                      |   ws message         | 
        |                      |    register          |  
        |                      |--------------------->|
    -----------------------SDP の交換-----------------------
        |                      |                      |
        + peer.createOffer(),  |                      |
        | peer.setLocalDescription()                  |
        |  offerSDP を生成     |                      |
        |                      |                      |
        +--------------------->|                      |
        |      ws message      |--------------------> |
        |      offerSDP        |   ws message         | offerSDP をもとに Remote Description をセット
        |                      |    offerSDP          |  answerSDP を生成し、それをもとに localDescription を生成する
        |                      |                      |　peer.setRemoteDescription(offerSDP),
        |                      |                      |  peer.createAnswer(),
        |                      | <--------------------+  peer.setLocalDescription(answer)
        | <--------------------+    ws message        |
        |      ws message      |    answerSDP         |
        |     　answerSDP      |                      |
        |                      |                      |
        + setRemoteDescription(answerSDP)             |
        | Remote Description をセット                 |
        |                      |                      |
　   　 |                      |                      |
    ------------------ ICE candidate の交換 -----------------------
        |                      |                      |
　　　　+onicecandidate()の発火|                      |
        |  candidate の取得    |                      |
        +--------------------->+                      |
        |      ws message      +--------------------> | peerConnection に　ice candidate を追加する
        |  {type: "candidate", |   ws message         | peer.addIceCandidate(candidate)
        |    ice: candidate}   |   {type: "candidate",|
        |                      |   ice: candidate}    |　
      ==== 同様に browser2 から browser1 への ICE candidate の交換を行う ====
        |                      |                      |
     ========= ICE negotiation があれば 再び SDP をやりとり ================
        |                      |                      |
        + onaddstream()の発火  |                      + onaddstream()の発火
        | remoteStream をセット（browser2）           | remoteStream をセット(browser1)
    ------------------ Peer　Connection 確立 -----------------------
 　　   |                      |                      |　
```


### プロトコル

WS のメッセージはJSONフォーマットでやり取りします。
すべてのメッセージはプロパティに `type` を持ちます。
`type` は以下の5つです。

- register
- accept
- reject
- offer
- answer
- candidate
- close

#### type: register

クライアントが Ayame Server に room id, client id を登録するメッセージです。

```
{type: "register", "roomId" "<string>", "clientId": "<string>"}
```

これを受け取った Ayame Server はそのクライアントが指定した room に入室可能か検査して、可能であれば accept, 不可であれば reject を返却します。

#### type: accept

Ayame Server がregister に込められている情報を検査して入室可能であることをクライアントに知らせるメッセージです。

```
{type: "accept"}
```

将来的に、TURN などのメタデータもここで返却される予定です。
これを受け取ったらクライアントは offer のやり取りを開始します。

#### type: reject

Ayame Server がregister に込められている情報を検査して入室不可能であることをクライアントに知らせるメッセージです。

```
{type: "reject"}
```

これを受け取ったらクライアントは peerConnection, websocket を閉じて初期化します。

#### type: offer

offer SDP を送信するメッセージです。

```
{type: "offer", sdp: "v=0\r\no=- 4765067307885144980..."}
```

これを受け取ったクライアントはこのSDP をもとに peer connection に remote description をセットします。
また、このタイミングで local description を生成し、anwser SDP を送信します。

#### type: answer

answer SDP を送信するメッセージです。

```
{type: "answer", sdp: "v=0\r\no=- 4765067307885144980..."}
```

これを受け取ったクライアントはこれをもとに peer connection に remote description をセットします。

#### type: candidate

ice candidate を交換するメッセージです。

```
{type: "candidate", ice: {candidate: "...."}}
```

これを受け取ったクライアントは peer connection に ice candidate を追加します。

#### type: close

peer connection を切断したことを知らせるメッセージです。

```
{type: "close"}
```

これを受け取ったクライアントは peer connection を閉じて、リモート(受信側)の video element を破棄します。
