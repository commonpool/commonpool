package domain

import "github.com/commonpool/backend/pkg/eventsource"

type OfferCompletedPayload struct {
}

type OfferCompleted struct {
	eventsource.EventEnvelope
	OfferCompletedPayload `json:"payload"`
}

func NewOfferCompleted() *OfferCompleted {
	return &OfferCompleted{
		eventsource.NewEventEnvelope(OfferCompletedEvent, 1),
		OfferCompletedPayload{},
	}
}

var _ eventsource.Event = &OfferCompleted{}
