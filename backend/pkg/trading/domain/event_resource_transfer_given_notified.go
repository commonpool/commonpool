package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceTransferGivenNotifiedPayload struct {
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type ResourceTransferGivenNotified struct {
	eventsource.EventEnvelope
	ResourceTransferGivenNotifiedPayload `json:"payload"`
}

func NewResourceGivenNotified(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *ResourceTransferGivenNotified {
	return &ResourceTransferGivenNotified{
		eventsource.NewEventEnvelope(ResourceTransferGivenNotifiedEvent, 1),
		ResourceTransferGivenNotifiedPayload{
			notifiedBy,
			resourceOfferItemKey,
		},
	}
}

var _ eventsource.Event = &ResourceTransferGivenNotified{}
