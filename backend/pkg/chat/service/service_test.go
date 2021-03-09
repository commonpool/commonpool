package service

import (
	"context"
	mock2 "github.com/commonpool/backend/mock/amqp_channel"
	mock3 "github.com/commonpool/backend/mock/amqp_client"
	mock4 "github.com/commonpool/backend/mock/chat_store"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/stretchr/testify/suite"
	"testing"
)

type serviceTestSuite struct {
	suite.Suite
	Service     *ChatService
	AmqpChannel *mock2.AmqpChannel
	AmqpClient  *mock3.AmqpClient
	ChatStore   *mock4.ChatStore
}

func (s *serviceTestSuite) SetupTest() {
	s.Service = &ChatService{}
	s.AmqpChannel = &mock2.AmqpChannel{
		CloseFunc: func() error { return nil },
		PublishFunc: func(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, publishing mq.Message) error {
			return nil
		},
	}
	s.AmqpClient = &mock3.AmqpClient{
		GetChannelFunc: func() (mq.Channel, error) { return s.AmqpChannel, nil },
	}
	s.Service.amqpClient = s.AmqpClient

	s.ChatStore = &mock4.ChatStore{
		SaveMessageFunc: func(ctx context.Context, message *chat.Message) error {
			return nil
		},
	}
	s.Service.chatStore = s.ChatStore
}

func TestMessageService(t *testing.T) {
	suite.Run(t, &serviceTestSuite{})
}
