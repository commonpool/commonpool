package domain

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type JSONDuration struct {
	time.Duration
}

func (d *SubmitOfferItemBase) UnmarshalJSON(data []byte) error {

	type Temp struct {
		OfferItemType OfferItemType     `json:"type"`
		ResourceKey   *keys.ResourceKey `json:"resourceId"`
		From          *keys.Target      `json:"from"`
		To            *keys.Target      `json:"to"`
		Amount        *string           `json:"amount"`
		Duration      *string           `json:"duration"`
	}

	var tmp Temp
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	var amount *time.Duration
	if tmp.Amount != nil && *tmp.Amount != "" {
		parsedAmount, err := time.ParseDuration(*tmp.Amount)
		if err != nil {
			return err
		}
		amount = &parsedAmount
	}

	var duration *time.Duration
	if tmp.Duration != nil && *tmp.Duration != "" {
		parsedDuration, err := time.ParseDuration(*tmp.Duration)
		if err != nil {
			return err
		}
		duration = &parsedDuration
	}

	*d = SubmitOfferItemBase{
		OfferItemType: tmp.OfferItemType,
		ResourceKey:   tmp.ResourceKey,
		From:          tmp.From,
		To:            tmp.To,
		Amount:        amount,
		Duration:      duration,
	}
	return nil
}

type SubmitOfferItems []SubmitOfferItem

func NewSubmitOfferItems(offerItems ...SubmitOfferItem) SubmitOfferItems {
	return offerItems
}

type SubmitOfferItemBase struct {
	OfferItemType OfferItemType     `json:"type"`
	ResourceKey   *keys.ResourceKey `json:"resourceId"`
	From          *keys.Target      `json:"from"`
	To            *keys.Target      `json:"to"`
	Amount        *time.Duration    `json:"amount"`
	Duration      *time.Duration    `json:"duration"`
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
