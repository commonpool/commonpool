package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleGroupMembersView(c echo.Context) error {
	return c.Render(http.StatusOK, "group_members", map[string]interface{}{
		"Title":      "Hello",
	})
}
