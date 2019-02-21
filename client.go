package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	uuid   string
	roomId string
	send   chan []byte
	sync.Mutex
}
