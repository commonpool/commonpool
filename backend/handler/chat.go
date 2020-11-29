package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/chat"
	. "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

/**

For designing a scalable chat system, we will de-normalize the chat data.

The chat messages are stored in threads.
The threads are owned by a single user.
Messages have an authorId, and a threadId, and a message.
The message userId could be retrieved by foreign key with the threadId, but we keep a message.userId for quicker querying

If a user sends a message to another user, the message will be stored in the sender thread and in the recipient thread.
This way we can partition whatever data persistence we want upon the recipient id.

*/

// GetRecentlyActiveSubscriptions
// @Summary Returns the latest user message threads
// @Description This endpoint returns the latest messaging threads for the currently logged in user.
// @ID getLatestThreads
// @Param take query int false "Number of threads to take" minimum(0) maximum(100) default(10)
// @Param skip query int false "Number of threads to skip" minimum(0) default(0)
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {object} web.GetLatestSubscriptionsResponse
// @Failure 400 {object} utils.Error
// @Router /chat/subscriptions [get]
func (h *Handler) GetRecentlyActiveSubscriptions(c echo.Context) error {

	ctx, l := GetEchoContext(c, "GetRecentlyActiveSubscriptions")

	l.Debug("getting user auth session")

	authUser := h.authorization.GetAuthUserSession2(ctx)
	userKey := model.NewUserKey(authUser.Subject)

	l.Debug("parsing skip query param")

	skip, err := utils.ParseSkip(c)
	if err != nil {
		l.Error("could not parse 'skip' query param", zap.Error(err))
		return err
	}

	l.Debug("parsing take query param")

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		l.Error("could not parse 'take' query param", zap.Error(err))
		return err
	}

	l.Debug("getting subscriptions",
		zap.Int("take", take),
		zap.Int("skip", skip))

	userSubscriptions, err := h.chatStore.GetSubscriptionsForUser(ctx, chat.NewGetSubscriptions(userKey, take, skip))
	if err != nil {
		l.Error("could not get subscriptions", zap.Error(err))
		return err
	}

	l.Debug("mapping subscriptions to web response", zap.Int("subscriptionCount", len(userSubscriptions.Items)))

	var items []web.Subscription
	for _, subscription := range userSubscriptions.Items {

		l := l.With(zap.String("channelId", subscription.ChannelID))

		l.Debug("getting channel")

		channel, err := h.chatStore.GetChannel(ctx, subscription.GetChannelKey())
		if err != nil {
			return err
		}

		l.Debug("mapping to web response")

		items = append(items, web.Subscription{
			ChannelID:           channel.ID,
			UserID:              subscription.UserID,
			HasUnreadMessages:   subscription.LastMessageAt.After(subscription.LastTimeRead),
			CreatedAt:           subscription.CreatedAt,
			UpdatedAt:           subscription.UpdatedAt,
			LastMessageAt:       subscription.LastMessageAt,
			LastTimeRead:        subscription.LastTimeRead,
			LastMessageChars:    subscription.LastMessageChars,
			LastMessageUserId:   subscription.LastMessageUserId,
			LastMessageUserName: subscription.LastMessageUserName,
			Name:                subscription.Name,
			Type:                channel.Type,
		})

	}

	if items == nil {
		items = []web.Subscription{}
	}

	return c.JSON(http.StatusOK, web.GetLatestSubscriptionsResponse{
		Subscriptions: items,
	})

}

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
func (h *Handler) GetMessages(c echo.Context) error {

	ctx, l := GetEchoContext(c, "GetMessages")

	l.Debug("getting user session")

	loggedInSession := h.authorization.GetAuthUserSession(c)
	loggedInUserKey := loggedInSession.GetUserKey()

	l.Debug("parsing 'channel' query param")

	channelSrt := c.QueryParam("channel")
	if channelSrt == "" {
		return fmt.Errorf("'channel' query param is required")
	}

	l.Debug("parsing 'take' query param")

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		l.Error("could not parse 'take' query param", zap.Error(err))
		return err
	}

	l.Debug("parsing 'before' query param")

	before, err := utils.ParseBefore(c)
	if err != nil {
		l.Error("could not parse 'before' query param", zap.Error(err))
		return err
	}

	channelKey := model.NewConversationKey(channelSrt)

	l.Debug("getting messages")

	getMessages, err := h.chatStore.GetMessages(ctx, chat.NewGetMessages(loggedInUserKey, channelKey, *before, take))
	if err != nil {
		l.Error("could not get messages query param", zap.Error(err))
		return err
	}

	l.Debug("mapping to web response")

	items := make([]*web.Message, len(getMessages.Messages.Items))
	for i, message := range getMessages.Messages.Items {
		items[i] = mapMessage(&message)
	}

	return c.JSON(http.StatusOK, web.GetTopicMessagesResponse{
		Messages: items,
	})
}

