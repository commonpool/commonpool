package domain

import "github.com/commonpool/backend/pkg/keys"

type ResourceLentNotified struct {
	Type         OfferEvent        `json:"type"`
	NotifiedBy   keys.UserKey      `json:"notified_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
	Version      int               `json:"version"`
}

func NewResourceLentNotified(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) *ResourceLentNotified {
	return &ResourceLentNotified{
		Type:         ResourceLentNotifiedEvent,
		NotifiedBy:   notifiedBy,
		OfferItemKey: resourceOfferItemKey,
		Version:      1,
	}
}

func (o *ResourceLentNotified) GetType() OfferEvent {
	return o.Type
}

func (o *ResourceLentNotified) GetVersion() int {
	return o.Version
}

var _ Event = &ResourceLentNotified{}
