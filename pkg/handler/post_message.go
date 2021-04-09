package handler

import (
	"cp/pkg/api"
	"cp/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type SubmitMessage struct {
	Content string `form:"content"`
}

func (h *Handler) handlePostMessage(c echo.Context) error {

	authenticatedUser, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	group, err := h.getGroup(c)
	if err != nil {
		return err
	}

	post, err := h.getPost(c)
	if err != nil {
		return err
	}

	var payload SubmitMessage
	if err := c.Bind(&payload); err != nil {
		return err
	}

	var message = &api.Message{
		ID:       uuid.NewV4().String(),
		AuthorID: authenticatedUser.ID,
		Content:  payload.Content,
		ThreadID: post.ID,
	}
	if err := h.messageStore.SendMessage(message); err != nil {
		return err
	}

	userIds, err := h.messageStore.FindUserIdsInThread(post.ID)
	if err != nil {
		return err
	}
	userIds = utils.UniqueStrings(append(userIds, post.AuthorID))
	users, err := h.userStore.GetByKeys(userIds)
	if err != nil {
		return err
	}
	userMap := utils.UserMap(users)
	delete(userMap, authenticatedUser.ID)

	var notifications []*api.Notification
	for _, user := range userMap {
		notifications = append(notifications, &api.Notification{
			ID:     uuid.NewV4().String(),
			UserID: user.ID,
			Title:  fmt.Sprintf("Post %s - New Message", post.HTMLLink()),
			Message: fmt.Sprintf("%s replied to post %s in group %s",
				user.HTMLLink(),
				post.HTMLLink(),
				group.HTMLLink()),
			Link: post.HTMLLink(),
		})
	}

	if err := h.notificationStore.AddNotifications(notifications); err != nil {
		return err
	}

	c.Response().Header().Set("Location", fmt.Sprintf("%s://%s/groups/%s/posts/%s", c.Scheme(), c.Request().Host, group.ID, post.ID))
	c.Response().WriteHeader(http.StatusSeeOther)
	return nil

}
