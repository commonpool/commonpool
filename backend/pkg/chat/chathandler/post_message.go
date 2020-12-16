package chathandler

import (
	model2 "github.com/commonpool/backend/pkg/chat/chathandler/model"
	chatmodel "github.com/commonpool/backend/pkg/chat/chatmodel"
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
	channelKey := chatmodel.NewConversationKey(channelId)
	// todo verify that user has permission to post on topic

	loggedInUser, err := chatHandler.auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	err = chatHandler.chatService.SendMessage(ctx, &chatmodel.Message{
		Key:            chatmodel.NewMessageKey(uuid.NewV4()),
		ChannelKey:     channelKey,
		MessageType:    chatmodel.NormalMessage,
		MessageSubType: chatmodel.UserMessage,
		SentBy: chatmodel.MessageSender{
			Type:     chatmodel.UserMessageSender,
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
