package ayame

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSignalingHandler(t *testing.T) {
	t.Skip("")
}

func TestSignalingHandlerNoWebSocketHeaders(t *testing.T) {
	e := echo.New()

	s := &Server{
		Server: http.Server{
			Handler: e,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/signaling", nil)
	//req.Header.Set("connection", "upgrade")
	//req.Header.Set("upgrade", "websocket")
	//req.Header.Set("Sec-WebSocket-Version", "13")
	//req.Header.Set("Sec-WebSocket-Key", "9KSB2xxx0G/vaPz4e5+ACw==")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.Error(t, s.signalingHandler(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
