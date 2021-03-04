package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/eventstore"
	"time"
)

type ReplayListener struct {
	eventStore          eventstore.EventStore
	name                string
	eventTypes          []string
	getCurrentTimestamp func() time.Time
	lastIds             []string
}

func NewReplayListener(
	eventStore eventstore.EventStore,
	getCurrentTimestamp func() time.Time) *ReplayListener {
	return &ReplayListener{
		eventStore:          eventStore,
		getCurrentTimestamp: getCurrentTimestamp,
	}
}

func (s *ReplayListener) Initialize(ctx context.Context, name string, eventTypes []string) error {
	s.name = name
	s.eventTypes = eventTypes
	return nil
}

func (s *ReplayListener) Listen(ctx context.Context, listenerFunc ListenerFunc) error {
	return s.eventStore.ReplayEventsByType(ctx, s.eventTypes, s.getCurrentTimestamp(), listenerFunc)
}
