package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	roomId   string
	clientId string
	send     chan []byte
	sync.Mutex
}

// json を送る
func (c *Client) SendJSON(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.conn.WriteJSON(v)
}

func (c *Client) Setup(roomId string, clientId string) *Client {
	c.Lock()
	defer c.Unlock()
	c.roomId = roomId
	c.clientId = clientId
	return c
}
