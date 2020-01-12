package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	ID            string
	roomID        string
	authnMetadata *interface{}
	signalingKey  *string

	// WebSocket コネクション
	conn *websocket.Conn

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

func (c *client) SendJSON(v interface{}) error {
	if err := c.conn.WriteJSON(v); err != nil {
		c.errLog().Err(err).Interface("msg", v).Msg("FailedToSendMsg")
		return err
	}
	return nil
}

func (c *client) sendPingMessage() error {
	msg := &pingMessage{
		Type: "ping",
	}

	if err := c.SendJSON(msg); err != nil {
		return err
	}

	return nil
}

// reason の長さが不十分そうな場合は CloseMessage ではなく TextMessage を使用するように変更する
func (c *client) sendCloseMessage(code int, reason string) error {
	deadline := time.Now().Add(writeWait)
	closeMessage := websocket.FormatCloseMessage(code, reason)
	return c.conn.WriteControl(websocket.CloseMessage, closeMessage, deadline)
}

func (c *client) sendAcceptMessage(isExistClient bool, iceServers *[]iceServer, authzMetadata *interface{}) error {
	msg := &acceptMessage{
		Type:          "accept",
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

func (c *client) sendRejectMessage(reason string) error {
	msg := &rejectMessage{
		Type:   "reject",
		Reason: reason,
	}

	if err := c.SendJSON(msg); err != nil {
		return err
	}
	return nil
}

func (c *client) sendByeMessage() error {
	msg := &byeMessage{
		Type: "bye",
	}

	if err := c.SendJSON(msg); err != nil {
		return err
	}
	return nil
}

func (c *client) closeWs() {
	c.conn.Close()
	c.debugLog().Msg("CLOSED-WS-CONN")
}

func (c *client) register() int {
	resultChannel := make(chan int)
	registerChannel <- &register{
		client:        c,
		resultChannel: resultChannel,
	}
	// ここはブロックする candidate とかを並列で来てるかもしれないが知らん
	result, _ := <-resultChannel
	// もう server で触ることはないのでここで閉じる
	// これ server で閉じてもいいのかも
	close(resultChannel)
	return result
}

func (c *client) unregister() {
	if c.registered {
		unregisterChannel <- &unregister{
			client: c,
		}
	}
}

func (c *client) forward(msg []byte) {
	// グローバルにあるチャンネルに対して投げ込む
	forwardChannel <- forward{
		client:     c,
		rawMessage: msg,
	}
}

func (c *client) main(cancel context.CancelFunc, messageChannel chan []byte) {
	pongTimeoutTimer := time.NewTimer(pongTimeout * time.Second)
	pingTimer := time.NewTimer(pingInterval * time.Second)

	defer timerStop(pongTimeoutTimer)
	defer timerStop(pingTimer)
	defer c.unregister()
	defer cancel()

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
			// チャンネルが閉じられたら loop を抜ける
			if !ok {
				break loop
			}
			if err := c.handleWsMessage(rawMessage, pongTimeoutTimer); err != nil {
				// エラーになったら抜ける
				break loop
			}
		case forward, ok := <-c.forwardChannel:
			if !ok {
				// server 側でforwardChannel を閉じた
				c.debugLog().Msg("UNREGISTERED")
				if err := c.sendByeMessage(); err != nil {
				}
				c.debugLog().Msg("SENT-BYE-MESSAGE")
				break loop
			}
			c.conn.WriteMessage(websocket.TextMessage, forward.rawMessage)
		}
	}
	// 終了するので Websocket 終了のお知らせを送る
	if err := c.sendCloseMessage(websocket.CloseNormalClosure, ""); err != nil {
	}
	c.debugLog().Msg("EXIT-MAIN")
}

func (c *client) wsRecv(ctx context.Context, messageChannel chan []byte) {
loop:
	for {
		readDeadline := time.Now().Add(time.Duration(readTimeout) * time.Second)
		if err := c.conn.SetReadDeadline(readDeadline); err != nil {
			c.errLog().Err(err).Msg("FailedSetReadDeadLine")
			break loop
		}
		// 定期的に main が終了していないかチェックする
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			c.debugLog().Err(err).Msg("WS-READ-MESSAGE-ERROR")
			break loop
		}

		select {
		case <-ctx.Done():
			// メインが死んでたらループ抜けて終了
			c.debugLog().Msg("EXITED-MAIN")
			break loop
		default:
			messageChannel <- rawMessage
		}
	}
	close(messageChannel)
	c.debugLog().Msg("CLOSED-MESSAGE-CHANNEL")
	c.closeWs()
	c.debugLog().Msg("CLOSED-WS")
	c.debugLog().Msg("EXIT-WS-RECV")
}

