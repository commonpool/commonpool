package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ServiceReceivedNotifiedPayload struct {
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type ServiceReceivedNotified struct {
	eventsource.EventEnvelope
	ServiceReceivedNotifiedPayload `json:"payload"`
}

func NewServiceReceivedNotified(notifiedBy keys.UserKey, serviceOfferItemKey keys.OfferItemKey) *ServiceReceivedNotified {
	return &ServiceReceivedNotified{
		eventsource.NewEventEnvelope(ServiceReceivedNotifiedEvent, 1),
		ServiceReceivedNotifiedPayload{
			notifiedBy,
			serviceOfferItemKey,
		},
	}
}

var _ eventsource.Event = &ServiceReceivedNotified{}
