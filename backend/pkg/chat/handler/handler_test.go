package handler

import (
	"github.com/commonpool/backend/mock"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type chatHandlerSuite struct {
	test.HandlerSuite
	ChatService *mock.ChatService
	Handler     *Handler
}

func TestChatHandler(t *testing.T) {
	suite.Run(t, &chatHandlerSuite{})
}

func (s *chatHandlerSuite) SetupTest() {
	s.HandlerSuite.SetupTest()
	s.ChatService = &mock.ChatService{}

	s.Handler = &Handler{
		chatService: s.ChatService,
		auth:        s.Authenticator,
	}
	group := s.Echo.Group("/api/v1")
	s.Handler.Register(group)
}
