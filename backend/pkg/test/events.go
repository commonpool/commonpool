package test

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/eventsource"
)

type MockEvent struct {
	eventsource.EventEnvelope
}

func NewMockEvent(id string) MockEvent {
	evt := MockEvent{
		eventsource.NewEventEnvelope("mock", 1),
	}
	evt.EventID = id
	evt.AggregateType = "mock"
	evt.AggregateID = "mock-id"
	return evt
}

func NewMockEvents(events ...eventsource.Event) []eventsource.Event {
	return events
}

func RegisterMockEvents(mapper *eventsource.EventMapper) error {
	return mapper.RegisterMapper("mock", func(eventType string, bytes []byte) (eventsource.Event, error) {
		var dest MockEvent
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	})
}
