package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeduplicateListener(t *testing.T) {

	memoryDeduplicator := NewMemoryDeduplicator(10)

	l1 := NewStaticListener(
		test.NewMockEvents(
			test.NewMockEvent("1"),
			test.NewMockEvent("2"),
			test.NewMockEvent("1"),
		),
	)

	l2 := NewStaticListener(
		test.NewMockEvents(
			test.NewMockEvent("1"),
			test.NewMockEvent("2"),
			test.NewMockEvent("3"),
			test.NewMockEvent("2"),
			test.NewMockEvent("1"),
			test.NewMockEvent("1"),
		),
	)

	l3 := NewStaticListener(
		test.NewMockEvents(
			test.NewMockEvent("2"),
			test.NewMockEvent("1"),
		),
	)

	ls := NewSequenceListener([]Listener{l1, l2, l3})

	ds := NewDeduplicateListener(memoryDeduplicator, ls)

	if !assert.NoError(t, ds.Initialize(context.TODO(), "test", []string{"typ1"})) {
		return
	}

	var calls [][]eventsource.Event
	if !assert.NoError(t, ds.Listen(context.TODO(), func(events []eventsource.Event) error {
		calls = append(calls, events)
		return nil
	})) {
		return
	}

	assert.Len(t, calls, 3)
	assert.Len(t, calls[0], 1)
	assert.Len(t, calls[1], 1)
	assert.Len(t, calls[1], 1)
	assert.Equal(t, calls[0][0].GetEventID(), "1")
	assert.Equal(t, calls[1][0].GetEventID(), "2")
	assert.Equal(t, calls[2][0].GetEventID(), "3")

}
