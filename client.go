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
