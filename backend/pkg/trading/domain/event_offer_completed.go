package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type OfferCompletedPayload struct {
	OfferItems *OfferItems   `json:"offer_items"`
	GroupKey   keys.GroupKey `json:"group_key"`
}

type OfferCompleted struct {
	eventsource.EventEnvelope
	OfferCompletedPayload `json:"payload"`
}

func NewOfferCompleted(offerItems *OfferItems, groupKey keys.GroupKey) *OfferCompleted {
	return &OfferCompleted{
		eventsource.NewEventEnvelope(OfferCompletedEvent, 1),
		OfferCompletedPayload{
			OfferItems: offerItems,
			GroupKey:   groupKey,
		},
	}
}

var _ eventsource.Event = &OfferCompleted{}
