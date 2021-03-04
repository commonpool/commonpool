package domain

import "github.com/commonpool/backend/pkg/keys"

type LenderReceivedBackResourceNotified struct {
	Type         OfferEvent        `json:"type"`
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
	Version      int               `json:"version"`
}

func NewLenderReceivedBackResource(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *LenderReceivedBackResourceNotified {
	return &LenderReceivedBackResourceNotified{
		Type:         LenderReceivedBackResourceEvent,
		NotifiedBy:   notifiedBy,
		OfferItemKey: resourceOfferItemKey,
		Version:      1,
	}
}

func (o *LenderReceivedBackResourceNotified) GetType() OfferEvent {
	return o.Type
}

func (o *LenderReceivedBackResourceNotified) GetVersion() int {
	return o.Version
}

var _ Event = &LenderReceivedBackResourceNotified{}
