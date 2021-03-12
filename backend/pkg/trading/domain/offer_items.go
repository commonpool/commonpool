package domain

import (
	"encoding/json"
	"errors"
)

type OfferItems struct {
	Items []OfferItem
}

func NewEmptyOfferItems() *OfferItems {
	return NewOfferItems([]OfferItem{})
}

func NewOfferItems(offerItems []OfferItem) *OfferItems {
	copied := make([]OfferItem, len(offerItems))
	copy(copied, offerItems)
	return &OfferItems{
		Items: copied,
	}
}

func (i *OfferItems) Append(offerItem OfferItem) {
	i.Items = append(i.Items, offerItem)
}

func (i *OfferItems) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Items)
}

func (i *OfferItems) UnmarshalJSON(bytes []byte) error {
	var slice []*json.RawMessage
	err := json.Unmarshal(bytes, &slice)
	if err != nil {
		return err
	}

	var offerItems = make([]OfferItem, len(slice))
	for i, jsonRaw := range slice {
		var m map[string]interface{}
		err := json.Unmarshal(*jsonRaw, &m)
		if err != nil {
			return err
		}
		if m["type"] == string(ProvideService) {
			var ps ProvideServiceItem
			if err := json.Unmarshal(*jsonRaw, &ps); err != nil {
				return err
			}
			offerItems[i] = &ps
		} else if m["type"] == string(BorrowResource) {
			var ps BorrowResourceItem
			if err := json.Unmarshal(*jsonRaw, &ps); err != nil {
				return err
			}
			offerItems[i] = &ps
		} else if m["type"] == string(ResourceTransfer) {
			var ps ResourceTransferItem
			if err := json.Unmarshal(*jsonRaw, &ps); err != nil {
				return err
			}
			offerItems[i] = &ps
		} else if m["type"] == string(CreditTransfer) {
			var ps CreditTransferItem
			if err := json.Unmarshal(*jsonRaw, &ps); err != nil {
				return err
			}
			offerItems[i] = &ps
		} else {
			return errors.New("unsupported type found")
		}
	}

	i.Items = offerItems

	return nil
}
