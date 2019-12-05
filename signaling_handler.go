package main

import (
	"context"
	"encoding/json"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 4,
	WriteBufferSize: 1024 * 4,
	CheckOrigin: func(r *http.Request) bool {
		if Options.AllowOrigin == "" {
			return true
		}
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
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
	},
}

type Message struct {
	Type     string       `json:"type"`
	RoomID   string       `json:"roomId"`
	ClientID string       `json:"clientId"`
	Metadata *interface{} `json:"authnMetadata,omitempty"`
	Key      *string      `json:"key,omitempty"`
}

type PingMessage struct {
	Type string `json:"type"`
}

func (c *Client) listen(cancel context.CancelFunc) {
	defer func() {
		cancel()
		c.hub.unregister <- &RegisterInfo{
			client: c,
			roomID: c.roomID,
		}
		c.conn.Close()
	}()

	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		logger.Warnf("failed to set read deadline, err=%v", err)
	}
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			logger.Warnf("failed to set read deadline, err=%v", err)
		}
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			logger.Warnf("Error while read message, err=%v", err)
			break
		}
		msg := &Message{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			logger.Warnf("Invalid JSON, err=%v", err)
			break
		}
		if msg.Type == "" {
			logger.Warnf("Invalid Signaling Type")
			break
		}
		if msg.Type == "pong" {
			logger.Printf("recv ping over WS")
			err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				logger.Warnf("failed to set read deadline, err=%v", err)
			}
		} else {
			if msg.Type == "register" && msg.RoomID != "" {
				logger.Printf("register: %v", msg)
				c.hub.register <- &RegisterInfo{
					clientID: msg.ClientID,
					client:   c,
					roomID:   msg.RoomID,
					key:      msg.Key,
					metadata: msg.Metadata,
				}
			} else {
				logger.Printf("onmessage: %s", message)
				logger.Printf("client roomID: %s", c.roomID)
				if c.roomID == "" {
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
					roomID:   c.roomID,
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
			err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			if err != nil {
				logger.Warnf("failed to write close message, err=%v", err)
			}
			return
		case message, ok := <-c.send:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logger.Warnf("failed to write close message, err=%v", err)
				}
				return
			}
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.Warnf("failed to set write deadline, err=%v", err)
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				logger.Warnf("failed to write message, err=%v", err)
				return
			}
			w.Close()
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.Warnf("failed to set write deadline, err=%v", err)
			}
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

func signalingHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Println(err)
		return
	}
	client := &Client{hub: hub, conn: c, send: make(chan []byte, 256)}
	logger.Printf("[WS] connected")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go client.listen(cancel)
	go client.broadcast(ctx)
}
