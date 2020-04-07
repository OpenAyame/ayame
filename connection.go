package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
)

type connection struct {
	ID            string
	roomID        string
	clientID      string
	authnMetadata *interface{}
	signalingKey  *string

	// クライアント情報
	ayameClient *string
	environment *string
	libwebrtc   *string

	authzMetadata *interface{}

	// WebSocket コネクション
	wsConn *websocket.Conn

	// レジスターされているかどうか
	registered bool

	// 転送用のチャネル
	forwardChannel chan forward
}

const (
	// socket の待ち受け時間
	readTimeout = 90
	// pong が送られてこないためタイムアウトにするまでの時間
	pongTimeout = 60
	// ping 送信の時間間隔
	pingInterval = 5
)

func (c *connection) SendJSON(v interface{}) error {
	if err := c.wsConn.WriteJSON(v); err != nil {
		c.errLog().Err(err).Interface("msg", v).Msg("FailedToSendMsg")
		return err
	}
	return nil
}

func (c *connection) sendPingMessage() error {
	msg := &pingMessage{
		Type: "ping",
	}

	if err := c.SendJSON(msg); err != nil {
		return err
	}

	return nil
}

// reason の長さが不十分そうな場合は CloseMessage ではなく TextMessage を使用するように変更する
func (c *connection) sendCloseMessage(code int, reason string) error {
	deadline := time.Now().Add(writeWait)
	closeMessage := websocket.FormatCloseMessage(code, reason)
	return c.wsConn.WriteControl(websocket.CloseMessage, closeMessage, deadline)
}

func (c *connection) sendAcceptMessage(isExistClient bool, iceServers *[]iceServer, authzMetadata *interface{}) error {
	msg := &acceptMessage{
		Type:          "accept",
		ConnectionID:  c.ID,
		IsExistClient: isExistClient,
		// 下位互換性
		IsExistUser:   isExistClient,
		AuthzMetadata: authzMetadata,
		IceServers:    iceServers,
	}

	if err := c.SendJSON(msg); err != nil {
		return err
	}
	return nil
}

func (c *connection) sendRejectMessage(reason string) error {
	msg := &rejectMessage{
		Type:   "reject",
		Reason: reason,
	}

	if err := c.SendJSON(msg); err != nil {
		return err
	}
	return nil
}

func (c *connection) sendByeMessage() error {
	msg := &byeMessage{
		Type: "bye",
	}

	if err := c.SendJSON(msg); err != nil {
		return err
	}
	return nil
}

func (c *connection) closeWs() {
	c.wsConn.Close()
	c.debugLog().Msg("CLOSED-WS")
}

func (c *connection) register() int {
	resultChannel := make(chan int)
	registerChannel <- &register{
		connection:    c,
		resultChannel: resultChannel,
	}
	// ここはブロックする candidate とかを並列で来てるかもしれないが知らん
	result := <-resultChannel
	// もう server で触ることはないのでここで閉じる
	close(resultChannel)
	return result
}

func (c *connection) unregister() {
	if c.registered {
		unregisterChannel <- &unregister{
			connection: c,
		}
	}
}

func (c *connection) forward(msg []byte) {
	// グローバルにあるチャンネルに対して投げ込む
	forwardChannel <- forward{
		connection: c,
		rawMessage: msg,
	}
}

