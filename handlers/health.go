package handlers

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleUp(c echo.Context) error {
	return c.JSON(200, map[string]any{
		"version": h.Version,
	})
}
