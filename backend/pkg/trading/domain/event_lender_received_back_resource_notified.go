package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type LenderReceivedBackResourceNotifiedPayload struct {
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type LenderReceivedBackResourceNotified struct {
	eventsource.EventEnvelope
	LenderReceivedBackResourceNotifiedPayload `json:"payload"`
}

func NewLenderReceivedBackResource(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *LenderReceivedBackResourceNotified {
	return &LenderReceivedBackResourceNotified{
		eventsource.NewEventEnvelope(LenderReceivedBackResourceEvent, 1),
		LenderReceivedBackResourceNotifiedPayload{
			notifiedBy,
			resourceOfferItemKey,
		},
	}
}

var _ eventsource.Event = &LenderReceivedBackResourceNotified{}
