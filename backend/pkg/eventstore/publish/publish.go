package publish

import (
	"context"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
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

func (p PublishEventStore) Load(ctx context.Context, streamKey eventstore.StreamKey) ([]eventsource.Event, error) {
	return p.eventStore.Load(ctx, streamKey)
}

func (p PublishEventStore) Save(ctx context.Context, streamKey eventstore.StreamKey, expectedRevision int, events []eventsource.Event) ([]eventsource.Event, error) {
	publishedEvents, err := p.eventStore.Save(ctx, streamKey, expectedRevision, events)
	if err != nil {
		return nil, err
	}
	err = p.eventPublisher.PublishEvents(ctx, publishedEvents)
	return publishedEvents, err
}

func (p PublishEventStore) ReplayEventsByType(ctx context.Context, eventTypes []string, timestamp time.Time, replayFunc func(events []eventsource.Event) error, options ...eventstore.ReplayEventsByTypeOptions) error {
	return p.eventStore.ReplayEventsByType(ctx, eventTypes, timestamp, replayFunc)
}

var _ eventstore.EventStore = &PublishEventStore{}
