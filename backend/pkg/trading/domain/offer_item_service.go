package domain

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type ServiceOfferItem struct {
	OfferItemKey keys.OfferItemKey `json:"key"`
	Duration     time.Duration     `json:"duration"`
	ResourceKey  keys.ResourceKey  `json:"resource_key"`
	To           *OfferItemTarget  `json:"to"`
	From         *OfferItemTarget  `json:"from"`
}

func (c *ServiceOfferItem) MarshalJSON() ([]byte, error) {
	a := struct {
		ServiceOfferItem
		Type string `json:"type"`
	}{
		ServiceOfferItem: *c,
		Type:             string(ProvideServiceItemType),
	}
	return json.Marshal(a)
}

func (r ServiceOfferItem) Type() OfferItemType {
	return ProvideServiceItemType
}

func (r ServiceOfferItem) Key() keys.OfferItemKey {
	return r.OfferItemKey
}

var _ OfferItem = &ServiceOfferItem{}
