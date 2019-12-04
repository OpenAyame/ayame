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

func (c *Client) Setup(roomID string, clientID string) *Client {
	c.Lock()
	defer c.Unlock()
	c.roomID = roomID
	c.clientID = clientID
	return c
}
