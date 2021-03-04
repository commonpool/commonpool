package domain

import "github.com/commonpool/backend/pkg/keys"

type ResourceTransferReceivedNotified struct {
	Type         OfferEvent        `json:"type"`
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
	Version      int               `json:"version"`
}

func NewResourceReceivedNotified(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *ResourceTransferReceivedNotified {
	return &ResourceTransferReceivedNotified{
		Type:         ResourceTransferReceivedNotifiedEvent,
		NotifiedBy:   notifiedBy,
		OfferItemKey: resourceOfferItemKey,
		Version:      1,
	}
}

func (o *ResourceTransferReceivedNotified) GetType() OfferEvent {
	return o.Type
}

func (o *ResourceTransferReceivedNotified) GetVersion() int {
	return o.Version
}

var _ Event = &ResourceTransferReceivedNotified{}
