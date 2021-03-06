package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSequenceListener(t *testing.T) {

	l1 := NewStaticListener(test.NewMockEvents(test.NewMockEvent("1")))
	l2 := NewStaticListener(test.NewMockEvents(test.NewMockEvent("2")))
	l3 := NewStaticListener(test.NewMockEvents(test.NewMockEvent("3"), test.NewMockEvent("4")))

	l := NewSequenceListener([]Listener{l1, l2, l3})

	if !assert.NoError(t, l.Initialize(context.TODO(), "name", []string{"bla"})) {
		return
	}

	var calls [][]eventsource.Event
	if !assert.NoError(t, l.Listen(context.TODO(), func(events []eventsource.Event) error {
		calls = append(calls, events)
		return nil
	})) {
		return
	}

	assert.Len(t, calls, 3)
	assert.Len(t, calls[0], 1)
	assert.Len(t, calls[1], 1)
	assert.Len(t, calls[2], 2)
	assert.Equal(t, "1", calls[0][0].GetEventID())
	assert.Equal(t, "2", calls[1][0].GetEventID())
	assert.Equal(t, "3", calls[2][0].GetEventID())
	assert.Equal(t, "4", calls[2][1].GetEventID())

}
