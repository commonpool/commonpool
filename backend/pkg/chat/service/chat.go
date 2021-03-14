package service

import (
	"github.com/commonpool/backend/pkg/auth/queries"
	store2 "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/mq"
)

type ChatService struct {
	amqpClient     mq.Client
	chatStore      store2.Store
	getUsersByKeys *queries.GetUsersByKeys
}

func NewChatService(mq mq.Client, cs store2.Store, getUsersByKeys *queries.GetUsersByKeys) *ChatService {
	return &ChatService{
		amqpClient:     mq,
		chatStore:      cs,
		getUsersByKeys: getUsersByKeys,
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
