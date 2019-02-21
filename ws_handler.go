package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

const (
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
	Type   string `json:"type"`
	RoomId string `json:"room_id"`
}

func (c *Client) listen() {
	defer func() {
		c.hub.unregister <- &RegisterInfo{
			client: c,
			roomId: c.roomId,
		}
		c.conn.Close()
	}()

	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		msg := &Message{}
		json.Unmarshal(message, &msg)
		if msg.Type == "register" && msg.RoomId != "" {
			log.Printf("register: %v", msg)
			c.Lock()
			defer c.Unlock()
			c.roomId = msg.RoomId
			c.hub.register <- &RegisterInfo{
				client: c,
				roomId: msg.RoomId,
			}
		} else {
			log.Printf("onmessage: %v", string(message))
			log.Printf("client roomId: %s", c.roomId)
			if c.roomId == "" {
				log.Printf("client does not registered: %v", c)
				return
			}
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
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
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
		}
	}
}

func wsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	uuid := uuid.NewV4().String()
	client := &Client{hub: hub, conn: c, send: make(chan []byte, 256), uuid: uuid}
	log.Printf("[WS] connected")
	go client.listen()
	go client.broadcast()
}
