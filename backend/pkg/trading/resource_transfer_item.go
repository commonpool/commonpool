package trading

import (
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceTransferItem struct {
	OfferItemBase
	ResourceKey  keys.ResourceKey
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
