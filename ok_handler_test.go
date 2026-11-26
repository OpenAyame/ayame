package ayame

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOkHandler(t *testing.T) {
	s := &Server{}

	req := httptest.NewRequest(http.MethodGet, "/.ok", nil)
	rec := httptest.NewRecorder()
	s.okHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Body.String())
}
