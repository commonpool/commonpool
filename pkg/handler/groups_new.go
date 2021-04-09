package handler

import (
	"cp/pkg/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
)

type CreateGroup struct {
	Name string `form:"name"`
}

func (h *Handler) handleNewGroup(c echo.Context) error {

	profile, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "groupform", map[string]interface{}{
			"Title": "Hello",
		})
	}

	var payload CreateGroup
	if err := c.Bind(&payload); err != nil {
		return err
	}

	group := &api.Group{
		ID:   uuid.NewV4().String(),
		Name: payload.Name,
	}
	if err := h.groupStore.Create(group); err != nil {
		return err
	}

	membership := &api.Membership{
		GroupID:    group.ID,
		UserID:     profile.ID,
		Permission: api.Owner,
		MemberConfirmed: true,
		GroupConfirmed: true,
	}
	if err := h.membershipStore.Create(membership); err != nil {
		return err
	}

	c.Response().Header().Set("Location", fmt.Sprintf("%s://%s/groups/%s", c.Scheme(), c.Request().Host, group.ID))
	c.Response().WriteHeader(http.StatusSeeOther)
	return nil

}
