package chat

import (
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// SendMessage
// @Summary Sends a message to a topic
// @Description This endpoint sends a message to the given thread
// @ID sendMessage
// @Param message body web.SendMessageRequest true "Message to send"
// @Param id path string true "channel id"
// @Tags chat
// @Accept json
// @Success 202
// @Failure 400 {object} utils.Error
// @Router /chat/:id [post]
func (h *ChatHandler) SendMessage(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "SendMessage")

	// Unmarshal request
	req := web.SendMessageRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	// retrieve the thread
	channelId := c.Param("id")
	channelKey := model.NewConversationKey(channelId)
	// todo verify that user has permission to post on topic

	_, err := h.chatService.SendChannelMessage(ctx, channelKey, req.Message)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)

}