func (c *connection) main(cancel context.CancelFunc, messageChannel chan []byte) {
	pongTimeoutTimer := time.NewTimer(pongTimeout * time.Second)
	pingTimer := time.NewTimer(pingInterval * time.Second)

	defer func() {
		timerStop(pongTimeoutTimer)
		timerStop(pingTimer)
		// キャンセルを呼ぶ
		cancel()
		c.debugLog().Msg("CANCEL")
		// アンレジはここでやる
		c.unregister()
		c.debugLog().Msg("UNREGISTER")
		c.debugLog().Msg("EXIT-MAIN")
	}()

loop:
	for {
		select {
		case <-pingTimer.C:
			if err := c.sendPingMessage(); err != nil {
				break loop
			}
			pingTimer.Reset(pingInterval * time.Second)
		case <-pongTimeoutTimer.C:
			// タイマーが発火してしまったので切断する
			c.errLog().Msg("PongTimeout")
			break loop
		case rawMessage, ok := <-messageChannel:
			// message チャンネルが閉じられた、main 終了待ち
			if !ok {
				c.debugLog().Msg("CLOSED-MESSAGE-CHANNEL")
				// メッセージチャネルが閉じてるので return でもう抜けてしまう
				return
			}
			if err := c.handleWsMessage(rawMessage, pongTimeoutTimer); err != nil {
				// ここのエラーのログはすでに handleWsMessage でとってあるので不要
				// エラーになったら抜ける
				break loop
			}
		case forward, ok := <-c.forwardChannel:
			if !ok {
				// server 側で forwardChannel を閉じた
				c.debugLog().Msg("UNREGISTERED")
				if err := c.sendByeMessage(); err != nil {
					c.errLog().Err(err).Msg("FailedSendByeMessage")
					// 送れなかったら閉じるメッセージも送れないので return
					return
				}
				c.debugLog().Msg("SENT-BYE-MESSAGE")
				break loop
			}
			if err := c.wsConn.WriteMessage(websocket.TextMessage, forward.rawMessage); err != nil {
				c.errLog().Err(err).Msg("FailedWriteMessage")
				// 送れなかったら閉じるメッセージも送れないので return
				return
			}
		}
	}

	// こちらの都合で終了するので Websocket 終了のお知らせを送る
	if err := c.sendCloseMessage(websocket.CloseNormalClosure, ""); err != nil {
		c.debugLog().Err(err).Msg("FAILED-SEND-CLOSE-MESSAGE")
		// 送れなかったら return する
		return
	}
	c.debugLog().Msg("SENT-CLOSE-MESSAGE")
}

func (c *connection) wsRecv(ctx context.Context, messageChannel chan []byte) {
loop:
	for {
		readDeadline := time.Now().Add(time.Duration(readTimeout) * time.Second)
		if err := c.wsConn.SetReadDeadline(readDeadline); err != nil {
			c.errLog().Err(err).Msg("FailedSetReadDeadLine")
			break loop
		}
		_, rawMessage, err := c.wsConn.ReadMessage()
		if err != nil {
			// ここに来るのはほぼ WebSocket が切断されたとき
			c.debugLog().Err(err).Msg("WS-READ-MESSAGE-ERROR")
			break loop
		}
		messageChannel <- rawMessage
	}
	close(messageChannel)
	c.debugLog().Msg("CLOSE-MESSAGE-CHANNEL")
	// メインが死ぬまで待つ
	<-ctx.Done()
	c.debugLog().Msg("EXITED-MAIN")
	c.closeWs()
	c.debugLog().Msg("EXIT-WS-RECV")

	if err := c.disconnectWebhook(); err != nil {
		c.errLog().Err(err).Caller().Msg("DisconnectWebhookError")
		return
	}
}

