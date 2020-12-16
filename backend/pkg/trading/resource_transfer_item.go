package trading

import (
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
)

type ResourceTransferItem struct {
	OfferItemBase
	ResourceKey  resourcemodel.ResourceKey
	ItemGiven    bool
	ItemReceived bool
}

func (r ResourceTransferItem) IsCompleted() bool {
	return r.ItemGiven && r.ItemReceived
}

func (r ResourceTransferItem) Type() OfferItemType {
	return ResourceTransfer
}

var _ OfferItem = &ResourceTransferItem{}
