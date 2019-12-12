package main

import (
	"sync"
)

type Broadcast struct {
	client   *Client
	roomID   string
	messages []byte
}

type registerInfo struct {
	roomID        string
	clientID      string
	client        *Client
	authnMetadata *interface{}
	signalingKey  *string
}

type Room struct {
	clients map[*Client]bool
	roomID  string
	sync.Mutex
}

func (r *Room) newClient(client *Client) {
	r.Lock()
	defer r.Unlock()
	r.clients[client] = true
}

func (r *Room) deleteClient(client *Client) {
	r.Lock()
	defer r.Unlock()
	close(client.send)
	delete(r.clients, client)
}

type Hub struct {
	rooms map[string]*Room

	broadcast chan *Broadcast

	register chan *registerInfo

	unregister chan *registerInfo
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Broadcast),
		register:   make(chan *registerInfo),
		unregister: make(chan *registerInfo),
		rooms:      make(map[string]*Room),
	}
}

func (h *Hub) run() {
	for {
		select {
		case registerInfo := <-h.register:
			client := registerInfo.client
			clientID := registerInfo.clientID
			roomID := registerInfo.roomID
			if len(roomID) == 0 || len(clientID) == 0 {
				msg := &rejectMessage{
					Type:   "reject",
					Reason: "INVALID-ROOM-ID-OR-CLIENT-ID",
				}
				err := client.SendJSON(msg)
				if err != nil {
					logger.Warnf("failed to send msg=%v", msg)
				}
				client.conn.Close()
				break
			}
			client = client.Setup(roomID, clientID)
			room := h.rooms[roomID]
			if _, ok := h.rooms[roomID]; !ok {
				room = &Room{
					clients: make(map[*Client]bool),
					roomID:  roomID,
				}
				h.rooms[roomID] = room
			}
			ok := len(room.clients) < 2
			if !ok {
				msg := &rejectMessage{
					Type:   "reject",
					Reason: "TOO-MANY-USERS",
				}
				err := client.SendJSON(msg)
				if err != nil {
					logger.Warnf("failed to send msg=%v", msg)
				}
				client.conn.Close()
				break
			}
			if options.AuthWebhookURL != "" {
				resp, err := authWebhookRequest(roomID, clientID, registerInfo.authnMetadata, registerInfo.signalingKey)
				// インターナルエラー
				if err != nil {
					msg := &rejectMessage{
						Type:   "reject",
						Reason: "AUTH-WEBHOOK-INTERNAL-ERROR",
					}
					err = client.SendJSON(msg)
					if err != nil {
						logger.Warnf("Failed to send msg=%v", msg)
					}
					client.conn.Close()
					break
				}

				// 認証失敗
				if !resp.Allowed {
					msg := &rejectMessage{
						Type:   "reject",
						Reason: resp.Reason,
					}
					err = client.SendJSON(msg)
					if err != nil {
						logger.Warnf("Failed to send msg=%v", msg)
					}
					client.conn.Close()
					break
				}

				// 認証成功
				isExistUser := len(room.clients) > 0
				msg := &acceptMessage{
					Type:        "accept",
					IsExistUser: isExistUser,
					IceServers:  resp.IceServers,
					// TODO(nakai): authz を個々に入れる
				}
				room.newClient(client)
				err = client.SendJSON(msg)
				if err != nil {
					logger.Warnf("Failed to send msg=%v", msg)
				}
				client.conn.Close()
			}
		case registerInfo := <-h.unregister:
			roomID := registerInfo.roomID
			client := registerInfo.client
			if room, ok := h.rooms[roomID]; ok {
				if _, ok := room.clients[client]; ok {
					room.deleteClient(client)
				}
			}
		case broadcast := <-h.broadcast:
			if room, ok := h.rooms[broadcast.roomID]; ok {
				for client := range room.clients {
					if client.clientID != broadcast.client.clientID {
						select {
						case client.send <- broadcast.messages:
						default:
							room.deleteClient(client)
						}
					}
				}
			}
		}
	}
}
