package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleGetUserProfile(c echo.Context) error {
	return c.Render(http.StatusOK, "user_profile_view", map[string]interface{}{
		"Title": "Hello",
	})
}
