package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type OfferDeclinedPayload struct {
	DeclinedBy keys.UserKey `json:"declined_by"`
}

type OfferDeclined struct {
	eventsource.EventEnvelope
	OfferDeclinedPayload `json:"payload"`
}

func NewOfferDeclined(declinedBy keys.UserKey) *OfferDeclined {
	return &OfferDeclined{
		eventsource.NewEventEnvelope(OfferDeclinedEvent, 1),
		OfferDeclinedPayload{
			declinedBy,
		},
	}
}

var _ eventsource.Event = &OfferDeclined{}
