package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ServiceGivenNotifiedPayload struct {
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type ServiceGivenNotified struct {
	eventsource.EventEnvelope
	ServiceGivenNotifiedPayload
}

func NewServiceGivenNotified(notifiedBy keys.UserKey, serviceOfferItemKey keys.OfferItemKey) *ServiceGivenNotified {
	return &ServiceGivenNotified{
		eventsource.NewEventEnvelope(ServiceGivenNotifiedEvent, 1),
		ServiceGivenNotifiedPayload{
			notifiedBy,
			serviceOfferItemKey,
		},
	}
}

var _ eventsource.Event = &ServiceGivenNotified{}
