package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleGetGroupAcknowledgements(c echo.Context) error {
	group, err := h.getGroup(c)
	if err != nil {
		return err
	}

	acknowledgements, err := h.acknowledgementStore.GetForGroup(group.ID)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "group_acknowledgements_view", map[string]interface{}{
		"Title": "Hello",
		"Acknowledgements": acknowledgements,
	})
}
