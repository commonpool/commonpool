package handler

import (
	"cp/pkg/api"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func (h *Handler) handleGroupSettings(c echo.Context) error {
	return c.Render(http.StatusOK, "group_settings", map[string]interface{}{
		"Title": "Hello",
	})
}

func (h *Handler) handleGroupDelete(c echo.Context) error {

	authenticatedUserMembership, err := h.getAuthenticatedUserMembership(c)
	if err != nil {
		return err
	}

	if !authenticatedUserMembership.IsOwner() {
		return echo.ErrForbidden
	}

	var groupID = authenticatedUserMembership.GroupID

	if err := h.db.Where("group_id = ?", groupID).Delete(&api.Acknowledgement{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("group_id = ?", groupID).Delete(&api.Credits{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("group_id = ?", groupID).Delete(&api.Membership{}).Error; err != nil {
		return err
	}

	var posts []*api.Post
	if err := h.db.Unscoped().Model(&api.Post{}).Find(&posts, "group_id = ?", groupID).Error; err != nil {
		return err
	}

	if len(posts) > 0 {
		var sb strings.Builder
		sb.WriteString("thread_id in (")
		var params []interface{}
		for i, post := range posts {
			sb.WriteString("?")
			params = append(params, post.ID)
			if i < len(posts)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")

		if err := h.db.Where(sb.String()).Delete(&api.Message{}).Error; err != nil {
			return err
		}
	}

	if err := h.db.Unscoped().Where("group_id = ?", groupID).Delete(&api.Post{}).Error; err != nil {
		return err
	}

	if err := h.db.Where("id = ?", groupID).Delete(&api.Group{}).Error; err != nil {
		return err
	}

	c.Response().Header().Set("Location", "/")
	c.Response().WriteHeader(http.StatusSeeOther)
	return nil
}
