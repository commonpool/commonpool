package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleGetUserNotifications(c echo.Context) error {

	user, err := h.getUser(c)
	if err != nil {
		return err
	}

	authenticatedUser, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	if authenticatedUser.ID != user.ID {
		return echo.ErrForbidden
	}

	notifications, err := h.notificationStore.GetNotifications(authenticatedUser.ID)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "user_notifications_view", map[string]interface{}{
		"Title": "Hello",
		"Notifications": notifications,
	})
}
