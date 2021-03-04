package publish

import (
	"context"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventstore"
	"time"
)

type PublishEventStore struct {
	eventPublisher eventbus.EventPublisher
	eventStore     eventstore.EventStore
}

func NewPublishEventStore(eventStore eventstore.EventStore, eventPublisher eventbus.EventPublisher) *PublishEventStore {
	return &PublishEventStore{
		eventPublisher: eventPublisher,
		eventStore:     eventStore,
	}
}

func (p PublishEventStore) Load(ctx context.Context, streamKey eventstore.StreamKey) ([]*eventstore.StreamEvent, error) {
	return p.eventStore.Load(ctx, streamKey)
}

func (p PublishEventStore) Save(ctx context.Context, streamKey eventstore.StreamKey, expectedRevision int, events []*eventstore.StreamEvent) error {
	if err := p.eventStore.Save(ctx, streamKey, expectedRevision, events); err != nil {
		return err
	}
	return p.eventPublisher.PublishEvents(ctx, events)
}

func (p PublishEventStore) ReplayEventsByType(ctx context.Context, eventTypes []string, timestamp time.Time, replayFunc func(events []*eventstore.StreamEvent) error, options ...eventstore.ReplayEventsByTypeOptions) error {
	return p.eventStore.ReplayEventsByType(ctx, eventTypes, timestamp, replayFunc)
}

var _ eventstore.EventStore = &PublishEventStore{}
