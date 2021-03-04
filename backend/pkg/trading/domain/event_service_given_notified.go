package domain

import "github.com/commonpool/backend/pkg/keys"

type ServiceGivenNotified struct {
	Type         OfferEvent        `json:"type"`
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
	Version      int               `json:"version"`
}

func NewServiceGivenNotified(notifiedBy keys.UserKey, serviceOfferItemKey keys.OfferItemKey) *ServiceGivenNotified {
	return &ServiceGivenNotified{
		Type:         ServiceGivenNotifiedEvent,
		NotifiedBy:   notifiedBy,
		OfferItemKey: serviceOfferItemKey,
		Version:      1,
	}
}

func (o *ServiceGivenNotified) GetType() OfferEvent {
	return o.Type
}

func (o *ServiceGivenNotified) GetVersion() int {
	return o.Version
}

var _ Event = &ServiceGivenNotified{}
