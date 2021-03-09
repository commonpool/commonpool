package domain

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type BorrowResourceItem struct {
	OfferItemBase
	ResourceKey      keys.ResourceKey `json:"resourceId"`
	Duration         time.Duration    `json:"duration"`
	ItemTaken        bool             `json:"itemTaken"`
	ItemGiven        bool             `json:"itemGiven"`
	ItemReturnedBack bool             `json:"itemReturnedBack"`
	ItemReceivedBack bool             `json:"itemReceivedBack"`
}

func (b *BorrowResourceItem) AsCreditTransfer() (*CreditTransferItem, bool) {
	return nil, false
}

func (b *BorrowResourceItem) AsProvideService() (*ProvideServiceItem, bool) {
	return nil, false
}

func (b *BorrowResourceItem) AsBorrowResource() (*BorrowResourceItem, bool) {
	return b, true
}

func (b *BorrowResourceItem) AsResourceTransfer() (*ResourceTransferItem, bool) {
	return nil, false
}

func (b BorrowResourceItem) IsCompleted() bool {
	return b.ItemTaken && b.ItemGiven && b.ItemReturnedBack && b.ItemReceivedBack
}

func (b BorrowResourceItem) Type() OfferItemType {
	return BorrowResource
}

func (b BorrowResourceItem) GetResourceKey() keys.ResourceKey {
	return b.ResourceKey
}

func (b BorrowResourceItem) GetTo() keys.Target {
	return *b.To
}

var _ OfferItem = &BorrowResourceItem{}

type NewBorrowResourceItemOptions struct {
	ItemGiven        bool
	ItemReceived     bool
	ItemReturnedBack bool
	ItemReceivedBack bool
	ReceiverAccepted bool
	GiverAccepted    bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewBorrowResourceItem(
	offerKey keys.OfferKey,
	offerItemKey keys.OfferItemKey,
	resourceKey keys.ResourceKey,
	to *keys.Target,
	duration time.Duration,
	options ...NewBorrowResourceItemOptions) *BorrowResourceItem {

	now := time.Now()
	defaultOptions := &NewBorrowResourceItemOptions{
		ItemGiven:        false,
		ItemReceived:     false,
		ItemReturnedBack: false,
		ItemReceivedBack: false,
		ReceiverAccepted: false,
		GiverAccepted:    false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if len(options) > 0 {
		option := options[0]

		if option.ItemGiven != false {
			defaultOptions.ItemGiven = option.ItemGiven
		}
		if option.ItemReceived != false {
			defaultOptions.ItemReceived = option.ItemReceived
		}
		if option.ItemReturnedBack != false {
			defaultOptions.ItemReturnedBack = option.ItemReturnedBack
		}
		if option.ItemReceivedBack != false {
			defaultOptions.ItemReceivedBack = option.ItemReceivedBack
		}
		if option.GiverAccepted != false {
			defaultOptions.GiverAccepted = option.GiverAccepted
		}
		if option.ReceiverAccepted != false {
			defaultOptions.ReceiverAccepted = option.ReceiverAccepted
		}
		if option.CreatedAt != time.Unix(0, 0) {
			defaultOptions.CreatedAt = option.CreatedAt
		}
		if option.UpdatedAt != time.Unix(0, 0) {
			defaultOptions.UpdatedAt = option.UpdatedAt
		}
	}

	item := &BorrowResourceItem{
		OfferItemBase: OfferItemBase{
			Type:             BorrowResource,
			Key:              offerItemKey,
			OfferKey:         offerKey,
			To:               to,
			ApprovedOutbound: defaultOptions.GiverAccepted,
			ApprovedInbound:  defaultOptions.ReceiverAccepted,
			CreatedAt:        defaultOptions.CreatedAt,
			UpdatedAt:        defaultOptions.UpdatedAt,
		},
		ResourceKey:      resourceKey,
		Duration:         duration,
		ItemTaken:        defaultOptions.ItemReceived,
		ItemGiven:        defaultOptions.ItemGiven,
		ItemReturnedBack: defaultOptions.ItemReturnedBack,
		ItemReceivedBack: defaultOptions.ItemReceivedBack,
	}

	return item
}
