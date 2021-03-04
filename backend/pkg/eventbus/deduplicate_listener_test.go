package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeduplicateListener(t *testing.T) {

	memoryDeduplicator := NewMemoryDeduplicator(10)

	l1 := NewStaticListener(
		evts(
			evt("typ1", "1"),
			evt("typ1", "2"),
			evt("typ1", "1"),
		),
	)

	l2 := NewStaticListener(
		evts(
			evt("typ1", "1"),
			evt("typ1", "2"),
			evt("typ1", "3"),
			evt("typ1", "2"),
			evt("typ1", "1"),
			evt("typ1", "1"),
		),
	)

	l3 := NewStaticListener(
		evts(
			evt("typ1", "2"),
			evt("typ1", "1"),
		),
	)

	ls := NewSequenceListener([]Listener{l1, l2, l3})

	ds := NewDeduplicateListener(memoryDeduplicator, ls)

	if !assert.NoError(t, ds.Initialize(context.TODO(), "test", []string{"typ1"})) {
		return
	}

	var calls [][]*eventstore.StreamEvent
	if !assert.NoError(t, ds.Listen(context.TODO(), func(events []*eventstore.StreamEvent) error {
		calls = append(calls, events)
		return nil
	})) {
		return
	}

	assert.Len(t, calls, 3)
	assert.Len(t, calls[0], 1)
	assert.Len(t, calls[1], 1)
	assert.Len(t, calls[1], 1)
	assert.Equal(t, calls[0][0].EventID, "1")
	assert.Equal(t, calls[1][0].EventID, "2")
	assert.Equal(t, calls[2][0].EventID, "3")

}
