package service

import (
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	res "github.com/commonpool/backend/resource"
)

type ChatService struct {
	us auth.Store
	gs group.Store
	mq amqp.Client
	rs res.Store
	cs chat.Store
}

func NewChatService(us auth.Store, gs group.Store, rs res.Store, mq amqp.Client, cs chat.Store) *ChatService {
	return &ChatService{
		us: us,
		gs: gs,
		mq: mq,
		rs: rs,
		cs: cs,
	}
}

var _ chat.Service = &ChatService{}

// getChannelBindingHeaders returns the binding headers to link the websocket messages exchange and a given user exchange
// The user will receive messages on his exchange if the message has a "channel_id" = "subscribed channel id" header
func (c ChatService) getChannelBindingHeaders(channelSubscriptionKey model.ChannelSubscriptionKey) map[string]interface{} {
	return map[string]interface{}{
		"event_type": "chat.message",
		"channel_id": channelSubscriptionKey.ChannelKey.String(),
		"x-match":    "all",
	}
}
