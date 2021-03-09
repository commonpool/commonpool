package domain

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type SubmitOfferItems []SubmitOfferItem

func NewSubmitOfferItems(offerItems ...SubmitOfferItem) SubmitOfferItems {
	return offerItems
}

type SubmitOfferItemBase struct {
	OfferItemType OfferItemType
	ResourceKey   *keys.ResourceKey
	From          *keys.Target
	To            *keys.Target
	Amount        *time.Duration
	Duration      *time.Duration
}

type SubmitOfferItem struct {
	SubmitOfferItemBase
	OfferItemKey keys.OfferItemKey
}

func NewResourceTransferItemInputBase(to keys.Targetter, resourceKey keys.ResourceKeyGetter) SubmitOfferItemBase {
	rkey := resourceKey.GetResourceKey()
	return SubmitOfferItemBase{
		OfferItemType: ResourceTransfer,
		ResourceKey:   &rkey,
		To:            to.Target(),
	}
}

func NewCreditTransferItemInputBase(from keys.Targetter, to keys.Targetter, amount time.Duration) SubmitOfferItemBase {
	return SubmitOfferItemBase{
		OfferItemType: CreditTransfer,
		From:          from.Target(),
		To:            to.Target(),
		Amount:        &amount,
	}
}

func NewProvideServiceItemInputBase(from keys.Targetter, to keys.Targetter, resourceKey keys.ResourceKey, duration time.Duration) SubmitOfferItemBase {
	return SubmitOfferItemBase{
		OfferItemType: ProvideService,
		From:          from.Target(),
		To:            to.Target(),
		ResourceKey:   &resourceKey,
		Duration:      &duration,
	}
}

func NewBorrowResourceInputBase(to keys.Targetter, resourceKey keys.ResourceKey, duration time.Duration) SubmitOfferItemBase {
	return SubmitOfferItemBase{
		OfferItemType: BorrowResource,
		To:            to.Target(),
		ResourceKey:   &resourceKey,
		Duration:      &duration,
	}
}

func NewResourceTransferItemInput(offerItemKey keys.OfferItemKey, to keys.Targetter, resourceKey keys.ResourceKey) SubmitOfferItem {
	return SubmitOfferItem{
		OfferItemKey:        offerItemKey,
		SubmitOfferItemBase: NewResourceTransferItemInputBase(to, resourceKey),
	}
}

func NewCreditTransferItemInput(offerItemKey keys.OfferItemKey, from keys.Targetter, to keys.Targetter, amount time.Duration) SubmitOfferItem {
	return SubmitOfferItem{
		OfferItemKey:        offerItemKey,
		SubmitOfferItemBase: NewCreditTransferItemInputBase(from, to, amount),
	}
}

func NewProvideServiceItemInput(offerItemKey keys.OfferItemKey, from keys.Targetter, to keys.Targetter, resourceKey keys.ResourceKey, duration time.Duration) SubmitOfferItem {
	return SubmitOfferItem{
		OfferItemKey:        offerItemKey,
		SubmitOfferItemBase: NewProvideServiceItemInputBase(from, to, resourceKey, duration),
	}
}

func NewBorrowResourceInput(offerItemKey keys.OfferItemKey, to keys.Targetter, resourceKey keys.ResourceKey, duration time.Duration) SubmitOfferItem {
	return SubmitOfferItem{
		OfferItemKey:        offerItemKey,
		SubmitOfferItemBase: NewBorrowResourceInputBase(to, resourceKey, duration),
	}
}
