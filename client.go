package main

import (
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
