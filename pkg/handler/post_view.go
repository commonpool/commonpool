package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handlePostView(c echo.Context) error {

	post, err := h.getPost(c)
	if err != nil {
		return err
	}

	messages, err := h.messageStore.GetMessages(post.ID)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "post_view", map[string]interface{}{
		"Title":    "Hello",
		"Messages": messages,
	})
}
