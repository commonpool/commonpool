package handler

import (
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
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
func (h *Handler) SendMessage(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "SendMessage")

	// Unmarshal request
	req := SendMessageRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	// retrieve the thread
	channelId := c.Param("id")
	channelKey := chat.NewConversationKey(channelId)
	// todo verify that user has permission to post on topic

	loggedInUser, err := h.auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	err = h.chatService.SendMessage(ctx, &chat.Message{
		Key:            chat.NewMessageKey(uuid.NewV4()),
		ChannelKey:     channelKey,
		MessageType:    chat.NormalMessage,
		MessageSubType: chat.UserMessage,
		SentBy: chat.MessageSender{
			Type:     chat.UserMessageSender,
			UserKey:  loggedInUser.GetUserKey(),
			Username: loggedInUser.GetUsername(),
		},
		SentAt:        time.Now().UTC(),
		Text:          req.Message,
		Blocks:        nil,
		Attachments:   nil,
		VisibleToUser: nil,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)

}

type SendMessageRequest struct {
	Message string `json:"message,omitempty" validate:"notblank,required,min=1,max=2000"`
}
