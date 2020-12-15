package model

import (
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	"time"
)

type BorrowResourceItem struct {
	OfferItemBase
	ResourceKey      resourcemodel.ResourceKey
	Duration         time.Duration
	ItemTaken        bool
	ItemGiven        bool
	ItemReturnedBack bool
	ItemReceivedBack bool
}

func (b BorrowResourceItem) IsCompleted() bool {
	return b.ItemTaken && b.ItemGiven && b.ItemReturnedBack && b.ItemReceivedBack
}

func (b BorrowResourceItem) Type() OfferItemType {
	return BorrowResource
}

var _ OfferItem = &BorrowResourceItem{}
