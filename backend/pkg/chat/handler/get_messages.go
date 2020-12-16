package handler

import (
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/utils"
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
// @Success 200 {object} web.GetMessagesResponse
// @Failure 400 {object} utils.Error
// @Router /chat/messages [get]
func (h *Handler) GetMessages(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetMessages")

	channelSrt := c.QueryParam("channel")
	if channelSrt == "" {
		return exceptions.ErrQueryParamRequired("channel")
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	before, err := utils.ParseBefore(c)
	if err != nil {
		return err
	}

	channelKey := chat.NewConversationKey(channelSrt)

	messages, err := h.chatService.GetMessages(ctx, channelKey, *before, take)
	if err != nil {
		return err
	}

	items := make([]*Message, len(messages.Messages.Items))
	for i, message := range messages.Messages.Items {
		items[i] = MapMessage(&message)
	}

	return c.JSON(http.StatusOK, GetMessagesResponse{
		Messages: items,
	})
}

type GetMessagesResponse struct {
	Messages []*Message `json:"messages"`
}
