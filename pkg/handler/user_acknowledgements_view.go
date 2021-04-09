package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleGetUserAcknowledgements(c echo.Context) error {
	user, err := h.getUser(c)
	if err != nil {
		return err
	}

	acknowledgements, err := h.acknowledgementStore.GetForUser(user.ID)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "user_acknowledgements_view", map[string]interface{}{
		"Title": "Hello",
		"Acknowledgements": acknowledgements,
	})
}
