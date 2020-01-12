package main

var (
	registerChannel   = make(chan *register)
	unregisterChannel = make(chan *unregister)
	forwardChannel    = make(chan forward)
)

type room struct {
	id      string
	clients map[string]*client
}

func server() {
	// Rooms を管理するマップはここに用意する
	// roomId がキーになる
	var m = make(map[string]room)
	// ここはシングルなのでロックは不要、多分
	for {
		select {
		case register := <-registerChannel:
			c := register.client
			rch := register.resultChannel
			r, ok := m[c.roomID]
			if ok {
				// room があった
				if len(r.clients) == 1 {
					// room に 自分を追加する
					// 登録しているのが同じ ID だった場合はエラーにする
					_, ok := r.clients[c.ID]
					if ok {
						// 重複エラー
						rch <- dup
					} else {
						r.clients[c.ID] = c
						m[c.roomID] = r
						rch <- two
					}
				} else {
					// room あったけど満杯
					rch <- full
				}
			} else {
				// room がなかった
				var clients = make(map[string]*client)
				clients[c.ID] = c
				// room を追加
				m[c.roomID] = room{
					id:      c.roomID,
					clients: clients,
				}
				c.debugLog().Msg("CREATED-ROOM")
				rch <- one
			}
		case unregister := <-unregisterChannel:
			c := unregister.client
			// 部屋を探す
			r, ok := m[c.roomID]
			if ok {
				_, ok := r.clients[c.ID]
				if ok {
					for _, client := range r.clients {
						// 両方の forwardChannel を閉じる
						close(client.forwardChannel)
						c.debugLog().Msg("CLOSED-FORWARD-CHANNEL")
						c.debugLog().Msg("REMOVED-CLIENT")
					}
					// room を削除
					delete(m, c.roomID)
					c.debugLog().Msg("DELETED-ROOM")
				}
			} else {
				// 部屋がなかった
				// 何もしない
			}
		case forward := <-forwardChannel:
			r, ok := m[forward.client.roomID]
			if ok {
				// room があった
				for clientID, client := range r.clients {
					if clientID != forward.client.ID {
						client.forwardChannel <- forward
					}
				}
			} else {
				// room がなかった
				logger.Warn().Msg("MISSING-ROOM")
			}
		}
	}
}
