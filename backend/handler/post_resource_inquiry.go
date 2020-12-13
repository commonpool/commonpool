package handler

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// InquireAboutResource
// @Summary Sends a message to the user about a resource
// @Description This endpoint sends a message to the resource owner
// @ID inquireAboutResource
// @Param message body web.InquireAboutResourceRequest true "Message to send"
// @Param id path string true "Resource id"
// @Tags resources
// @Accept json
// @Success 202
// @Failure 400 {object} utils.Error
// @Router /resources/:id/inquire [post]
func (h *Handler) InquireAboutResource(c echo.Context) error {
	var err error

	ctx, _ := GetEchoContext(c, "InquireAboutResource")

	// Get current user
	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := model.NewUserKey(loggedInUser.Subject)

	resourceKey, err := model.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}

	req := web.InquireAboutResourceRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	// todo: send the channel id back to the client so he can redirect
	_, err = h.chatService.NotifyUserInterestedAboutResource(
		ctx, chat.NewNotifyUserInterestedAboutResource(loggedInUserKey, resourceKey, req.Message))

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)

}
