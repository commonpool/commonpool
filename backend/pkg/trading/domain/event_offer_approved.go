package domain

import "github.com/commonpool/backend/pkg/eventsource"

type OfferApprovedPayload struct {
}

type OfferApproved struct {
	eventsource.EventEnvelope
	OfferApprovedPayload `json:"payload"`
}

func NewOfferApproved() *OfferApproved {
	return &OfferApproved{
		eventsource.NewEventEnvelope(OfferApprovedEvent, 1),
		OfferApprovedPayload{},
	}
}

var _ eventsource.Event = &OfferApproved{}
