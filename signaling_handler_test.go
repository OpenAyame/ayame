package ayame

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignalingHandler(t *testing.T) {
	t.Skip("")
}

func TestSignalingHandlerNoWebSocketHeaders(t *testing.T) {
	s := &Server{}

	req := httptest.NewRequest(http.MethodGet, "/signaling", nil)
	rec := httptest.NewRecorder()
	s.signalingHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// hijackableRecorder は Hijack で conn と brw を返す ResponseWriter である。
type hijackableRecorder struct {
	*httptest.ResponseRecorder
	conn net.Conn
	brw  *bufio.ReadWriter
}

func (h *hijackableRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return h.conn, h.brw, nil
}

func TestSignalingHandlerUpgradeFailedAfterHijack(t *testing.T) {
	cases := []struct {
		name string
		// newWriter は Hijack した後に Upgrade が失敗する ResponseWriter を返す。
		newWriter func(t *testing.T) http.ResponseWriter
	}{
		{
			name: "data before handshake",
			newWriter: func(t *testing.T) http.ResponseWriter {
				server, client := net.Pipe()
				t.Cleanup(func() { client.Close() })
				// クライアントがハンドシェイク前に送ったデータを読み込み済みにする。
				br := bufio.NewReader(io.MultiReader(strings.NewReader("data"), server))
				_, err := br.Peek(1)
				require.NoError(t, err)
				return &hijackableRecorder{
					ResponseRecorder: httptest.NewRecorder(),
					conn:             server,
					brw:              bufio.NewReadWriter(br, bufio.NewWriter(server)),
				}
			},
		},
		{
			name: "write 101 failed",
			newWriter: func(t *testing.T) http.ResponseWriter {
				server, client := net.Pipe()
				// 相手側を閉じて 101 の書き込みを失敗させる。
				require.NoError(t, client.Close())
				return &hijackableRecorder{
					ResponseRecorder: httptest.NewRecorder(),
					conn:             server,
					brw:              bufio.NewReadWriter(bufio.NewReader(server), bufio.NewWriter(server)),
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestServer(t)
			req := httptest.NewRequest(http.MethodGet, "/signaling", nil)
			req.Host = "ayame.local"
			req.Header.Set("Connection", "Upgrade")
			req.Header.Set("Upgrade", "websocket")
			req.Header.Set("Sec-WebSocket-Version", "13")
			req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
			s.Handler.ServeHTTP(tc.newWriter(t), req)

			assert.Contains(t, scrapeMetrics(t, s), `ayame_requests_total{code="500",host="ayame.local",method="GET",url="/signaling"} 1`)
		})
	}
}
