package service

import (
	"context"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/mock"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/stretchr/testify/suite"
	"testing"
)

type serviceTestSuite struct {
	suite.Suite
	Service     *ChatService
	AmqpChannel *mock.AmqpChannel
	AmqpClient  *mock.AmqpClient
	ChatStore   *mock.ChatStore
}

func (s *serviceTestSuite) SetupTest() {
	s.Service = &ChatService{}
	s.AmqpChannel = &mock.AmqpChannel{
		CloseFunc: func() error { return nil },
		PublishFunc: func(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, publishing amqp.Message) error {
			return nil
		},
	}
	s.AmqpClient = &mock.AmqpClient{
		GetChannelFunc: func() (amqp.Channel, error) { return s.AmqpChannel, nil },
	}
	s.Service.amqpClient = s.AmqpClient

	s.ChatStore = &mock.ChatStore{
		SaveMessageFunc: func(ctx context.Context, message *chat.Message) error {
			return nil
		},
	}
	s.Service.chatStore = s.ChatStore
}

func TestMessageService(t *testing.T) {
	suite.Run(t, &serviceTestSuite{})
}
