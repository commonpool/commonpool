package keys

import (
	"go.uber.org/zap/zapcore"
)

type ChannelSubscriptionKey struct {
	ChannelKey ChannelKey
	UserKey    UserKey
}

func (c ChannelSubscriptionKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	err := encoder.AddObject("channel", c.ChannelKey)
	if err != nil {
		return err
	}
	return encoder.AddObject("user", c.UserKey)
}

var _ zapcore.ObjectMarshaler = ChannelSubscriptionKey{}

func NewChannelSubscriptionKey(channelKey ChannelKey, userKey UserKey) ChannelSubscriptionKey {
	return ChannelSubscriptionKey{
		ChannelKey: channelKey,
		UserKey:    userKey,
	}
}