func mapMessage(message *chat.Message) *web.Message {
	var visibleToUser *string = nil
	if message.VisibleToUser != nil {
		visibleToUserStr := message.VisibleToUser.String()
		visibleToUser = &visibleToUserStr
	}
	return &web.Message{
		ID:             message.Key.String(),
		ChannelID:      message.ChannelKey.String(),
		MessageType:    message.MessageType,
		MessageSubType: message.MessageSubType,
		SentById:       message.SentBy.UserKey.String(),
		SentByUsername: message.SentBy.USername,
		SentAt:         message.SentAt,
		Text:           message.Text,
		Blocks:         message.Blocks,
		Attachments:    message.Attachments,
		VisibleToUser:  visibleToUser,
	}
}

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

	ctx, l := GetEchoContext(c, "InquireAboutResource")

	l.Debug("getting logged in user")

	// Get current user
	loggedInUser := h.authorization.GetAuthUserSession(c)
	loggedInUserKey := model.NewUserKey(loggedInUser.Subject)

	l.Debug("getting 'id' resource key query param")

	resourceKey, err := model.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}

	l.Debug("unmarshaling request")

	req := web.InquireAboutResourceRequest{}
	if err := c.Bind(&req); err != nil {

		l.Error("could not unmarshal request", zap.Error(err))

		response := ErrSendResourceMsgBadRequest(err)
		return &response
	}

	l.Debug("validating request")

	if err := c.Validate(req); err != nil {
		l.Warn("bad request payload", zap.Error(err))
		return ErrValidation(err.Error())
	}

	// todo: send the channel id back to the client so he can redirect
	_, err = h.chatService.NotifyUserInterestedAboutResource(
		ctx, chat.NewNotifyUserInterestedAboutResource(loggedInUserKey, *resourceKey, req.Message))

	if err != nil {
		l.Error("could not notify user interested about resource", zap.Error(err))
		return err
	}

	return c.NoContent(http.StatusAccepted)

}

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

	ctx, l := GetEchoContext(c, "SendMessage")

	l.Debug("unmarshaling request")

	// Unmarshal request
	req := web.SendMessageRequest{}
	if err := c.Bind(&req); err != nil {
		response := ErrSendResourceMsgBadRequest(err)
		return &response
	}

	l.Debug("validating request")

	if err := c.Validate(req); err != nil {
		return ErrValidation(err.Error())
	}

	l.Debug("getting channel 'id' query param")

	// retrieve the thread
	channelId := c.Param("id")
	channelKey := model.NewConversationKey(channelId)
	// todo verify that user has permission to post on topic

	l.Debug("sending message")

	_, err := h.chatService.SendChannelMessage(ctx, channelKey, req.Message)
	if err != nil {
		l.Error("could not send message", zap.Error(err))
		return err
	}

	return c.NoContent(http.StatusAccepted)

}

// SubmitInteraction
// @Summary Sends a message to a topic
// @Description This endpoint is for user interactions through the chat box
// @ID submitInteraction
// @Param message body web.SubmitInteractionRequest true "Message to send"
// @Tags chat
// @Accept json
// @Success 200
// @Failure 400 {object} utils.Error
// @Router /chat/interaction [post]
func (h *Handler) SubmitInteraction(c echo.Context) error {

	ctx, l := GetEchoContext(c, "SubmitInteraction")

	// Get current user
	authSession := h.authorization.GetAuthUserSession(c)
	authUserKey := model.NewUserKey(authSession.Subject)

	// Unmarshal request
	req := web.SubmitInteractionRequest{}
	if err := c.Bind(&req); err != nil {
		l.Error("could not unmarshal request", zap.Error(err))
		response := ErrSendResourceMsgBadRequest(err)
		return &response
	}
	if err := c.Validate(req); err != nil {
		l.Warn("error validating request", zap.Error(err))
		return ErrValidation(err.Error())
	}

	// Getting the message
	uid, err := uuid.FromString(req.Payload.MessageID)
	if err != nil {
		l.Warn("could not convert message id to uuid", zap.Error(err))
		return err
	}
	message, err := h.chatStore.GetMessage(ctx, model.NewMessageKey(uid))
	if err != nil {
		l.Warn("could not get message", zap.Error(err))
		return err
	}

	now := time.Now()

	// Mapping message actions
	var actions []web.Action
	for _, action := range req.Payload.Actions {
		actions = append(actions, web.Action{
			SubmitAction: web.SubmitAction{
				ElementState: action.ElementState,
				BlockID:      action.BlockID,
				ActionID:     action.ActionID,
			},
			ActionTimestamp: now,
		})
	}

	webMessage := mapMessage(message)

	// Creating interaction payload message
	interactionPayload := web.InteractionCallback{
		Token: h.config.CallbackToken,
		Payload: web.InteractionCallbackPayload{
			Type:        web.BlockActions,
			TriggerId:   "",
			ResponseURL: "",
			User: web.InteractionPayloadUser{
				ID:       authUserKey.String(),
				Username: authSession.Username,
			},
			Message: webMessage,
			Actions: actions,
			State:   req.Payload.State,
		},
	}

	requestBody, err := json.Marshal(interactionPayload)
	if err != nil {
		l.Error("could not convert interaction payload to request body", zap.Error(err))
		return err
	}

	httpRequest, err := http.NewRequest("POST", "http://localhost:8585/api/v1/chatback", bytes.NewBuffer(requestBody))
	if err != nil {
		l.Error("error occurred while creating the chatback query", zap.Error(err))
		return err
	}

	l.Debug("Token: " + c.Get("token").(string))

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+c.Get("token").(string))

	response, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		l.Error("error occurred while calling the chatback api", zap.Error(err))
		return err
	}

	if response.StatusCode != 200 {
		l.Error("unexpected chatback return code", zap.String("status", response.Status))
		return fmt.Errorf("unexpected status code")
	}

	return c.String(http.StatusOK, "OK")

}
