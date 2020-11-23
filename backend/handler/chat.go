package handler

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/chat"
	. "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/utils"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
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

// GetLatestThreads
// @Summary Returns the latest user message threads
// @Description This endpoint returns the latest messaging threads for the currently logged in user.
// @ID getLatestThreads
// @Param take query int false "Number of threads to take" minimum(0) maximum(100) default(10)
// @Param skip query int false "Number of threads to skip" minimum(0) default(0)
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {object} web.GetLatestThreadsResponse
// @Failure 400 {object} utils.Error
// @Router /chat/threads [get]
func (h *Handler) GetLatestThreads(c echo.Context) error {

	var err error

	authUser := h.authorization.GetAuthUserSession(c)
	userKey := model.NewUserKey(authUser.Subject)

	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	threads, err := h.chatStore.GetLatestThreads(userKey, take, skip)
	if err != nil {
		return err
	}

	var items = make([]web.Thread, len(threads))
	for i, thread := range threads {

		topic, err := h.chatStore.GetTopic(thread.GetKey().TopicKey)
		if err != nil {
			return err
		}

		fmt.Println("!!!")
		fmt.Println(userKey.String())
		fmt.Println(thread.LastTimeRead)
		fmt.Println(thread.LastMessageAt)

		before := thread.LastMessageAt.After(thread.LastTimeRead)
		fmt.Println("Has unread : ", before)

		items[i] = web.Thread{
			TopicID:             thread.TopicID,
			RecipientID:         thread.UserID,
			LastChars:           thread.LastMessageChars,
			HasUnreadMessages:   before,
			LastMessageAt:       thread.LastMessageAt,
			LastMessageUserId:   thread.LastMessageUserId,
			LastMessageUsername: thread.LastMessageUserName,
			Title:               topic.Title,
		}
	}

	return c.JSON(http.StatusOK, web.GetLatestThreadsResponse{
		Threads: items,
	})

}