// メッセージ系のエラーログはここですべて取る
func (c *connection) handleWsMessage(rawMessage []byte, pongTimeoutTimer *time.Timer) error {
	message := &message{}
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		c.errLog().Err(err).Bytes("rawMessage", rawMessage).Msg("InvalidJSON")
		return errInvalidJSON
	}

	if message == nil {
		c.errLog().Bytes("rawMessage", rawMessage).Msg("UnexpectedJSON")
		return errUnexpectedJSON
	}

	// 受信したメッセージで message type がパースできたものをログとして保存する
	c.signalingLog(*message, rawMessage)

	switch message.Type {
	case "pong":
		timerStop(pongTimeoutTimer)
		pongTimeoutTimer.Reset(pongTimeout * time.Second)
	case "register":
		// すでに登録されているのにもう一度登録しに来た
		if c.registered {
			c.errLog().Bytes("rawMessage", rawMessage).Msg("InternalServer")
			return errInternalServer
		}

		c.ID = getULID()

		registerMessage := &registerMessage{}
		if err := json.Unmarshal(rawMessage, &registerMessage); err != nil {
			c.errLog().Err(err).Bytes("rawMessage", rawMessage).Msg("InvalidRegisterMessageJSON")
			return errInvalidJSON
		}

		if registerMessage.RoomID == "" {
			c.errLog().Bytes("rawMessage", rawMessage).Msg("MissingRoomID")
			return errMissingRoomID
		}
		c.roomID = registerMessage.RoomID

		c.clientID = registerMessage.ClientID
		if registerMessage.ClientID == "" {
			c.clientID = c.ID
		}

		// 下位互換性
		if registerMessage.Key != nil {
			c.signalingKey = registerMessage.Key
		}

		if registerMessage.SignalingKey != nil {
			c.signalingKey = registerMessage.SignalingKey
		}

		c.authnMetadata = registerMessage.AuthnMetadata

		// クライアント情報の登録
		c.ayameClient = registerMessage.AyameClient
		c.environment = registerMessage.Environment
		c.libwebrtc = registerMessage.Libwebrtc

		// Webhook 系のエラーログは Caller をつける
		resp, err := c.authnWebhook()
		if err != nil {
			c.errLog().Err(err).Caller().Msg("AuthnWebhookError")
			if err := c.sendRejectMessage("InternalServerError"); err != nil {
				c.errLog().Err(err).Caller().Msg("FailedSendRejectMessage")
				return err
			}
			return err
		}

		// 認証サーバの戻り値がおかしい場合は全部 Error にする
		if resp.Allowed == nil {
			c.errLog().Caller().Msg("AuthnWebhookResponseError")
			if err := c.sendRejectMessage("InternalServerError"); err != nil {
				c.errLog().Err(err).Caller().Msg("FailedSendRejectMessage")
				return err
			}
			return errAuthnWebhookResponse
		}

		if !*resp.Allowed {
			if resp.Reason == nil {
				c.errLog().Caller().Msg("AuthnWebhookResponseError")
				if err := c.sendRejectMessage("InternalServerError"); err != nil {
					c.errLog().Err(err).Caller().Msg("FailedSendRejectMessage")
					return err
				}
				return errAuthnWebhookResponse
			}
			if err := c.sendRejectMessage(*resp.Reason); err != nil {
				c.errLog().Err(err).Caller().Msg("FailedSendRejectMessage")
				return err
			}
			return errAuthnWebhookReject
		}

		c.authzMetadata = resp.AuthzMetadata

		// 戻り値は手抜き
		switch c.register() {
		case one:
			c.registered = true
			// room がまだなかった、accept を返す
			c.debugLog().Msg("REGISTERED-ONE")
			if err := c.sendAcceptMessage(false, resp.IceServers, resp.AuthzMetadata); err != nil {
				c.errLog().Err(err).Msg("FailedSendAcceptMessage")
				return err
			}
		case two:
			c.registered = true
			// room がすでにあって、一人いた、二人目
			c.debugLog().Msg("REGISTERED-TWO")
			if err := c.sendAcceptMessage(true, resp.IceServers, resp.AuthzMetadata); err != nil {
				c.errLog().Err(err).Msg("FailedSendAcceptMessage")
				return err
			}
		case full:
			// room が満杯だった
			c.errLog().Msg("RoomFilled")
			if err := c.sendRejectMessage("full"); err != nil {
				c.errLog().Err(err).Msg("FailedSendRejectMessage")
				return err
			}
			return errRoomFull
		}
	case "offer", "answer", "candidate":
		// register が完了していない
		if !c.registered {
			c.errLog().Msg("RegistrationIncomplete")
			return errRegistrationIncomplete
		}
		c.forward(rawMessage)
	default:
		c.errLog().Msg("InvalidMessageType")
		return errInvalidMessageType
	}
	return nil
}

func timerStop(timer *time.Timer) {
	// タイマー終了からのリセットへは以下参考にした
	// https://www.kaoriya.net/blog/2019/12/19/
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}

func getULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
