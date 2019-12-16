package main

import (
	"sync"

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

func (c *Client) Setup(roomID string, clientID string) {
	c.Lock()
	defer c.Unlock()
	c.roomID = roomID
	c.clientID = clientID
}

func (c *Client) sendRejectMessage(reason string) error {
	msg := &rejectMessage{
		Type:   "reject",
		Reason: "AUTH-WEBHOOK-INTERNAL-ERROR",
	}

	if err := c.SendJSON(msg); err != nil {
		logger.Warnf("Failed to send msg=%v", msg)
		return err
	}
	c.conn.Close()
	return nil
}
