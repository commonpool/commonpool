package trading

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type BorrowResourceItem struct {
	OfferItemBase
	ResourceKey      keys.ResourceKey
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
