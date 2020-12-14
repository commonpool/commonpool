package trading

import "github.com/commonpool/backend/model"

type ResourceTransferItem struct {
	OfferItemBase
	ResourceKey  model.ResourceKey
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
