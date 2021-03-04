package domain

import "github.com/commonpool/backend/pkg/keys"

type ResourceBorrowedNotified struct {
	Type         OfferEvent        `json:"type"`
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
	Version      int               `json:"version"`
}

func NewResourceBorrowedNotified(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *ResourceBorrowedNotified {
	return &ResourceBorrowedNotified{
		Type:         ResourceBorrowedNotifiedEvent,
		NotifiedBy:   notifiedBy,
		OfferItemKey: resourceOfferItemKey,
		Version:      1,
	}
}

func (o *ResourceBorrowedNotified) GetType() OfferEvent {
	return o.Type
}

func (o *ResourceBorrowedNotified) GetVersion() int {
	return o.Version
}

var _ Event = &ResourceBorrowedNotified{}
