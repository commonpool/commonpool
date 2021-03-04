package domain

import "github.com/commonpool/backend/pkg/keys"

type ResourceTransferGivenNotified struct {
	Type         OfferEvent        `json:"type"`
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
	Version      int               `json:"version"`
}

func NewResourceGivenNotified(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *ResourceTransferGivenNotified {
	return &ResourceTransferGivenNotified{
		Type:         ResourceTransferGivenNotifiedEvent,
		NotifiedBy:   notifiedBy,
		OfferItemKey: resourceOfferItemKey,
		Version:      1,
	}
}

func (o *ResourceTransferGivenNotified) GetType() OfferEvent {
	return o.Type
}

func (o *ResourceTransferGivenNotified) GetVersion() int {
	return o.Version
}

var _ Event = &ResourceTransferGivenNotified{}
