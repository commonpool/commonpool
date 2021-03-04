package domain

import "github.com/commonpool/backend/pkg/keys"

type ServiceReceivedNotified struct {
	Type         OfferEvent        `json:"type"`
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
	Version      int               `json:"version"`
}

func NewServiceReceivedNotified(notifiedBy keys.UserKey, serviceOfferItemKey keys.OfferItemKey) *ServiceReceivedNotified {
	return &ServiceReceivedNotified{
		Type:         ServiceReceivedNotifiedEvent,
		NotifiedBy:   notifiedBy,
		OfferItemKey: serviceOfferItemKey,
		Version:      1,
	}
}

func (o *ServiceReceivedNotified) GetType() OfferEvent {
	return o.Type
}

func (o *ServiceReceivedNotified) GetVersion() int {
	return o.Version
}

var _ Event = &ServiceReceivedNotified{}
