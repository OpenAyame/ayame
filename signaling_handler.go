package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ryanuber/go-glob"
)

const (
	writeWait = 10 * time.Second

	pongWait = 10 * time.Second

	pingPeriod = (pongWait * 9) / 10
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 4,
		WriteBufferSize: 1024 * 4,
		CheckOrigin:     checkOrigin,
	}
)



func signalingHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	logger.Printf("Websocket connected")
	client.conn.SetCloseHandler(client.closeHandler)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go client.listen(cancel)
	go client.broadcast(ctx)
}

func checkOrigin(r *http.Request) bool {
	if options.AllowOrigin == "" {
		return true
	}
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	// origin を trim
	host, err := trimOriginToHost(origin)
	if err != nil {
		logger.Warn("Invalid Origin Header, header=", origin)
		return false
	}
	// config.yaml で指定した Allow Origin と一致するかで検査する
	logger.Infof("[WS] Request Origin=%s, AllowOrigin=%s", origin, options.AllowOrigin)
	if options.AllowOrigin == host {
		return true
	}
	if glob.Glob(options.AllowOrigin, host) {
		return true
	}
	return false
}

func (c *Client) closeHandler(code int, text string) error {
	logger.Printf("Close code: %d, message: %s", code, text)
	logger.Printf("Client roomID: %s", c.roomID)
	if c.roomID == "" {
		msg := fmt.Sprintf("Client does not registered: %v", c)
		logger.Printf(msg)
		return errors.New(msg)
	}
	byeMessage := &byeMessage{Type: "bye"}
	message, err := json.Marshal(byeMessage)
	if err != nil {
		logger.Printf("error: %v", err)
		return err
	}
	broadcast := &Broadcast{
		client:   c,
		roomID:   c.roomID,
		messages: message,
	}
	c.hub.broadcast <- broadcast
	return nil
}
