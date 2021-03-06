package service

import (
	"github.com/commonpool/backend/pkg/auth/store"
	store2 "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/mq"
)

type ChatService struct {
	userStore  store.Store
	amqpClient mq.Client
	chatStore  store2.Store
}

func NewChatService(us store.Store, mq mq.Client, cs store2.Store) *ChatService {
	return &ChatService{
		userStore:  us,
		amqpClient: mq,
		chatStore:  cs,
	}
}

var _ Service = &ChatService{}

// getChannelBindingHeaders returns the binding headers to link the websocket messages exchange and a given user exchange
// The user will receive messages on his exchange if the message has a "channel_id" = "subscribed channel id" header
func (c ChatService) getChannelBindingHeaders(channelSubscriptionKey keys.ChannelSubscriptionKey) map[string]interface{} {
	return map[string]interface{}{
		"event_type": "chat.message",
		"channel_id": channelSubscriptionKey.ChannelKey.String(),
		"x-match":    "all",
	}
}
