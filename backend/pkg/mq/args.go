package mq

import (
	"go.uber.org/zap/zapcore"
)

type Args map[string]interface{}

func NewArgs() Args {
	return Args{}
}

func (a Args) With(key ArgKey, value string) Args {
	a[ChannelIdArg] = value
	return a
}

func (a Args) WithEventType(eventType EventType) Args {
	a[EventTypeArg] = eventType
	return a
}

func (a Args) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	if a == nil {
		return nil
	}

	for s, i := range a {
		err := encoder.AddReflected(s, i)
		if err != nil {
			return err
		}
	}

	return nil
}

var _ zapcore.ObjectMarshaler = Args{}
