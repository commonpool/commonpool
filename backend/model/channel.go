package model

import "go.uber.org/zap/zapcore"

type ChannelKey struct {
	ID string
}

func (tk ChannelKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("channel_id", tk.String())
	return nil
}

var _ zapcore.ObjectMarshaler = &ChannelKey{}

func (tk *ChannelKey) String() string {
	return tk.ID
}

func NewConversationKey(key string) ChannelKey {
	return ChannelKey{
		ID: key,
	}
}
