package domain

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type ResourceBorrowItem struct {
	OfferItemKey keys.OfferItemKey `json:"key"`
	Duration     time.Duration     `json:"duration"`
	ResourceKey  keys.ResourceKey  `json:"resource_key"`
	To           *OfferItemTarget  `json:"to"`
}

func (c *ResourceBorrowItem) MarshalJSON() ([]byte, error) {
	a := struct {
		ResourceBorrowItem
		Type string `json:"type"`
	}{
		ResourceBorrowItem: *c,
		Type:               string(BorrowResourceItemType),
	}
	return json.Marshal(a)
}

func (r ResourceBorrowItem) Type() OfferItemType {
	return BorrowResourceItemType
}

func (r ResourceBorrowItem) Key() keys.OfferItemKey {
	return r.OfferItemKey
}

var _ OfferItem = &ResourceBorrowItem{}
