package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
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
	Type     string `json:"type"`
	RoomId   string `json:"room_id"`
	ClientId string `json:"client_id"`
}

func (c *Client) listen() {
	defer func() {
		c.hub.unregister <- &RegisterInfo{
			client: c,
			roomId: c.roomId,
		}
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		msg := &Message{}
		json.Unmarshal(message, &msg)
		if msg.Type == "register" && msg.RoomId != "" {
			logger.Printf("register: %v", msg)
			c.Lock()
			defer c.Unlock()
			c.roomId = msg.RoomId
			c.clientId = msg.ClientId
			c.hub.register <- &RegisterInfo{
				client: c,
				roomId: msg.RoomId,
			}
		} else {
			logger.Printf("onmessage: %v", string(message))
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

func (c *Client) broadcast() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
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
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
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
	client := &Client{hub: hub, conn: c, send: make(chan []byte, 256)}
	logger.Printf("[WS] connected")
	go client.listen()
	go client.broadcast()
}
