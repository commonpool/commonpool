package domain

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceTransferItem struct {
	OfferItemKey keys.OfferItemKey `json:"key"`
	ResourceKey  keys.ResourceKey  `json:"resource_key"`
	To           *OfferItemTarget  `json:"to"`
}

func (c *ResourceTransferItem) MarshalJSON() ([]byte, error) {
	a := struct {
		ResourceTransferItem
		Type string `json:"type"`
	}{
		ResourceTransferItem: *c,
		Type:                 string(ResourceTransferItemType),
	}
	return json.Marshal(a)
}

func (r ResourceTransferItem) Type() OfferItemType {
	return ResourceTransferItemType
}

func (r ResourceTransferItem) Key() keys.OfferItemKey {
	return r.OfferItemKey
}

var _ OfferItem = &ResourceTransferItem{}
