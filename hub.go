package main

type Broadcast struct {
	client   *Client
	roomId   string
	messages []byte
}

type RegisterInfo struct {
	roomId   string
	clientId string
	client   *Client
	metadata *interface{}
	key      *string
}

type Room struct {
	clients map[*Client]bool
	roomId  string
}

type Hub struct {
	rooms map[string]*Room

	broadcast chan *Broadcast

	register chan *RegisterInfo

	unregister chan *RegisterInfo
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Broadcast),
		register:   make(chan *RegisterInfo),
		unregister: make(chan *RegisterInfo),
		rooms:      make(map[string]*Room),
	}
}

func (h *Hub) run() {
	for {
		select {
		case registerInfo := <-h.register:
			client := registerInfo.client
			clientId := registerInfo.clientId
			roomId := registerInfo.roomId
			client = client.Setup(roomId, clientId)
			room := h.rooms[roomId]
			if _, ok := h.rooms[roomId]; !ok {
				room = &Room{
					clients: make(map[*Client]bool),
					roomId:  roomId,
				}
				h.rooms[roomId] = room
			}
			ok := len(room.clients) < 2
			if !ok {
				msg := &RejectMessage{
					Type:   "reject",
					Reason: "TOO-MANY-USERS",
				}
				client.SendJSON(msg)
				client.conn.Close()
				break
			}
			// auth webhook を用いる場合
			if Options.AuthWebhookUrl != "" {
				resp, err := AuthWebhookRequest(registerInfo.key, roomId, registerInfo.metadata, client.host)
				if err != nil {
					msg := &RejectMessage{
						Type:   "reject",
						Reason: "AUTH-WEBHOOK-ERROR",
					}
					if resp != nil {
						msg.Reason = resp.Reason
					}
					client.SendJSON(msg)
					client.conn.Close()
					break
				}
				msg := &AcceptMetadataMessage{
					Type:       "accept",
					IceServers: resp.IceServers,
				}
				if resp.AuthzMetadata != nil {
					msg.Metadata = resp.AuthzMetadata
				}

				room.clients[client] = true
				client.SendJSON(msg)
			} else {
				msg := &AcceptMessage{
					Type: "accept",
				}
				room.clients[client] = true
				client.SendJSON(msg)
			}
		case registerInfo := <-h.unregister:
			roomId := registerInfo.roomId
			client := registerInfo.client
			if room, ok := h.rooms[roomId]; ok {
				if _, ok := room.clients[client]; ok {
					delete(room.clients, client)
					close(client.send)
				}
			}
		case broadcast := <-h.broadcast:
			if room, ok := h.rooms[broadcast.roomId]; ok {
				for client := range room.clients {
					if client.clientId != broadcast.client.clientId {
						select {
						case client.send <- broadcast.messages:
						default:
							close(client.send)
							delete(room.clients, client)
						}
					}
				}
			}
		}
	}
}
