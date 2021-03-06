package store

import (
	"context"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
)

type EventSourcedGroupRepository struct {
	eventStore eventstore.EventStore
}

func NewEventSourcedGroupRepository(eventStore eventstore.EventStore) *EventSourcedGroupRepository {
	return &EventSourcedGroupRepository{
		eventStore: eventStore,
	}
}

func (e EventSourcedGroupRepository) Load(ctx context.Context, groupKey keys.GroupKey) (*domain.Group, error) {
	events, err := e.eventStore.Load(ctx, eventstore.NewStreamKey("offer", groupKey.String()))
	if err != nil {
		return nil, err
	}
	return domain.NewFromEvents(groupKey, events), nil
}

func (e *EventSourcedGroupRepository) Save(ctx context.Context, group *domain.Group) error {
	streamKey := eventstore.NewStreamKey("group", group.GetKey().String())
	if err := e.eventStore.Save(ctx, streamKey, group.GetVersion(), group.GetChanges()); err != nil {
		return err
	}
	group.MarkAsCommitted()
	return nil
}

var _ domain.GroupRepository = &EventSourcedGroupRepository{}
