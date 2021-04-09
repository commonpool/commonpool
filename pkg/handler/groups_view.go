package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleGroupsView(c echo.Context) error {
	groups, err := h.groupStore.Search()
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "groups", map[string]interface{}{
		"Title":  "Hello",
		"Groups": groups,
	})
}
