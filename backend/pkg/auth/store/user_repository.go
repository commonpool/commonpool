package store

import (
	"context"
	userdomain "github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/keys"
)

type EventSourcedUserRepository struct {
	eventStore eventstore.EventStore
}

func NewEventSourcedUserRepository(eventStore eventstore.EventStore) *EventSourcedUserRepository {
	return &EventSourcedUserRepository{
		eventStore: eventStore,
	}
}

var _ userdomain.UserRepository = &EventSourcedUserRepository{}

func (e EventSourcedUserRepository) Load(ctx context.Context, userKey keys.UserKey) (*userdomain.User, error) {
	events, err := e.eventStore.Load(ctx, eventstore.NewStreamKey("user", userKey.String()))
	if err != nil {
		return nil, err
	}
	return userdomain.NewFromEvents(userKey, events), nil
}

func (e EventSourcedUserRepository) Save(ctx context.Context, user *userdomain.User) error {
	streamKey := eventstore.NewStreamKey("user", user.GetKey().String())
	if err := e.eventStore.Save(ctx, streamKey, user.GetVersion(), user.GetChanges()); err != nil {
		return err
	}
	user.MarkAsCommitted()
	return nil

}