// GetMessages
// @Summary Gets topic messages
// @Description This endpoint returns the messages for the given topic.
// @ID getMessages
// @Param take query int false "Number of messages to take" minimum(0) maximum(100) default(10)
// @Param skip query int false "Number of messages to skip" minimum(0) default(0)
// @Param topic query string true "Topic id"
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {object} web.GetTopicMessagesResponse
// @Failure 400 {object} utils.Error
// @Router /chat/messages [get]
func (h *Handler) GetMessages(c echo.Context) error {

	var err error

	authUser := h.authorization.GetAuthUserSession(c)
	userKey := model.NewUserKey(authUser.Subject)

	topicStr := c.QueryParam("topic")
	if topicStr == "" {
		return fmt.Errorf("'topic' query param is required")
	}

	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	topicKey := model.NewTopicKey(topicStr)
	threadKey := model.NewThreadKey(topicKey, userKey)

	messages, err := h.chatStore.GetThreadMessages(threadKey, take, skip)
	if err != nil {
		return err
	}

	authors := map[string]string{}

	items := make([]web.Message, len(messages))
	for i, message := range messages {

		if message.MessageType == model.NormalMessage && message.MessageSubType == model.UserMessage {
			author := model.User{}
			err = h.authStore.GetByKey(message.GetAuthorKey(), &author)
			if err != nil {
				return err
			}
			authors[author.ID] = author.Username
		}

		var blocks []model.Block
		_ = json.Unmarshal([]byte(message.Blocks), &blocks)

		var attachments []model.Attachment
		_ = json.Unmarshal([]byte(message.Attachments), &attachments)

		item := web.Message{
			ID:             message.ID.String(),
			SentAt:         message.SentAt,
			TopicID:        message.TopicID,
			Blocks:         blocks,
			Attachments:    attachments,
			MessageType:    message.MessageType,
			MessageSubType: message.MessageSubType,
			Text:           message.Text,
			UserID:         message.UserID,
			BotID:          message.BotID,
			IsPersonal:     message.IsPersonal,
			SentBy:         message.SentBy,
			SentByUsername: message.SentByUsername,
		}
		items[i] = item
	}

	return c.JSON(http.StatusOK, web.GetTopicMessagesResponse{
		Messages: items,
	})
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

	/**

	this one is tricky. The problem is in designing the data model, and how the chat thread between users are
	permitted.

	if we only allow one thread between 2 users (a user can only send messages to another user on the same thread),
	then it is simple. This step would only have to retrieve or create the thread between the 2 users, and post
	a message to it.

	if we want to allow multiple threads between 2 users (eg. a user inquires about 2 resources from the same owner.
	these two convos would be on 2 threads instead of a single one) - then we have to create a mapping between
	the (interested user ⋅ resource) → thread id.

	We first check if there already exist a (interested user ⋅ resource) → thread id mapping. If not, we create
	a new thread and mapping.

	*/

	var err error

	// Get current user
	loggedInUser := h.authorization.GetAuthUserSession(c)
	loggedInUserKey := model.NewUserKey(loggedInUser.Subject)

	// Unmarshal request
	req := web.InquireAboutResourceRequest{}
	if err := c.Bind(&req); err != nil {
		response := ErrSendResourceMsgBadRequest(err)
		return &response
	}
	if err := c.Validate(req); err != nil {
		response := ErrValidation(err.Error())
		return &response
	}

	// get the resource key
	resourceKeyStr := c.Param("id")
	resourceKey, err := model.ParseResourceKey(resourceKeyStr)
	if err != nil {
		return err
	}

	// retrieve the resource
	getResourceByKeyResponse := h.resourceStore.GetByKey(resource.NewGetResourceByKeyQuery(*resourceKey))
	if getResourceByKeyResponse.Error != nil {
		return getResourceByKeyResponse.Error
	}
	res := getResourceByKeyResponse.Resource

	// make sure auth user is not resource owner
	// doesn't make sense for one to inquire about his own stuff
	if res.GetUserKey() == loggedInUserKey {
		err := ErrCannotInquireAboutOwnResource()
		return NewErrResponse(c, &err)
	}

	// get or create the (interested user ⋅ resource) → topic mapping
	userResourceTopic, err := h.chatStore.GetOrCreateResourceTopicMapping(*resourceKey, loggedInUserKey, h.resourceStore)
	if err != nil {
		return err
	}

	// send a message on that topic
	sendMessageToThreadRequest := chat.NewSendMessageToThreadRequest(
		model.NewThreadKey(userResourceTopic.GetTopicKey(), res.GetUserKey()),
		loggedInUserKey,
		loggedInUser.Username,
		req.Message,
		[]model.Block{
			*model.NewHeaderBlock(model.NewMarkdownObject("Someone is interested in your stuff!"), nil),
			*model.NewContextBlock([]model.BlockElement{
				model.NewMarkdownObject(
					fmt.Sprintf("[%s](%s) is interested by your post [%s](%s).",
						loggedInUser.Username,
						h.config.BaseUri+"/users/"+loggedInUser.Subject,
						res.Summary,
						h.config.BaseUri+"/users/"+res.CreatedBy+"/"+res.ID.String(),
					),
				),
			}, nil),
		},
		[]model.Attachment{},
	)
	sendMessageToThreadResponse := h.chatStore.SendMessageToThread(&sendMessageToThreadRequest)
	if sendMessageToThreadResponse.Error != nil {
		return sendMessageToThreadResponse.Error
	}

	sendMessageRequest := chat.NewSendMessageRequest(
		userResourceTopic.GetTopicKey(),
		loggedInUserKey,
		loggedInUser.Username,
		req.Message,
		[]model.Block{},
		[]model.Attachment{},
	)
	sendMessageResponse := h.chatStore.SendMessage(&sendMessageRequest)

	if sendMessageResponse.Error != nil {
		return NewErrResponse(c, sendMessageResponse.Error)
	}

	return c.NoContent(http.StatusAccepted)

}

// SendMessage
// @Summary Sends a message to a topic
// @Description This endpoint sends a message to the given thread
// @ID sendMessage
// @Param message body web.SendMessageRequest true "Message to send"
// @Param id path string true "Topic id"
// @Tags chat
// @Accept json
// @Success 202
// @Failure 400 {object} utils.Error
// @Router /chat/:id [post]
func (h *Handler) SendMessage(c echo.Context) error {

	/**
	In this method, the user is replying to an existing topic.
	*/

	// Get current user
	authUser := h.authorization.GetAuthUserSession(c)
	userKey := model.NewUserKey(authUser.Subject)

	// Unmarshal request
	req := web.SendMessageRequest{}
	if err := c.Bind(&req); err != nil {
		response := ErrSendResourceMsgBadRequest(err)
		return &response
	}
	if err := c.Validate(req); err != nil {
		response := ErrValidation(err.Error())
		return &response
	}

	// retrieve the thread
	topicIdStr := c.Param("id")
	topicKey := model.NewTopicKey(topicIdStr)

	// todo verify that user has permission to post on topic

	block := model.NewSectionBlock(model.NewPlainTextObject(req.Message), nil, nil, nil)
	sendMessageRequest := chat.NewSendMessageRequest(
		topicKey,
		userKey,
		authUser.Username,
		req.Message,
		[]model.Block{*block},
		[]model.Attachment{},
	)

	sendMessageResponse := h.chatStore.SendMessage(&sendMessageRequest)
	if sendMessageResponse.Error != nil {
		return NewErrResponse(c, sendMessageResponse.Error)
	}

	return c.NoContent(http.StatusAccepted)

}


func (h *Handler) SubmitInteraction(c echo.Context) error {

	authUser := h.authorization.GetAuthUserSession(c)
	userKey := model.NewUserKey(authUser.Subject)

	// Unmarshal request
	req := web.Message{}


	return nil

}