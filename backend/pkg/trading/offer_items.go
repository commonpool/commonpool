package trading

import (
	"encoding/json"
	"errors"
	"github.com/commonpool/backend/pkg/keys"
)

type OfferItems struct {
	Items []OfferItem
}

func (i *OfferItems) AllUserActionsCompleted() bool {
	for _, item := range i.Items {
		if !item.IsCreditTransfer() && !item.IsCompleted() {
			return false
		}
	}
	return true
}

func (i *OfferItems) AllApproved() bool {
	for _, item := range i.Items {
		if !item.IsAccepted() {
			return false
		}
	}
	return true
}

func NewOfferItems(offerItems []OfferItem) *OfferItems {
	copied := make([]OfferItem, len(offerItems))
	copy(copied, offerItems)
	return &OfferItems{
		Items: copied,
	}
}

func NewOfferItemsFrom(offerItems ...OfferItem) *OfferItems {
	return &OfferItems{
		Items: offerItems,
	}
}

func (i *OfferItems) GetOfferItem(key keys.OfferItemKey) OfferItem {
	for _, offerItem := range i.Items {
		if offerItem.GetKey() == key {
			return offerItem
		}
	}
	return nil
}

func (i *OfferItems) ItemCount() int {
	return len(i.Items)
}

func (i *OfferItems) GetResourceKeys() *keys.ResourceKeys {
	var resourceKeys []keys.ResourceKey
	for _, item := range i.Items {
		if item.IsBorrowingResource() {
			resourceKeys = append(resourceKeys, item.(*BorrowResourceItem).ResourceKey)
		} else if item.IsServiceProviding() {
			resourceKeys = append(resourceKeys, item.(*ProvideServiceItem).ResourceKey)
		} else if item.IsResourceTransfer() {
			resourceKeys = append(resourceKeys, item.(*ResourceTransferItem).ResourceKey)
		}
	}
	if resourceKeys == nil {
		resourceKeys = []keys.ResourceKey{}
	}
	return keys.NewResourceKeys(resourceKeys)
}

func (i *OfferItems) GetUserKeys() *keys.UserKeys {
	var userKeys []keys.UserKey
	for _, offerItem := range i.Items {
		if offerItem.GetReceiverKey().IsForUser() {
			userKeys = append(userKeys, offerItem.GetReceiverKey().GetUserKey())
		}
		if offerItem.IsCreditTransfer() {
			creditTransfer := offerItem.(*CreditTransferItem)
			if creditTransfer.From.IsForUser() {
				userKeys = append(userKeys, creditTransfer.From.GetUserKey())
			}
		}
	}
	return keys.NewUserKeys(userKeys)
}

func (i *OfferItems) GetGroupKeys() *keys.GroupKeys {
	var groupKeys []keys.GroupKey
	for _, offerItem := range i.Items {
		if offerItem.GetReceiverKey().IsForGroup() {
			groupKeys = append(groupKeys, offerItem.GetReceiverKey().GetGroupKey())
		}
		if offerItem.IsCreditTransfer() {
			creditTransfer := offerItem.(*CreditTransferItem)
			if creditTransfer.From.IsForGroup() {
				groupKeys = append(groupKeys, creditTransfer.From.GetGroupKey())
			}
		}
	}
	return keys.NewGroupKeys(groupKeys)
}

func (i *OfferItems) IsEmpty() bool {
	return i.Items == nil || len(i.Items) == 0
}

func (i *OfferItems) Count() int {
	if i.Items == nil {
		return 0
	}
	return len(i.Items)
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
