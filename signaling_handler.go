package ayame

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

func (s *Server) signalingHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		markHandlerError(w)
		zlog.Debug().Err(err).Send()
		return
	}

	wsConn.SetReadLimit(readLimit)

	// s.config.CopyWebSocketHeaderNames にマッチするヘッダーを取得して connection.copyHeaders に設定する
	copyHeaders := make(map[string]string)
	for _, headerName := range s.config.CopyWebSocketHeaderNames {
		// 設定ファイルのヘッダー名を http.CanonicalHeaderKey で正規化してから Header.Get で取得する
		canonicalName := http.CanonicalHeaderKey(headerName)
		if headerValue := r.Header.Get(canonicalName); headerValue != "" {
			copyHeaders[canonicalName] = headerValue
		}
	}

	connection := connection{
		// connectionId を設定する
		ID: getULID(),

		wsConn: wsConn,
		// 複数箇所でブロックした時を考えて少し余裕をもたせる
		forwardChannel: make(chan forward, 100),

		// config を connection でも触れるように渡しておく
		config:          *s.config,
		signalingLogger: *s.signalingLogger,
		webhookLogger:   *s.webhookLogger,
		metrics:         s.Metrics,

		copyHeaders: copyHeaders,
	}

	// client.conn.SetCloseHandler(client.closeHandler)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// ブロックしないよう余裕をもたせておく
	messageChannel := make(chan []byte, 100)
	go connection.wsRecv(ctx, messageChannel)
	go connection.main(cancel, messageChannel)
}
