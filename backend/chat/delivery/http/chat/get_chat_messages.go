package chat

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetMessages
// @Summary Gets topic messages
// @Description This endpoint returns the messages for the given topic.
// @ID getMessages
// @Param take query int false "Number of messages to take" minimum(0) maximum(100) default(10)
// @Param skip query int false "Number of messages to skip" minimum(0) default(0)
// @Param channel query string true "Subscription id"
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {object} web.GetTopicMessagesResponse
// @Failure 400 {object} utils.Error
// @Router /chat/messages [get]
func (h *ChatHandler) GetMessages(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetMessages")

	loggedInSession, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return errors.ErrUnauthorized
	}
	loggedInUserKey := loggedInSession.GetUserKey()

	channelSrt := c.QueryParam("channel")
	if channelSrt == "" {
		return errors.ErrQueryParamRequired("channel")
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	before, err := utils.ParseBefore(c)
	if err != nil {
		return err
	}

	channelKey := model.NewConversationKey(channelSrt)

	messages, err := h.chatService.GetMessages(ctx, loggedInUserKey, channelKey, *before, take)
	if err != nil {
		return err
	}

	items := make([]*web.Message, len(messages.Messages.Items))
	for i, message := range messages.Messages.Items {
		items[i] = web.MapMessage(&message)
	}

	return c.JSON(http.StatusOK, web.GetTopicMessagesResponse{
		Messages: items,
	})
}
