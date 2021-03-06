package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSequenceListener(t *testing.T) {

	l1 := NewStaticListener(test.NewMockEvents(evt("typ1", "1")))
	l2 := NewStaticListener(test.NewMockEvents(evt("typ1", "2")))
	l3 := NewStaticListener(test.NewMockEvents(evt("typ1", "3"), evt("typ1", "4")))

	l := NewSequenceListener([]Listener{l1, l2, l3})

	if !assert.NoError(t, l.Initialize(context.TODO(), "name", []string{"bla"})) {
		return
	}

	var calls [][]*eventstore.StreamEvent
	if !assert.NoError(t, l.Listen(context.TODO(), func(events []*eventstore.StreamEvent) error {
		calls = append(calls, events)
		return nil
	})) {
		return
	}

	assert.Len(t, calls, 3)
	assert.Len(t, calls[0], 1)
	assert.Len(t, calls[1], 1)
	assert.Len(t, calls[2], 2)
	assert.Equal(t, "1", calls[0][0].EventID)
	assert.Equal(t, "2", calls[1][0].EventID)
	assert.Equal(t, "3", calls[2][0].EventID)
	assert.Equal(t, "4", calls[2][1].EventID)

}
