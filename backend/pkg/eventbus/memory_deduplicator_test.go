package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryDeduplicatorWithSmallBufferSize(t *testing.T) {

	d := NewMemoryDeduplicator(1)

	var calls []eventsource.Event

	if !assert.NoError(t, d.Deduplicate(context.TODO(), test.NewMockEvents(
		test.NewMockEvent("1"),
		test.NewMockEvent("2"),
		test.NewMockEvent("1"),
	), func(evt eventsource.Event) error {
		calls = append(calls, evt)
		return nil
	})) {
		return
	}

	assert.Len(t, calls, 3)
	assert.Equal(t, "1", calls[0].GetEventID())
	assert.Equal(t, "2", calls[1].GetEventID())
	assert.Equal(t, "1", calls[2].GetEventID())

}

func TestMemoryDeduplicator(t *testing.T) {

	d := NewMemoryDeduplicator(10)

	var calls []eventsource.Event

	if !assert.NoError(t, d.Deduplicate(context.TODO(), test.NewMockEvents(
		test.NewMockEvent("1"),
		test.NewMockEvent("2"),
		test.NewMockEvent("1"),
		test.NewMockEvent("3"),
	), func(evt eventsource.Event) error {
		calls = append(calls, evt)
		return nil
	})) {
		return
	}

	assert.Len(t, calls, 3)
	assert.Equal(t, "1", calls[0].GetEventID())
	assert.Equal(t, "2", calls[1].GetEventID())
	assert.Equal(t, "3", calls[2].GetEventID())

}
