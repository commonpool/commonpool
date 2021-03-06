package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceBorrowedNotifiedPayload struct {
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type ResourceBorrowedNotified struct {
	eventsource.EventEnvelope
	ResourceBorrowedNotifiedPayload `json:"payload"`
}

func NewResourceBorrowedNotified(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *ResourceBorrowedNotified {
	return &ResourceBorrowedNotified{
		eventsource.NewEventEnvelope(ResourceBorrowedNotifiedEvent, 1),
		ResourceBorrowedNotifiedPayload{
			notifiedBy,
			resourceOfferItemKey,
		},
	}
}

var _ eventsource.Event = &ResourceBorrowedNotified{}
