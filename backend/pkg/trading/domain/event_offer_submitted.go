package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type OfferSubmittedPayload struct {
	OfferItems  *OfferItems   `json:"offer_items"`
	GroupKey    keys.GroupKey `json:"group_key"`
	SubmittedBy keys.UserKey  `json:"submitted_by"`
}

type OfferSubmitted struct {
	eventsource.EventEnvelope
	OfferSubmittedPayload `json:"payload"`
}

func NewOfferSubmitted(by keys.UserKey, offerItems *OfferItems, groupKey keys.GroupKey) *OfferSubmitted {
	return &OfferSubmitted{
		eventsource.NewEventEnvelope(OfferSubmittedEvent, 1),
		OfferSubmittedPayload{
			offerItems,
			groupKey,
			by,
		},
	}
}

var _ eventsource.Event = &OfferSubmitted{}
