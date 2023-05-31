package ayame

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) okHandler(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
