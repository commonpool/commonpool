package handler

import (
	"cp/pkg/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleAdmin(c echo.Context) error {
	return c.Render(http.StatusOK, "admin", map[string]interface{}{
		"Title": "Admin",
	})
}

func (h *Handler) handleAdminClearAll(c echo.Context) error {
	if err := h.db.Where("1 = 1").Delete(&api.Acknowledgement{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&api.Credits{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&api.Message{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&api.Notification{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&api.Image{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&api.Post{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&api.Membership{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&api.Group{}).Error; err != nil {
		return err
	}
	if err := h.db.Where("1 = 1").Delete(&api.User{}).Error; err != nil {
		return err
	}
	c.Response().Header().Set("Location", "/auth/logout")
	c.Response().WriteHeader(http.StatusSeeOther)
	return nil
}
