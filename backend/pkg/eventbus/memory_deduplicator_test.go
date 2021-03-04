package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryDeduplicatorWithSmallBufferSize(t *testing.T) {

	d := NewMemoryDeduplicator(1)

	var calls []*eventstore.StreamEvent

	if !assert.NoError(t, d.Deduplicate(context.TODO(), evts(
		evt("t1", "1"),
		evt("t1", "2"),
		evt("t1", "1"),
	), func(evt *eventstore.StreamEvent) error {
		calls = append(calls, evt)
		return nil
	})) {
		return
	}

	assert.Len(t, calls, 3)
	assert.Equal(t, "1", calls[0].EventID)
	assert.Equal(t, "2", calls[1].EventID)
	assert.Equal(t, "1", calls[2].EventID)

}

func TestMemoryDeduplicator(t *testing.T) {

	d := NewMemoryDeduplicator(10)

	var calls []*eventstore.StreamEvent

	if !assert.NoError(t, d.Deduplicate(context.TODO(), evts(
		evt("t1", "1"),
		evt("t1", "2"),
		evt("t1", "1"),
		evt("t1", "3"),
	), func(evt *eventstore.StreamEvent) error {
		calls = append(calls, evt)
		return nil
	})) {
		return
	}

	assert.Len(t, calls, 3)
	assert.Equal(t, "1", calls[0].EventID)
	assert.Equal(t, "2", calls[1].EventID)
	assert.Equal(t, "3", calls[2].EventID)

}
