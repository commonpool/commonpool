package handler

import (
	"cp/pkg/api"
	posts2 "cp/pkg/posts"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Query struct {
	Query *string       `query:"query"`
	Type  *api.PostType `query:"type"`
}

func (h *Handler) handleGroupPostsView(c echo.Context) error {

	group, err := h.getGroup(c)
	if err != nil {
		return err
	}

	var payload Query
	if err := c.Bind(&payload); err != nil {
		return err
	}

	if payload.Type != nil && string(*payload.Type) == "" {
		payload.Type = nil
	}
	if payload.Query != nil && string(*payload.Query) == "" {
		payload.Query = nil
	}

	posts, err := h.postStore.GetByGroup(group.ID, &posts2.FindPostsOptions{
		Query: payload.Query,
		Type:  payload.Type,
	})
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "group", map[string]interface{}{
		"Title": "Hello",
		"Posts": posts,
		"Query": payload.Query,
		"Type":  payload.Type,
	})
}
