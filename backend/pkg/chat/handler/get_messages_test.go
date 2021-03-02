package handler

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"time"
)

func (s *chatHandlerSuite) TestGetMessagesShouldFailIfChannelKeyNotPresent() {
	s.NewContext(http.MethodGet, "/api/v1/chat/messages?take=bla", nil)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		test.HasStatusCode(http.StatusBadRequest),
		test.HasCode("ErrQueryParamRequired"),
		test.HasMessage(`query parameter 'channel' is required`))
}

func (s *chatHandlerSuite) TestGetMessagesShouldFailIfTakeQueryParamNotValid() {
	s.NewContext(http.MethodGet, "/api/v1/chat/messages?channel=abc&take=bla", nil)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		test.HasStatusCode(http.StatusBadRequest),
		test.HasCode("ErrInvalidTakeQueryParam"),
		test.HasMessage(`query parameter 'take' is invalid`))
}

func (s *chatHandlerSuite) TestGetMessagesShouldFailIfBeforeQueryParamNotValid() {
	s.NewContext(http.MethodGet, "/api/v1/chat/messages?channel=abc&before=bla", nil)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		test.HasStatusCode(http.StatusBadRequest),
		test.HasCode("ErrInvalidBeforeQueryParam"),
		test.HasMessage(`query parameter 'before' is invalid`))
}

func (s *chatHandlerSuite) TestGetMessages() {
	s.NewContext(http.MethodGet, "/api/v1/chat/messages?channel=abc&before=0", nil)
	s.ChatService.GetMessagesFunc = func(ctx context.Context, channel keys.ChannelKey, before time.Time, take int) (*chat.GetMessagesResponse, error) {
		return &chat.GetMessagesResponse{
			Messages: chat.Messages{
				Items: []chat.Message{},
			},
			HasMore: false,
		}, nil
	}

	s.ServeHTTP()

	if !s.AssertOK() {
		return
	}
	response := GetMessagesResponse{}
	s.ReadResponse(&response)
	assert.NotNil(s.T(), response.Messages)
	assert.Len(s.T(), response.Messages, 0)
}
