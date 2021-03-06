package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceLentNotifiedPayload struct {
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type ResourceLentNotified struct {
	eventsource.EventEnvelope
	ResourceLentNotifiedPayload `json:"payload"`
}

func NewResourceLentNotified(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *ResourceLentNotified {
	return &ResourceLentNotified{
		eventsource.NewEventEnvelope(ResourceLentNotifiedEvent, 1),
		ResourceLentNotifiedPayload{
			notifiedBy,
			resourceOfferItemKey,
		},
	}
}

var _ eventsource.Event = &ResourceLentNotified{}
