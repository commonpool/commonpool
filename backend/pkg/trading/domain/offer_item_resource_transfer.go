package domain

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type ResourceTransferItem struct {
	OfferItemBase
	ResourceKey  keys.ResourceKey `json:"resourceId"`
	ItemGiven    bool             `json:"itemGiven"`
	ItemReceived bool             `json:"itemReceived"`
}

func (r *ResourceTransferItem) AsCreditTransfer() (*CreditTransferItem, bool) {
	return nil, false
}

func (r *ResourceTransferItem) AsProvideService() (*ProvideServiceItem, bool) {
	return nil, false
}

func (r *ResourceTransferItem) AsBorrowResource() (*BorrowResourceItem, bool) {
	return nil, false
}

func (r *ResourceTransferItem) AsResourceTransfer() (*ResourceTransferItem, bool) {
	return r, true
}

func (r ResourceTransferItem) IsCompleted() bool {
	return r.ItemGiven && r.ItemReceived
}

func (b ResourceTransferItem) GetResourceKey() keys.ResourceKey {
	return b.ResourceKey
}

func (b ResourceTransferItem) GetTo() keys.Target {
	return *b.To
}

func (r ResourceTransferItem) Type() OfferItemType {
	return ResourceTransfer
}

var _ OfferItem = &ResourceTransferItem{}

type NewResourceTransferItemOptions struct {
	ReceiverAccepted bool
	GiverAccepted    bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ItemGiven        bool
	ItemReceived     bool
}

func NewResourceTransferItem(
	offerKey keys.OfferKey,
	offerItemKey keys.OfferItemKey,
	to *keys.Target,
	resourceKey keys.ResourceKey,
	options ...NewResourceTransferItemOptions) *ResourceTransferItem {

	now := time.Now().UTC()

	defaultOptions := NewResourceTransferItemOptions{
		ReceiverAccepted: false,
		GiverAccepted:    false,
		CreatedAt:        now,
		UpdatedAt:        now,
		ItemGiven:        false,
		ItemReceived:     false,
	}

	if len(options) > 0 {
		option := options[0]

		if option.ReceiverAccepted != false {
			defaultOptions.ReceiverAccepted = option.ReceiverAccepted
		}
		if option.GiverAccepted != false {
			defaultOptions.GiverAccepted = option.GiverAccepted
		}
		if option.CreatedAt != time.Unix(0, 0).UTC() {
			defaultOptions.CreatedAt = option.CreatedAt.UTC()
		}
		if option.UpdatedAt != time.Unix(0, 0).UTC() {
			defaultOptions.UpdatedAt = option.UpdatedAt.UTC()
		}
		if option.ItemGiven != false {
			defaultOptions.ItemGiven = option.ItemGiven
		}
		if option.ItemReceived != false {
			defaultOptions.ItemReceived = option.ItemReceived
		}
	}

	return &ResourceTransferItem{
		OfferItemBase: OfferItemBase{
			Type:             ResourceTransfer,
			Key:              offerItemKey,
			OfferKey:         offerKey,
			To:               to,
			ApprovedInbound:  defaultOptions.ReceiverAccepted,
			ApprovedOutbound: defaultOptions.GiverAccepted,
			CreatedAt:        defaultOptions.CreatedAt,
			UpdatedAt:        defaultOptions.UpdatedAt,
		},
		ResourceKey:  resourceKey,
		ItemGiven:    defaultOptions.ItemGiven,
		ItemReceived: defaultOptions.ItemReceived,
	}

}
