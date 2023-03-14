package ayame

import (
	"context"
	"net/http"
	"time"

	zlog "github.com/rs/zerolog/log"
	"github.com/shiguredo/websocket"
)

const (
	writeWait = 10 * time.Second

	// ws の読み込みは最大 1MByte までにする
	readLimit = 1048576
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 4,
		WriteBufferSize: 1024 * 4,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (s *Server) signalingHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	wsConn.SetReadLimit(readLimit)
	if err != nil {
		zlog.Debug().Err(err).Send()
		return
	}
	// ここで connectionId みたいなの作るべき
	connection := connection{
		wsConn: wsConn,
		// 複数箇所でブロックした時を考えて少し余裕をもたせる
		forwardChannel: make(chan forward, 100),

		// config を connection でも触れるように渡しておく
		config:          *s.config,
		signalingLogger: *s.signalingLogger,
		webhookLogger:   *s.webhookLogger,
	}
	// client.conn.SetCloseHandler(client.closeHandler)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// ブロックしないよう余裕をもたせておく
	messageChannel := make(chan []byte, 100)
	go connection.wsRecv(ctx, messageChannel)
	go connection.main(cancel, messageChannel)

}
