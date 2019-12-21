package main

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	roomID   string
	clientID string
	send     chan []byte
	sync.Mutex
}

// json を送る
func (c *Client) SendJSON(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.conn.WriteJSON(v)
}

func (c *Client) sendRejectMessage(reason string) error {
	msg := &rejectMessage{
		Type:   "reject",
		Reason: reason,
	}

	if err := c.SendJSON(msg); err != nil {
		logger.Warnf("Failed to send msg=%v", msg)
		return err
	}
	return nil
}

// TODO(yoshida): reason の長さが不十分そうな場合は CloseMessage ではなく TextMessage を使用するように変更する
func (c *Client) sendCloseMessage(code int, reason string) error {
	c.Lock()
	defer c.Unlock()

	deadline := time.Now().Add(writeWait)
	closeMessage := websocket.FormatCloseMessage(code, reason)
	return c.conn.WriteControl(websocket.CloseMessage, closeMessage, deadline)
}

func (c *Client) listen(cancel context.CancelFunc) {
	defer func() {
		cancel()
		c.hub.unregister <- &registerInfo{
			client: c,
			roomID: c.roomID,
		}
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logger.Warnf("failed to set read deadline, err=%v", err)
	}

	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			logger.Warnf("Error while read message, err=%v", err)
			break
		}
		message := &message{}
		if err := json.Unmarshal(rawMessage, &message); err != nil {
			logger.Warnf("Invalid JSON, err=%v", err)
			break
		}

		switch message.Type {
		case "pong":
			logger.Printf("Recv ping over WS")
			if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
				logger.Warnf("Failed to set read deadline, err=%v", err)
			}
		case "register":
			registerMessage := &registerMessage{}
			if err := json.Unmarshal(rawMessage, &registerMessage); err != nil {
				logger.Warnf("Invalid JSON, err=%v", err)
				break
			}

			if registerMessage.RoomID == "" {
				reason := "missing roomId"
				logger.Error(reason)
				if err := c.sendRejectMessage(reason); err != nil {
					logger.Error(err)
				}
				c.conn.Close()
				break
			}
			c.roomID = registerMessage.RoomID

			if registerMessage.ClientID == "" {
				reason := "missing clientId"
				logger.Error(reason)
				if err := c.sendRejectMessage(reason); err != nil {
					logger.Error(err)
				}
				c.conn.Close()
				break
			}
			c.clientID = registerMessage.ClientID

			var signalingKey *string
			if registerMessage.Key != nil {
				signalingKey = registerMessage.Key
			}
			if registerMessage.SignalingKey != nil {
				signalingKey = registerMessage.SignalingKey
			}

			logger.Printf("Register: %v", message)
			c.hub.register <- &registerInfo{
				client:        c,
				roomID:        registerMessage.RoomID,
				signalingKey:  signalingKey,
				authnMetadata: registerMessage.AuthnMetadata,
			}
		case "offer", "answer", "candidate":
			logger.Printf("Onmessage: %s", rawMessage)
			logger.Printf("Client roomID: %s", c.roomID)

			if c.roomID == "" {
				logger.Printf("Client does not registered: %v", c)
				break
			}
			broadcast := &Broadcast{
				client:   c,
				roomID:   c.roomID,
				messages: rawMessage,
			}
			c.hub.broadcast <- broadcast
		default:
			logger.Warnf("Invalid Signaling Type")
		}
	}
}

func (c *Client) broadcast(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case <-ctx.Done():
			// channel がすでに close していた場合
			// ループを抜ける
			if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
				logger.Warnf("failed to write close message, err=%v", err)
			}
			return
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				logger.Warnf("Failed to set write deadline, err=%v", err)
				return
			}
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					logger.Warnf("failed to write close message, err=%v", err)
				}
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.Warnf("Failed to write message, err=%v", err)
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				logger.Warnf("Failed to set write deadline, err=%v", err)
			}
			logger.Info("Send ping over WS")
			pingMsg := &pingMessage{Type: "ping"}
			if err := c.SendJSON(pingMsg); err != nil {
				return
			}
		}
	}
}
