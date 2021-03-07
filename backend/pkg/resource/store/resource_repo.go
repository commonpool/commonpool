package store

import (
	"context"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
)

type EventSourcedResourceRepository struct {
	eventStore eventstore.EventStore
}

func NewEventSourcedResourceRepository(eventStore eventstore.EventStore) *EventSourcedResourceRepository {
	return &EventSourcedResourceRepository{
		eventStore: eventStore,
	}
}

func (e EventSourcedResourceRepository) Load(ctx context.Context, resourceKey keys.ResourceKey) (*domain.Resource, error) {
	events, err := e.eventStore.Load(ctx, eventstore.NewStreamKey("resource", resourceKey.String()))
	if err != nil {
		return nil, err
	}
	return domain.NewFromEvents(resourceKey, events), nil
}

func (e *EventSourcedResourceRepository) Save(ctx context.Context, resource *domain.Resource) error {
	streamKey := eventstore.NewStreamKey("resource", resource.GetKey().String())
	if _, err := e.eventStore.Save(ctx, streamKey, resource.GetVersion(), resource.GetChanges()); err != nil {
		return err
	}
	resource.MarkAsCommitted()
	return nil
}

var _ domain.ResourceRepository = &EventSourcedResourceRepository{}
