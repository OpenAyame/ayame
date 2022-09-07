package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	zlog "github.com/rs/zerolog/log"
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

func signalingHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	wsConn.SetReadLimit(readLimit)
	if err != nil {
		zlog.Debug().Err(err).Caller().Msg("")
		return
	}
	// ここで connectionId みたいなの作るべき
	connection := connection{
		wsConn: wsConn,
		// 複数箇所でブロックした時を考えて少し余裕をもたせる
		forwardChannel: make(chan forward, 100),
	}
	// client.conn.SetCloseHandler(client.closeHandler)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// ブロックしないよう余裕をもたせておく
	messageChannel := make(chan []byte, 100)
	go connection.wsRecv(ctx, messageChannel)
	go connection.main(cancel, messageChannel)

}
