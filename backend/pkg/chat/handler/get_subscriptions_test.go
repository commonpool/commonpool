package handler

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func (s *chatHandlerSuite) TestGetSubscriptions() {
	s.NewContext(http.MethodGet, "/api/v1/chat/subscriptions", nil)
	s.ChatService.GetSubscriptionsForUserFunc = func(ctx context.Context, take int, skip int) (*chat.ChannelSubscriptions, error) {
		return &chat.ChannelSubscriptions{
			Items: []*chat.ChannelSubscription{},
		}, nil
	}
	s.ServeHTTP()
	if !s.AssertOK() {
		return
	}
	response := &GetLatestSubscriptionsResponse{}
	s.ReadResponse(response)
	assert.Len(s.T(), response.Subscriptions, 0)
}

func (s *chatHandlerSuite) TestGetSubscriptionsShouldFailIfSkipInvalidInt() {
	s.NewContext(http.MethodGet, "/api/v1/chat/subscriptions?skip=bla", nil)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		test.HasStatusCode(http.StatusBadRequest),
		test.HasCode("ErrInvalidSkipQueryParam"),
		test.HasMessage(`query parameter 'skip' is invalid`))
}

func (s *chatHandlerSuite) TestGetSubscriptionsShouldFailIfTakeInvalidInt() {
	s.NewContext(http.MethodGet, "/api/v1/chat/subscriptions?take=bla", nil)
	s.ServeHTTP()
	s.AssertBadRequest()
	s.AssertErrorResponse(
		test.HasStatusCode(http.StatusBadRequest),
		test.HasCode("ErrInvalidTakeQueryParam"),
		test.HasMessage(`query parameter 'take' is invalid`))
}
