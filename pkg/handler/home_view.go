package handler

import (
	"cp/pkg/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleHomeView(c echo.Context) error {

	var groups []*api.Group

	authenticatedUser, err := h.getAuthenticatedUser(c)
	if err == nil && authenticatedUser != nil{
		groups, err = h.groupStore.Search()
		if err != nil {
			return err
		}
	}

	return c.Render(http.StatusOK, "index", map[string]interface{}{
		"Title":  "Hello",
		"Groups": groups,
	})
}
