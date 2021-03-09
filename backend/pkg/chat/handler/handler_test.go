package handler

import (
	mock2 "github.com/commonpool/backend/mock/chat_service"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type chatHandlerSuite struct {
	test.HandlerSuite
	ChatService *mock2.ChatService
	Handler     *Handler
}

func TestChatHandler(t *testing.T) {
	suite.Run(t, &chatHandlerSuite{})
}

func (s *chatHandlerSuite) SetupTest() {
	s.HandlerSuite.SetupTest()
	s.ChatService = &mock2.ChatService{}

	s.Handler = &Handler{
		chatService: s.ChatService,
		auth:        s.Authenticator,
	}
	group := s.Echo.Group("/api/v1")
	s.Handler.Register(group)
}
