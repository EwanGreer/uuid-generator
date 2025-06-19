package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleUp(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
