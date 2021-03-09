package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestParallelListener(t *testing.T) {

	evtCalledMap := map[string]bool{}
	var expectedEventIds []string
	mu := sync.Mutex{}
	var listeners []Listener
	listenerCount := 50
	batchSize := 20

	for i := 0; i < listenerCount; i++ {
		events := test.NewMockEvents()
		for j := 0; j < batchSize; j++ {
			evtId := "evt-" + strconv.Itoa(i*batchSize+j)
			events = append(events, test.NewMockEvent(evtId))
			evtCalledMap[evtId] = false
			expectedEventIds = append(expectedEventIds, evtId)
		}
		listeners = append(listeners, NewStaticListener(events))
	}

	p := NewParallelListener(listeners)

	if !assert.NoError(t, p.Initialize(context.TODO(), "parallel", []string{"hello"})) {
		return
	}

	var calls [][]eventsource.Event
	if !assert.NoError(t, p.Listen(context.TODO(), func(ctx context.Context, events []eventsource.Event) error {
		calls = append(calls, events)
		mu.Lock()
		for _, event := range events {
			evtCalledMap[event.GetEventID()] = true
		}
		mu.Unlock()
		return nil
	})) {
		return
	}

	for _, value := range expectedEventIds {
		assert.True(t, evtCalledMap[value], value)
	}
	assert.Equal(t, len(evtCalledMap), len(expectedEventIds))
	assert.Equal(t, len(evtCalledMap), batchSize*listenerCount)

}
