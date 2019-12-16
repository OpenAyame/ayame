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

type message struct {
	Type string `json:"type"`
}

type registerMessage struct {
	Type          string       `json:"type"`
	RoomID        string       `json:"roomId"`
	ClientID      string       `json:"clientId"`
	AuthnMetadata *interface{} `json:"authnMetadata,omitempty"`
	SignalingKey  *string      `json:"signalingKey,omitempty"`
	Key           *string      `json:"key,omitempty"`
}

type pingMessage struct {
	Type string `json:"type"`
}

type byeMessage struct {
	Type string `json:"type"`
}

func (c *Client) listen(cancel context.CancelFunc) {
	defer func() {
		cancel()
		c.hub.unregister <- &registerInfo{
			client: c,
			roomID: c.roomID,
		}
		c.conn.Close()
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logger.Warnf("failed to set read deadline, err=%v", err)
	}
	c.conn.SetCloseHandler(func(code int, text string) error {
		logger.Printf("Close code: %d, message: %s", code, text)
		logger.Printf("Client roomID: %s", c.roomID)
		if c.roomID == "" {
			msg := fmt.Sprintf("Client does not registered: %v", c)
			logger.Printf(msg)
			err := errors.New(msg)
			return err
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
	})

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
		case "":
			logger.Warnf("Invalid Signaling Type")
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

			var signalingKey *string
			if registerMessage.Key != nil {
				signalingKey = registerMessage.Key
			}
			if registerMessage.SignalingKey != nil {
				signalingKey = registerMessage.SignalingKey
			}
			if registerMessage.RoomID != "" {
				logger.Printf("Register: %v", message)
				c.hub.register <- &registerInfo{
					clientID:      registerMessage.ClientID,
					client:        c,
					roomID:        registerMessage.RoomID,
					signalingKey:  signalingKey,
					authnMetadata: registerMessage.AuthnMetadata,
				}
			}
		default:
			logger.Printf("Onmessage: %s", rawMessage)
			logger.Printf("Client roomID: %s", c.roomID)

			if c.roomID == "" {
				logger.Printf("Client does not registered: %v", c)
				break
			}
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Printf("error: %v", err)
				}
				break
			}
			broadcast := &Broadcast{
				client:   c,
				roomID:   c.roomID,
				messages: rawMessage,
			}
			c.hub.broadcast <- broadcast
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
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					logger.Warnf("failed to write close message, err=%v", err)
				}
				return
			}
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				logger.Warnf("Failed to set write deadline, err=%v", err)
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				logger.Warnf("Failed to write message, err=%v", err)
				return
			}
			w.Close()
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

func signalingHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Println(err)
		return
	}
	client := &Client{hub: hub, conn: c, send: make(chan []byte, 256)}
	logger.Printf("Websocket connected")
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
