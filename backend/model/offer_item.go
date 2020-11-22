package model

import uuid "github.com/satori/go.uuid"

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
