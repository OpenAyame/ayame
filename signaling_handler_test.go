package ayame

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
