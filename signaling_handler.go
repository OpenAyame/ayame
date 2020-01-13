package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 10 * time.Second

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
	conn, err := upgrader.Upgrade(w, r, nil)
	conn.SetReadLimit(readLimit)
	if err != nil {
		logger.Debug().Err(err).Caller().Msg("")
		return
	}
	// ここで connectionId みたいなの作るべき
	client := client{
		conn:           conn,
		forwardChannel: make(chan forward),
	}
	// client.conn.SetCloseHandler(client.closeHandler)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	messageChannel := make(chan []byte)
	go client.wsRecv(ctx, messageChannel)
	go client.main(cancel, messageChannel)

}
