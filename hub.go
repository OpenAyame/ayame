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
	roomID string
	// TODO(nakai): client がいるので clientID は削除
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

// TODO(nakai): registerClient
func (r *Room) newClient(client *Client) {
	r.Lock()
	defer r.Unlock()
	r.clients[client] = true
}

// TODO(nakai): unregisterClient
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
			// TODO(nakai): ROOM-ID と CLIENT-ID のエラーは分ける
			if len(roomID) == 0 || len(clientID) == 0 {
				reason := "INVALID-ROOM-ID-OR-CLIENT-ID"
				if err := client.sendRejectMessage(reason); err != nil {
					logger.Error(err)
				}
				break
			}
			client.Setup(roomID, clientID)
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
				reason := "TOO-MANY-USERS"
				if err := client.sendRejectMessage(reason); err != nil {
					logger.Error(err)
				}
				break
			}
			// 認証成功
			isExistUser := len(room.clients) > 0
			msg := &acceptMessage{
				Type:        "accept",
				IsExistUser: isExistUser,
			}
			if options.AuthWebhookURL != "" {
				resp, err := authWebhookRequest(roomID, clientID, registerInfo.authnMetadata, registerInfo.signalingKey)
				// インターナルエラー
				if err != nil {
					logger.Warnf("%s", err)
					reason := "AUTH-WEBHOOK-INTERNAL-ERROR"
					if err := client.sendRejectMessage(reason); err != nil {
						logger.Error(err)
					}
					break
				}

				// allowed が存在しない場合はエラー
				if resp.Allowed == nil {
					logger.Warn("missing allowed key")
					reason := "AUTH-WEBHOOK-INTERNAL-ERROR"
					if err := client.sendRejectMessage(reason); err != nil {
						logger.Error(err)
					}
					break
				}

				// 認証失敗
				if !*resp.Allowed {
					reason := "AUTH-WEBHOOK-INTERNAL-ERROR"
					if resp.Reason != nil {
						reason = *resp.Reason
					}
					if err := client.sendRejectMessage(reason); err != nil {
						logger.Error(err)
					}
					break
				}
				msg.IceServers = resp.IceServers
			}
			room.newClient(client)

			if err := client.SendJSON(msg); err != nil {
				logger.Warnf("Failed to send msg=%v", msg)
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
