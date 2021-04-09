package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleGetUserPosts(c echo.Context) error {
	user, err := h.getUser(c)
	if err != nil {
		return err
	}

	posts, err := h.postStore.GetByAuthor(user.ID)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "user_posts_view", map[string]interface{}{
		"Title": "Hello",
		"Posts": posts,
	})
}
