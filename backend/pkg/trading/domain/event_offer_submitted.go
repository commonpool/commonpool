package domain

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
)

type OfferItems2 struct {
	Items []trading.OfferItem
}

func (o *OfferItems2) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Items)
}

func (o *OfferItems2) UnmarshalJSON(data []byte) error {

	type str struct {
		Type OfferItemType `json:"type"`
	}

	var tempSlice []map[string]interface{}

	err := json.Unmarshal(data, &tempSlice)
	if err != nil {
		return err
	}

	for _, tempItem := range tempSlice {

		itemType, ok := tempItem["type"]
		if !ok {
			return fmt.Errorf("item does not have a 'type' property")
		}

		itemJs, err := json.Marshal(tempItem)
		if err != nil {
			return err
		}

		var destination trading.OfferItem
		switch itemType {
		case string(CreditTransferItemType):
			destination = &CreditTransferItem{}
		case string(ProvideServiceItemType):
			destination = &ServiceOfferItem{}
		case string(ResourceTransferItemType):
			destination = &ResourceTransferItem{}
		case string(BorrowResourceItemType):
			destination = &ResourceBorrowItem{}
		default:
			return fmt.Errorf("unexpected item type")
		}

		err = json.Unmarshal(itemJs, destination)
		if err != nil {
			return err
		}

		o.Items = append(o.Items, destination)

	}

	return nil

}

type OfferSubmitted struct {
	Type       OfferEvent         `json:"type"`
	OfferItems trading.OfferItems `json:"offer_items"`
	GroupKey   keys.GroupKey      `json:"group_key"`
	Version    int                `json:"version"`
}

func NewOfferSubmitted(offerItems trading.OfferItems, groupKey keys.GroupKey) *OfferSubmitted {
	return &OfferSubmitted{
		Type:       OfferSubmittedEvent,
		OfferItems: offerItems,
		GroupKey:   groupKey,
		Version:    1,
	}
}

func (o OfferSubmitted) GetType() OfferEvent {
	return OfferSubmittedEvent
}

func (o *OfferSubmitted) GetVersion() int {
	return o.Version
}

var _ Event = &OfferSubmitted{}
