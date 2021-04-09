package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handlePostDelete(c echo.Context) error {

	authenticatedUser, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	post, err := h.getPost(c)
	if err != nil {
		return err
	}

	if authenticatedUser.ID != post.AuthorID {
		return echo.ErrForbidden
	}

	if err := h.postStore.Delete(post.ID); err != nil {
		return err
	}

	if err := h.messageStore.DeleteThread(post.ID); err != nil {
		return err
	}

	referer := c.Request().Header.Get("Referer")
	postURL := fmt.Sprintf("%s://%s/groups/%s/posts/%s", c.Scheme(), c.Request().Host, post.GroupID, post.ID)

	if referer == postURL {
		c.Response().Header().Set("Location", fmt.Sprintf("%s://%s/users/%s", c.Scheme(), c.Request().Host, post.AuthorID))
	} else {
		c.Response().Header().Set("Location", referer)
	}

	c.Response().WriteHeader(http.StatusSeeOther)
	return nil

}
