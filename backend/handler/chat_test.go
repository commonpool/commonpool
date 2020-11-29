package handler

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/router"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRecentlyActiveSubscriptions(t *testing.T) {

	e := router.NewRouter()
	httpRequest := httptest.NewRequest(http.MethodGet, "", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpRequest, rec)

	h := Handler{
		amqp:           nil,
		resourceStore:  nil,
		authStore:      nil,
		authorization:  nil,
		chatStore:      nil,
		tradingStore:   nil,
		groupService:   nil,
		config:         config.AppConfig{},
		chatService:    nil,
		tradingService: nil,
	}

	h.GetRecentlyActiveSubscriptions()

}

// TestCreateResource
// Should be able to create a resource
func TestSendMessage(t *testing.T) {
	tearDown()
	setup()

	// user 1 creates a resource
	mockLoggedInAs(user1)
	res := createResource(t, "summary", "description", resource.Offer)

	// user 2 sends message about resource
	mockLoggedInAs(user2)
	inquireAboutResource(t, res.Resource.Id, "hello!")

	user2Threads := getLatestThreads(t, 0, 10).Subscriptions
	assert.Equal(t, 1, len(user2Threads))
	user2Messages := getThreadMessages(t, user2Threads[0].ChannelID).Messages
	assert.Equal(t, 1, len(user2Messages))

	js, _ := json.MarshalIndent(user2Messages, "", "   ")
	fmt.Println(string(js))

	// user 1 checks his messages
	mockLoggedInAs(user1)
	user1Threads := getLatestThreads(t, 0, 10).Subscriptions
	assert.Equal(t, 1, len(user1Threads))

	js, _ = json.MarshalIndent(user1Threads, "", "   ")
	fmt.Println("user 1 threads", string(js))

	user1Messages := getThreadMessages(t, user1Threads[0].ChannelID).Messages
	assert.Equal(t, 1, len(user1Messages))

	js, _ = json.MarshalIndent(user1Messages, "", "   ")
	fmt.Println("user 1 messages", string(js))

	// user 1 replies to user 2
	sendMessage(t, user1Threads[0].ChannelID, "hello back!")
	user1Threads = getLatestThreads(t, 0, 10).Subscriptions

	js, _ = json.MarshalIndent(user1Threads, "", "   ")
	fmt.Println("user 1 threads", string(js))

	assert.Equal(t, 1, len(user1Threads))
	user1Messages = getThreadMessages(t, user1Threads[0].ChannelID).Messages

	js, _ = json.MarshalIndent(user1Threads, "", "   ")
	fmt.Println("user 1 messages", string(js))
	assert.Equal(t, 2, len(user1Messages))

	js, _ = json.MarshalIndent(user1Messages, "", "   ")
	fmt.Println(string(js))

}

func newSendMessageRequest(js string, topicId string) (*httptest.ResponseRecorder, echo.Context) {
	_, _, rec, c := newRequest(echo.POST, fmt.Sprintf("/api/chat/%s", topicId), &js)
	c.SetParamNames("id")
	c.SetParamValues(topicId)
	return rec, c
}

func sendMessage(t *testing.T, topicId string, message string) {
	js := fmt.Sprintf(`{ "message" : "%s" }`, message)
	rec, c := newSendMessageRequest(js, topicId)
	assert.NoError(t, h.SendMessage(c))
	assert.Equal(t, http.StatusAccepted, rec.Code)
}

func newInquireAboutResourceRequest(js string, resourceId string) (*httptest.ResponseRecorder, echo.Context) {
	_, _, rec, c := newRequest(echo.POST, fmt.Sprintf("/api/resources/%s/inquire", resourceId), &js)
	c.SetParamNames("id")
	c.SetParamValues(resourceId)
	return rec, c
}

func inquireAboutResource(t *testing.T, resourceId string, message string) {
	js := fmt.Sprintf(`{ "message" : "%s" }`, message)
	rec, c := newInquireAboutResourceRequest(js, resourceId)
	assert.NoError(t, h.InquireAboutResource(c))
	assert.Equal(t, http.StatusAccepted, rec.Code)
}

func newGetLatestThreadsRequest(skip int, take int) (*httptest.ResponseRecorder, echo.Context) {
	_, _, rec, c := newRequest(echo.GET, fmt.Sprintf("/api/chat/threads?take=%d&skip=%d", take, skip), nil)
	return rec, c
}

func getLatestThreads(t *testing.T, skip int, take int) web.GetLatestSubscriptionsResponse {
	rec, c := newGetLatestThreadsRequest(skip, take)
	assert.NoError(t, h.GetRecentlyActiveSubscriptions(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	threads := web.GetLatestSubscriptionsResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &threads))
	return threads

}

func newGetMessagesRequest(thread string) (*httptest.ResponseRecorder, echo.Context) {
	_, _, rec, c := newRequest(echo.GET, "/api/chat/messages", nil)
	c.QueryParams().Set("topic", thread)
	return rec, c
}

func getThreadMessages(t *testing.T, threadId string) web.GetTopicMessagesResponse {
	rec, c := newGetMessagesRequest(threadId)
	assert.NoError(t, h.GetMessages(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	messages := web.GetTopicMessagesResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &messages))
	return messages

}
