package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type BorrowerReturnedResourceNotifiedPayload struct {
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type BorrowerReturnedResourceNotified struct {
	eventsource.EventEnvelope
	BorrowerReturnedResourceNotifiedPayload `json:"payload"`
}

func NewBorrowerReturnedResource(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *BorrowerReturnedResourceNotified {
	return &BorrowerReturnedResourceNotified{
		eventsource.NewEventEnvelope(BorrowerReturnedResourceEvent, 1),
		BorrowerReturnedResourceNotifiedPayload{
			notifiedBy,
			resourceOfferItemKey,
		},
	}
}

var _ eventsource.Event = &BorrowerReturnedResourceNotified{}