// メッセージ系のエラーログはここですべて取る
func (c *client) handleWsMessage(rawMessage []byte, pongTimeoutTimer *time.Timer) error {
	message := &message{}
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		c.errLog().Err(err).Msg("InvalidJSON")
		return errInvalidJSON
	}

	// 受信したメッセージで message type がパースできたものをログとして保存する
	c.signalingLog(*message, rawMessage)

	switch message.Type {
	case "pong":
		timerStop(pongTimeoutTimer)
		pongTimeoutTimer.Reset(pongTimeout * time.Second)
	case "register":
		registerMessage := &registerMessage{}
		if err := json.Unmarshal(rawMessage, &registerMessage); err != nil {
			c.errLog().Err(err).Msg("InvalidJSON")
			return errInvalidJSON
		}

		if registerMessage.RoomID == "" {
			// XXX(nakai): どんな JSON だったのか見たくなるはず
			c.errLog().Msg("MissingRoomID")
			return errMissingRoomID
		}
		c.roomID = registerMessage.RoomID

		if registerMessage.ClientID == "" {
			// XXX(nakai): どんな JSON だったのか見たくなるはず
			c.errLog().Msg("MissingClientID")
			return errMissingClientID
		}
		c.ID = registerMessage.ClientID

		// 下位互換性
		if registerMessage.Key != nil {
			c.signalingKey = registerMessage.Key
		}

		if registerMessage.SignalingKey != nil {
			c.signalingKey = registerMessage.SignalingKey
		}

		c.authnMetadata = registerMessage.AuthnMetadata

		resp, err := c.authnWebhook()
		if err != nil {
			c.errLog().Err(err).Caller().Msg("AuthnWebhookError")
			return err
		}

		// 認証サーバの戻り値がおかしい場合は全部 Error にする
		if resp.Allowed == nil {
			c.errLog().Caller().Msg("AuthnWebhookResponseError")
			return errAuthnWebhookResponse
		}

		if !*resp.Allowed {
			if resp.Reason == nil {
				c.errLog().Caller().Msg("AuthnWebhookResponseError")
				if err := c.sendRejectMessage("InternalServerError"); err != nil {
					return err
				}
				return errAuthnWebhookResponse
			}
			if err := c.sendRejectMessage(*resp.Reason); err != nil {
				c.errLog().Err(err).Msg("FailedSendRejectMessage")
				return err
			}
			return errAuthnWebhookReject
		}

		// 戻り値は手抜き
		switch c.register() {
		case one:
			// room がまだなかった、accept を返す
			c.debugLog().Msg("REGISTERED-ONE")
			c.registered = true
			if err := c.sendAcceptMessage(false, resp.IceServers, resp.AuthzMetadata); err != nil {
				c.errLog().Err(err).Msg("FailedSendAcceptMessage")
				return err
			}
		case two:
			// room がすでにあって、一人いた、二人目
			c.debugLog().Msg("REGISTERED-TWO")
			c.registered = true
			if err := c.sendAcceptMessage(true, resp.IceServers, resp.AuthzMetadata); err != nil {
				c.errLog().Err(err).Msg("FailedSendAcceptMessage")
				return err
			}
		case full:
			// room が満杯だった
			c.errLog().Msg("RoomFull")
			if err := c.sendRejectMessage("full"); err != nil {
				c.errLog().Err(err).Msg("FailedSendRejectMessage")
				return err
			}
			return errRoomFull
		case dup:
			// clientID が重複してた
			c.errLog().Msg("DuplicateClientID")
			if err := c.sendRejectMessage("dup"); err != nil {
				c.errLog().Err(err).Msg("FailedSendRejectMessage")
				return err
			}
			return errDuplicateClientID
		}
	case "offer", "answer", "candidate":
		if c.ID == "" {
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
