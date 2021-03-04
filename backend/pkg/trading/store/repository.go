package store

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	uuid "github.com/satori/go.uuid"
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

	events, err := e.eventStore.Load(ctx, eventstore.NewStreamKey("offer", offerKey.String()))
	if err != nil {
		return nil, err
	}

	var decodedEvents []domain.Event

	for _, event := range events {
		var decoded domain.Event
		switch event.EventType {
		case string(domain.OfferSubmittedEvent):
			decoded = &domain.OfferSubmitted{}
		case string(domain.OfferItemApprovedEvent):
			decoded = &domain.OfferItemApproved{}
		case string(domain.OfferDeclinedEvent):
			decoded = &domain.OfferDeclined{}
		case string(domain.OfferApprovedEvent):
			decoded = &domain.OfferApproved{}
		case string(domain.OfferCompletedEvent):
			decoded = &domain.OfferCompleted{}
		case string(domain.ServiceGivenNotifiedEvent):
			decoded = &domain.ServiceGivenNotified{}
		case string(domain.ServiceReceivedNotifiedEvent):
			decoded = &domain.ServiceReceivedNotified{}
		case string(domain.ResourceTransferGivenNotifiedEvent):
			decoded = &domain.ResourceTransferGivenNotified{}
		case string(domain.ResourceTransferReceivedNotifiedEvent):
			decoded = &domain.ResourceTransferReceivedNotified{}
		case string(domain.ResourceBorrowedNotifiedEvent):
			decoded = &domain.ResourceBorrowedNotified{}
		case string(domain.ResourceLentNotifiedEvent):
			decoded = &domain.ResourceLentNotified{}
		case string(domain.BorrowerReturnedResourceEvent):
			decoded = &domain.BorrowerReturnedResourceNotified{}
		case string(domain.LenderReceivedBackResourceEvent):
			decoded = &domain.LenderReceivedBackResourceNotified{}
		}

		err := json.Unmarshal([]byte(event.Payload), decoded)
		if err != nil {
			return nil, err
		}
		decodedEvents = append(decodedEvents, decoded)

	}

	return domain.NewFromEvents(offerKey, decodedEvents), nil

}

func (e EventSourcedOfferRepository) Save(ctx context.Context, offer *domain.Offer) error {

	var streamEvents []*eventstore.StreamEvent

	streamKey := eventstore.NewStreamKey("offer", offer.GetKey().String())
	for _, change := range offer.GetChanges() {

		evtJson, err := json.Marshal(change)
		if err != nil {
			return err
		}

		streamEvent := eventstore.NewStreamEvent(
			streamKey,
			eventstore.NewStreamEventKey(string(change.GetType()), uuid.NewV4().String()),
			string(evtJson),
			eventstore.NewStreamEventOptions{
				Version: change.GetVersion(),
			},
		)
		streamEvents = append(streamEvents, streamEvent)
	}

	if err := e.eventStore.Save(ctx, streamKey, offer.GetVersion(), streamEvents); err != nil {
		return err
	}

	offer.MarkAsCommitted()

	return nil

}

var _ domain.OfferRepository = &EventSourcedOfferRepository{}
