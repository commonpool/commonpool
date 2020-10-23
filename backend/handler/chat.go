package handler

import (
	"fmt"
	. "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
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

	skip, err := ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	threads, err := h.chatStore.GetLatestThreads(userKey, take, skip)
	if err != nil {
		return err
	}

	var items = make([]web.Thread, len(threads))
	for i, thread := range threads {

		var user = model.User{}
		err := h.authStore.GetByKey(model.NewUserKey(thread.UserID), &user)
		if err != nil {
			return err
		}

		items[i] = web.Thread{
			TopicID:             thread.TopicID.String(),
			RecipientID:         thread.UserID,
			RecipientUsername:   user.Username,
			LastChars:           thread.LastMessageChars,
			HasUnreadMessages:   thread.LastTimeRead.Before(thread.LastMessageAt),
			LastMessageAt:       thread.LastMessageAt,
			LastMessageUserId:   thread.LastMessageUserId,
			LastMessageUsername: thread.LastMessageUserName,
		}
	}

	return c.JSON(http.StatusOK, web.GetLatestThreadsResponse{
		Threads: items,
	})

}

// GetMessages
// @Summary Gets thread messages
// @Description This endpoint returns the messages for the given threads.
// @ID getThreadMessages
// @Param take query int false "Number of messages to take" minimum(0) maximum(100) default(10)
// @Param skip query int false "Number of messages to skip" minimum(0) default(0)
// @Param topic query string thread "Topic id"
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {object} web.GetThreadMessagesResponse
// @Failure 400 {object} utils.Error
// @Router /chat/messages [get]
func (h *Handler) GetMessages(c echo.Context) error {

	var err error

	authUser := h.authorization.GetAuthUserSession(c)
	userKey := model.NewUserKey(authUser.Subject)

	topicStr := c.Param("topic")
	if topicStr == "" {
		return fmt.Errorf("'topic' query param is required")
	}

	skip, err := ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	topicId, err := uuid.FromString(topicStr)
	if err != nil {
		return err
	}

	topicKey := model.NewTopicKey(topicId)
	threadKey := model.NewThreadKey(topicKey, userKey)

	messages, err := h.chatStore.GetThreadMessages(threadKey, take, skip)
	if err != nil {
		return err
	}

	authors := map[string]string{}

	items := make([]web.Message, len(messages))
	for i, message := range messages {

		if _, ok := authors[message.AuthorID]; !ok {
			author := model.User{}
			err = h.authStore.GetByKey(message.GetAuthorKey(), &author)
			if err != nil {
				return err
			}
			authors[author.ID] = author.Username
		}

		item := web.Message{
			ID:             message.ID,
			SentAt:         message.SentAt,
			TopicID:        message.TopicId.String(),
			Content:        message.Content,
			SentBy:         message.AuthorID,
			SentByUsername: authors[message.AuthorID],
		}
		items[i] = item
	}

	return c.JSON(http.StatusOK, web.GetThreadMessagesResponse{
		Messages: items,
	})
}

// InquireAboutResource
// @Summary Sends a message to the user about a resource
// @Description This endpoint sends a message to the resource owner
// @ID inquireAboutResource
// @Param message body web.InquireAboutResourceRequest true "Message to send"
// @Tags chat
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
	authUser := h.authorization.GetAuthUserSession(c)
	userKey := model.NewUserKey(authUser.Subject)

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
	resource := &model.Resource{}
	err = h.resourceStore.GetByKey(*resourceKey, resource)
	if err != nil {
		return err
	}

	// make sure auth user is not resource owner
	if resource.GetUserKey() == userKey {
		err := ErrCannotInquireAboutOwnResource()
		return NewErrResponse(c, &err)
	}

	// get or create the (interested user ⋅ resource) → topic mapping
	userResourceTopic, err := h.chatStore.GetOrCreateResourceTopicMapping(*resourceKey, userKey, h.resourceStore)
	if err != nil {
		return err
	}

	// send a message on that topic
	err = h.chatStore.SendMessage(userKey, authUser.Username, userResourceTopic.GetTopicKey(), req.Message)
	if err != nil {
		return NewErrResponse(c, err)
	}

	return c.NoContent(http.StatusAccepted)

}

// SendMessage
// @Summary Sends a message to a thread
// @Description This endpoint sends a message to the given thread
// @ID sendMessage
// @Param message body web.SendMessageRequest true "Message to send"
// @Tags chat
// @Accept json
// @Success 202
// @Failure 400 {object} utils.Error
// @Router /chat/:id [post]
func (h *Handler) SendMessage(c echo.Context) error {

	/**

	In this method, the user is replying to an existing topic.

	*/

	var err error

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
	topicId, err := uuid.FromString(topicIdStr)
	if err != nil {
		err := ErrInvalidTopicId(topicIdStr)
		return NewErrResponse(c, &err)
	}
	topicKey := model.NewTopicKey(topicId)

	// todo verify that user has permission to post on topic

	err = h.chatStore.SendMessage(userKey, authUser.Username, topicKey, req.Message)
	if err != nil {
		return NewErrResponse(c, err)
	}

	return c.NoContent(http.StatusAccepted)

}
