package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type OfferItemType int

const (
	ResourceItem OfferItemType = iota
	TimeItem
)

type OfferItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	OfferID    uuid.UUID
	ItemType   OfferItemType
	FromUserID string
	ToUserID   string
	// Only available when ItemType = TimeItem
	OfferedTimeInSeconds *int64
	// Only available when ItemType = ResourceItem
	ResourceID *uuid.UUID
}

func (o *OfferItem) FormatOfferedTimeInSeconds() string {
	return time.Duration(int64(time.Second) * *o.OfferedTimeInSeconds).Truncate(time.Minute * 1).String()
}

func NewTimeOfferItem(key OfferItemKey, fromUser UserKey, toUser UserKey, offeredSeconds int64) OfferItem {
	return OfferItem{
		ID:                   key.ID,
		ItemType:             TimeItem,
		OfferID:              key.OfferID,
		OfferedTimeInSeconds: &offeredSeconds,
		FromUserID:           fromUser.String(),
		ToUserID:             toUser.String(),
		ResourceID:           nil,
	}
}

func NewResourceOfferItem(key OfferItemKey, fromUser UserKey, toUser UserKey, offeredResource ResourceKey) OfferItem {
	return OfferItem{
		ID:                   key.ID,
		ItemType:             ResourceItem,
		OfferID:              key.OfferID,
		OfferedTimeInSeconds: nil,
		FromUserID:           fromUser.String(),
		ToUserID:             toUser.String(),
		ResourceID:           &offeredResource.uuid,
	}
}

func (i *OfferItem) IsTimeExchangeItem() bool {
	return i.ItemType == TimeItem
}

func (i *OfferItem) IsResourceExchangeItem() bool {
	return i.ItemType == ResourceItem
}

func (i *OfferItem) GetOfferKey() OfferKey {
	return NewOfferKey(i.OfferID)
}

func (i *OfferItem) GetKey() OfferItemKey {
	return NewOfferItemKey(i.ID, i.GetOfferKey())
}

func (i *OfferItem) GetFromUserKey() UserKey {
	return NewUserKey(i.FromUserID)
}

func (i *OfferItem) GetToUserKey() UserKey {
	return NewUserKey(i.ToUserID)
}

func (i *OfferItem) GetResourceKey() ResourceKey {
	return NewResourceKey(*i.ResourceID)
}

type OfferItems struct {
	Items        []OfferItem
	Users        []UserKey
	ItemsPerUser map[UserKey][]OfferItem
}

func NewOfferItems(offerItems []OfferItem) *OfferItems {

	itemsPerUser := map[UserKey][]OfferItem{}
	var items []OfferItem
	var userKeys []UserKey

	for _, item := range offerItems {
		userKey := item.GetToUserKey()
		if _, ok := itemsPerUser[userKey]; !ok {
			itemsPerUser[userKey] = []OfferItem{}
			userKeys = append(userKeys, userKey)
		}
		itemsPerUser[userKey] = append(itemsPerUser[userKey], item)
		items = append(items, item)
	}

	return &OfferItems{
		Items:        items,
		Users:        userKeys,
		ItemsPerUser: itemsPerUser,
	}

}

func (o *OfferItems) GetOfferItemsForUser(userKey UserKey) *OfferItems {
	itemsForUser, ok := o.ItemsPerUser[userKey]
	if !ok {
		return NewOfferItems([]OfferItem{})
	}
	return NewOfferItems(itemsForUser)
}

func (o *OfferItems) Append(offerItem OfferItem) *OfferItems {
	newOfferItems := append(o.Items, offerItem)
	return NewOfferItems(newOfferItems)
}

func (o *OfferItems) ItemCount() int {
	return len(o.Items)
}

func (o *OfferItems) HasItemsForUser(userKey UserKey) bool {
	itemsForUser, ok := o.ItemsPerUser[userKey]
	if !ok {
		return false
	}
	if len(itemsForUser) == 0 {
		return false
	}
	return true
}

func (o *OfferItems) GetUserKeys() UserKeys {
	var userKeys []UserKey
	for _, user := range o.Users {
		userKeys = append(userKeys, user)
	}
	return NewUserKeys(userKeys)
}