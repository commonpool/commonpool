package test

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

const (
	MockAggregateType = "mock-aggregate"
	MockEventType     = "mock-event"
	MockAggregateID   = "mock-id"
)

var (
	MockStreamKey = keys.NewStreamKey(MockAggregateType, MockAggregateID)
)

type MockEvent struct {
	eventsource.EventEnvelope
}

type MockEventOptions struct {
	EventTime *time.Time
}

func MockEventTime(time time.Time) MockEventOptions {
	return MockEventOptions{
		EventTime: &time,
	}
}

func NewMockEvent(id string, options ...MockEventOptions) MockEvent {
	evt := MockEvent{
		eventsource.NewEventEnvelope(MockEventType, 1),
	}
	evt.EventID = id
	evt.AggregateType = MockAggregateType
	evt.AggregateID = MockAggregateID
	if len(options) > 0 {
		option := options[0]
		if option.EventTime != nil {
			evt.EventTime = *option.EventTime
		}
	}
	return evt
}

func NewMockEvents(events ...eventsource.Event) []eventsource.Event {
	return events
}

func RegisterMockEvents(mapper *eventsource.EventMapper) error {
	return mapper.RegisterMapper(MockEventType, func(eventType string, bytes []byte) (eventsource.Event, error) {
		var dest MockEvent
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	})
}
