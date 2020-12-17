package trading

import (
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

func (i *OfferItems) AllPartiesAccepted() bool {
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

func (i *OfferItems) GetOfferItem(key keys.OfferItemKey) OfferItem {
	for _, offerItem := range i.Items {
		if offerItem.GetKey() == key {
			return offerItem
		}
	}
	return nil
}

func (i *OfferItems) GetOfferItemsReceivedByUser(userKey keys.UserKey) *OfferItems {
	var offerItems []OfferItem
	for _, offerItem := range i.Items {
		if offerItem.GetReceiverKey().IsForUser() && offerItem.GetReceiverKey().GetUserKey() == userKey {
			offerItems = append(offerItems, offerItem)
		}
	}
	return NewOfferItems(offerItems)
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
