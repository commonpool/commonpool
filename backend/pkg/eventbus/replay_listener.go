package eventbus

import (
	"context"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/eventsource"
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
	l := logging.WithContext(ctx).Named("ReplayListener " + s.name)
	return s.eventStore.ReplayEventsByType(ctx, s.eventTypes, s.getCurrentTimestamp(), func(events []eventsource.Event) error {
		l.Debug("received events from event store")
		return listenerFunc(ctx, events)
	})
}
