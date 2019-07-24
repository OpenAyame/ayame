package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/ryanuber/go-glob"
	"net/http"
	"time"
)

const (
	writeWait = 10 * time.Second

	pongWait = 10 * time.Second

	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// クロスオリジンを一旦許可する
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type     string      `json:"type"`
	RoomId   string      `json:"roomId"`
	ClientId string      `json:"clientId"`
	Metadata interface{} `json:"authnMetadata"`
	Key      string      `json:"key"`
}

type PingMessage struct {
	Type string `json:"type"`
}

func (c *Client) listen(cancel context.CancelFunc) {
	defer func() {
		cancel()
		c.hub.unregister <- &RegisterInfo{
			client: c,
			roomId: c.roomId,
		}
		c.conn.Close()
	}()

	upgrader.CheckOrigin = func(r *http.Request) bool {
		if Options.AllowOrigin == "" {
			return true
		}
		origin := r.Header.Get("Origin")
		// origin を trim
		host, err := TrimOriginToHost(origin)
		if err != nil {
			logger.Warn("Invalid Origin Header, header=", origin)
		}
		// config.yaml で指定した Allow Origin と一致するかで検査する
		logger.Infof("[WS] Request Origin=%s, AllowOrigin=%s", origin, Options.AllowOrigin)
		if &Options.AllowOrigin == host {
			return true
		}
		if glob.Glob(Options.AllowOrigin, *host) {
			return true
		}
		return false
	}
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		msg := &Message{}
		json.Unmarshal(message, &msg)
		if msg.Type == "" {
			logger.Warnf("Invalid Signaling Type")
			break
		}
		if msg.Type == "pong" {
			logger.Printf("recv ping over WS")
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
		} else {
			if msg.Type == "register" && msg.RoomId != "" {
				logger.Printf("register: %v", msg)
				c.hub.register <- &RegisterInfo{
					clientId: msg.ClientId,
					client:   c,
					roomId:   msg.RoomId,
					key:      msg.Key,
					metadata: msg.Metadata,
				}
			} else {
				logger.Printf("onmessage: %s", message)
				logger.Printf("client roomId: %s", c.roomId)
				if c.roomId == "" {
					logger.Printf("client does not registered: %v", c)
					return
				}
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						logger.Printf("error: %v", err)
					}
					break
				}
				broadcast := &Broadcast{
					client:   c,
					roomId:   c.roomId,
					messages: message,
				}
				c.hub.broadcast <- broadcast
			}
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
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// over Ws で ping-pong を設定している場合
			if Options.OverWsPingPong {
				logger.Info("send ping over WS")
				pingMsg := &PingMessage{Type: "ping"}
				if err := c.SendJSON(pingMsg); err != nil {
					return
				}
			} else {
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}

func wsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Println(err)
		return
	}
	origin := r.Header.Get("Origin")
	host, err := TrimOriginToHost(origin)
	if err != nil {
		logger.Println(err)
		return
	}
	client := &Client{hub: hub, conn: c, host: *host, send: make(chan []byte, 256)}
	logger.Printf("[WS] connected")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go client.listen(cancel)
	go client.broadcast(ctx)
}
