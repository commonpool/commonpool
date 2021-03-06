package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceTransferReceivedNotifiedPayload struct {
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type ResourceTransferReceivedNotified struct {
	eventsource.EventEnvelope
	ResourceTransferReceivedNotifiedPayload `json:"payload"`
}

func NewResourceReceivedNotified(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *ResourceTransferReceivedNotified {
	return &ResourceTransferReceivedNotified{
		eventsource.NewEventEnvelope(ResourceTransferReceivedNotifiedEvent, 1),
		ResourceTransferReceivedNotifiedPayload{
			notifiedBy,
			resourceOfferItemKey,
		},
	}
}

var _ eventsource.Event = &ResourceTransferReceivedNotified{}
