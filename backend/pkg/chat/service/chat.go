package service

import (
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/commonpool/backend/pkg/user"
)

type ChatService struct {
	userStore  user.Store
	amqpClient mq.Client
	chatStore  chat.Store
}

func NewChatService(us user.Store, mq mq.Client, cs chat.Store) *ChatService {
	return &ChatService{
		userStore:  us,
		amqpClient: mq,
		chatStore:  cs,
	}
}

var _ chat.Service = &ChatService{}

// getChannelBindingHeaders returns the binding headers to link the websocket messages exchange and a given user exchange
// The user will receive messages on his exchange if the message has a "channel_id" = "subscribed channel id" header
func (c ChatService) getChannelBindingHeaders(channelSubscriptionKey keys.ChannelSubscriptionKey) map[string]interface{} {
	return map[string]interface{}{
		"event_type": "chat.message",
		"channel_id": channelSubscriptionKey.ChannelKey.String(),
		"x-match":    "all",
	}
}
