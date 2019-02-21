package main

type Broadcast struct {
	client   *Client
	roomId   string
	messages []byte
}

type RegisterInfo struct {
	roomId string
	client *Client
}

type Hub struct {
	clients map[string]map[*Client]bool

	broadcast chan *Broadcast

	register chan *RegisterInfo

	unregister chan *RegisterInfo
}

type AcceptMessage struct {
	Type string `json:"type"`
}

type RejectMessage struct {
	Type string `json:"type"`
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
			roomId := registerInfo.roomId
			client := registerInfo.client
			if h.clients[roomId] == nil {
				h.clients[roomId] = make(map[*Client]bool)
			}
			h.clients[roomId][client] = true
			ok := len(h.clients[roomId]) < 3
			if !ok {
				msg := &RejectMessage{
					Type: "reject",
				}
				client.conn.WriteJSON(msg)
				client.conn.Close()
				break
			}
			msg := &AcceptMessage{
				Type: "accept",
			}
			client.conn.WriteJSON(msg)
		case registerInfo := <-h.unregister:
			roomId := registerInfo.roomId
			client := registerInfo.client
			if _, ok := h.clients[roomId][client]; ok {
				delete(h.clients[roomId], client)
				close(client.send)
			}
		case broadcast := <-h.broadcast:
			for client := range h.clients[broadcast.roomId] {
				if client.uuid != broadcast.client.uuid {
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
