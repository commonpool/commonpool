package domain

import "github.com/commonpool/backend/pkg/keys"

type BorrowerReturnedResourceNotified struct {
	Type         OfferEvent        `json:"type"`
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
	Version      int               `json:"version"`
}

func NewBorrowerReturnedResource(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *BorrowerReturnedResourceNotified {
	return &BorrowerReturnedResourceNotified{
		Type:         BorrowerReturnedResourceEvent,
		NotifiedBy:   notifiedBy,
		OfferItemKey: resourceOfferItemKey,
		Version:      1,
	}
}

func (o *BorrowerReturnedResourceNotified) GetType() OfferEvent {
	return o.Type
}

func (o *BorrowerReturnedResourceNotified) GetVersion() int {
	return o.Version
}

var _ Event = &BorrowerReturnedResourceNotified{}
