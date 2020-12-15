package handler

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	model2 "github.com/commonpool/backend/pkg/chat/handler/model"
	"github.com/commonpool/backend/pkg/chat/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"time"
)

func (s *chatHandlerSuite) TestGetMessagesShouldFailIfChannelKeyNotPresent() {
	s.NewContext(http.MethodGet, "/api/v1/chat/messages?take=bla", nil)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		HasStatusCode(http.StatusBadRequest),
		HasCode("ErrQueryParamRequired"),
		HasMessage(`query parameter 'channel' is required`))
}

func (s *chatHandlerSuite) TestGetMessagesShouldFailIfTakeQueryParamNotValid() {
	s.NewContext(http.MethodGet, "/api/v1/chat/messages?channel=abc&take=bla", nil)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		HasStatusCode(http.StatusBadRequest),
		HasCode("ErrInvalidTakeQueryParam"),
		HasMessage(`query parameter 'take' is invalid`))
}

func (s *chatHandlerSuite) TestGetMessagesShouldFailIfBeforeQueryParamNotValid() {
	s.NewContext(http.MethodGet, "/api/v1/chat/messages?channel=abc&before=bla", nil)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		HasStatusCode(http.StatusBadRequest),
		HasCode("ErrInvalidBeforeQueryParam"),
		HasMessage(`query parameter 'before' is invalid`))
}

func (s *chatHandlerSuite) TestGetMessages() {
	s.NewContext(http.MethodGet, "/api/v1/chat/messages?channel=abc&before=0", nil)
	s.ChatService.GetMessagesFunc = func(ctx context.Context, channel model.ChannelKey, before time.Time, take int) (*chat.GetMessagesResponse, error) {
		return &chat.GetMessagesResponse{
			Messages: model.Messages{
				Items: []model.Message{},
			},
			HasMore: false,
		}, nil
	}

	s.ServeHTTP()

	if !s.AssertOK() {
		return
	}
	response := model2.GetTopicMessagesResponse{}
	s.ReadResponse(&response)
	assert.NotNil(s.T(), response.Messages)
	assert.Len(s.T(), response.Messages, 0)
}
