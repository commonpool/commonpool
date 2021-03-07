package store

import (
	"context"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type EventSourcedOfferRepository struct {
	eventStore eventstore.EventStore
}

func NewEventSourcedOfferRepository(eventStore eventstore.EventStore) *EventSourcedOfferRepository {
	return &EventSourcedOfferRepository{
		eventStore: eventStore,
	}
}

func (e EventSourcedOfferRepository) Load(ctx context.Context, offerKey keys.OfferKey) (*domain.Offer, error) {
	events, err := e.eventStore.Load(ctx, keys.NewStreamKey("offer", offerKey.String()))
	if err != nil {
		return nil, err
	}
	return domain.NewFromEvents(offerKey, events), nil
}

func (e EventSourcedOfferRepository) Save(ctx context.Context, offer *domain.Offer) error {
	streamKey := keys.NewStreamKey("offer", offer.GetKey().String())
	if _, err := e.eventStore.Save(ctx, streamKey, offer.GetSequenceNo(), offer.GetChanges()); err != nil {
		return err
	}
	offer.MarkAsCommitted()
	return nil
}

var _ domain.OfferRepository = &EventSourcedOfferRepository{}
