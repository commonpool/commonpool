package service

import (
	"github.com/commonpool/backend/pkg/chat"
	chatmodel "github.com/commonpool/backend/pkg/chat/model"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/commonpool/backend/pkg/resource"
	"github.com/commonpool/backend/pkg/user"
)

type ChatService struct {
	us         user.Store
	gs         group2.Store
	amqpClient mq.Client
	rs         resource.Store
	chatStore  chat.Store
}

func NewChatService(us user.Store, gs group2.Store, rs resource.Store, mq mq.Client, cs chat.Store) *ChatService {
	return &ChatService{
		us:         us,
		gs:         gs,
		amqpClient: mq,
		rs:         rs,
		chatStore:  cs,
	}
}

var _ chat.Service = &ChatService{}

// getChannelBindingHeaders returns the binding headers to link the websocket messages exchange and a given user exchange
// The user will receive messages on his exchange if the message has a "channel_id" = "subscribed channel id" header
func (c ChatService) getChannelBindingHeaders(channelSubscriptionKey chatmodel.ChannelSubscriptionKey) map[string]interface{} {
	return map[string]interface{}{
		"event_type": "chat.message",
		"channel_id": channelSubscriptionKey.ChannelKey.String(),
		"x-match":    "all",
	}
}
