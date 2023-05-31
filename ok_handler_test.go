package ayame

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestOkHandler(t *testing.T) {
	e := echo.New()

	s := &Server{
		Server: http.Server{
			Handler: e,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/.ok", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, s.okHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
