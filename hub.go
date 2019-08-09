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

type Hub struct {
	clients map[string]map[*Client]bool

	broadcast chan *Broadcast

	register chan *RegisterInfo

	unregister chan *RegisterInfo
}

type AcceptMessage struct {
	Type       string        `json:"type"`
	IceServers []interface{} `json:"iceServers,omitempty"`
}

type RejectMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type AcceptMetadataMessage struct {
	Type       string        `json:"type"`
	Metadata   interface{}   `json:"authzMetadata,omitempty"`
	IceServers []interface{} `json:"iceServers,omitempty"`
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Broadcast),
		register:   make(chan *RegisterInfo),
		unregister: make(chan *RegisterInfo),
		clients:    make(map[string]map[*Client]bool),
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
			if h.clients[roomId] == nil {
				h.clients[roomId] = make(map[*Client]bool)
			}
			ok := len(h.clients[roomId]) < 2
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
				h.clients[roomId][client] = true
				client.SendJSON(msg)
			} else {
				msg := &AcceptMessage{
					Type: "accept",
				}
				h.clients[roomId][client] = true
				client.SendJSON(msg)
			}
		case registerInfo := <-h.unregister:
			roomId := registerInfo.roomId
			client := registerInfo.client
			if _, ok := h.clients[roomId][client]; ok {
				delete(h.clients[roomId], client)
				close(client.send)
			}
		case broadcast := <-h.broadcast:
			for client := range h.clients[broadcast.roomId] {
				if client.clientId != broadcast.client.clientId {
					select {
					case client.send <- broadcast.messages:
					default:
						close(client.send)
						delete(h.clients[broadcast.roomId], client)
					}
				}
			}
		}
	}
}
