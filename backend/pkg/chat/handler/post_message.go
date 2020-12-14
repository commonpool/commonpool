package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/chat"
	model2 "github.com/commonpool/backend/pkg/chat/handler/model"
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
func (chatHandler *ChatHandler) SendMessage(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "SendMessage")

	// Unmarshal request
	req := model2.SendMessageRequest{}
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

	loggedInUser, err := chatHandler.auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	err = chatHandler.chatService.SendMessage(ctx, &chat.Message{
		Key:            model.NewMessageKey(uuid.NewV4()),
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
