package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleEditGroup(c echo.Context) error {
	return c.Render(http.StatusOK, "groupform", map[string]interface{}{
		"Title": "Hello",
	})
}
