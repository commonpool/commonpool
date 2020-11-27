package trading

import (
	"github.com/commonpool/backend/model"
	"github.com/satori/go.uuid"
	"time"
)

type OfferStatus int

const (
	PendingOffer OfferStatus = iota
	AcceptedOffer
	CanceledOffer
	DeclinedOffer
	ExpiredOffer
	CompletedOffer
)

type Offer struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	AuthorID       string
	Status         OfferStatus
	CreatedAt      time.Time
	ExpirationTime *time.Time
	// Completion time (when either accepted, expired or declined)
	CompletedAt *time.Time
	Message     string
}

func NewOffer(offerKey model.OfferKey, author model.UserKey, message string, expiration *time.Time) Offer {
	return Offer{
		ID:             offerKey.ID,
		AuthorID:       author.String(),
		Status:         PendingOffer,
		ExpirationTime: expiration,
		Message:        message,
	}
}

func (o *Offer) GetKey() model.OfferKey {
	return model.NewOfferKey(o.ID)
}

func (o *Offer) GetAuthorKey() model.UserKey {
	return model.NewUserKey(o.AuthorID)
}

type Decision int

const (
	PendingDecision Decision = iota
	AcceptedDecision
	DeclinedDecision
)

type OfferDecision struct {
	OfferID  uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID   string    `gorm:"primary_key"`
	Decision Decision
}

func (d *OfferDecision) GetKey() model.OfferDecisionKey {
	return model.NewOfferDecisionKey(d.GetOfferKey(), d.GetUserKey())
}

func (d *OfferDecision) GetOfferKey() model.OfferKey {
	return model.NewOfferKey(d.OfferID)
}
func (d *OfferDecision) GetUserKey() model.UserKey {
	return model.NewUserKey(d.UserID)
}

type OfferItemType int

const (
	ResourceItem OfferItemType = iota
	TimeItem
)

type OfferItemBond string

const (
	OfferItemGiving    OfferItemBond = "giving"
	OfferItemReceiving OfferItemBond = "receiving"
	OfferItemNeither   OfferItemBond = "none"
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
	Received   bool
	Given      bool
}

func (i *OfferItem) FormatOfferedTimeInSeconds() string {
	return time.Duration(int64(time.Second) * *i.OfferedTimeInSeconds).Truncate(time.Minute * 1).String()
}

func (i *OfferItem) GetUserBondDirection(user model.UserKey) OfferItemBond {
	if i.GetFromUserKey() == user {
		return OfferItemGiving
	} else if i.GetToUserKey() == user {
		return OfferItemReceiving
	} else {
		return OfferItemNeither
	}
}

func (i *OfferItem) IsReceived() bool {
	return i.Received
}

func (i *OfferItem) IsGiven() bool {
	return i.Given
}

func NewTimeOfferItem(offerKey model.OfferKey, key model.OfferItemKey, fromUser model.UserKey, toUser model.UserKey, offeredSeconds int64) OfferItem {
	return OfferItem{
		ID:                   key.ID,
		ItemType:             TimeItem,
		OfferID:              offerKey.ID,
		OfferedTimeInSeconds: &offeredSeconds,
		FromUserID:           fromUser.String(),
		ToUserID:             toUser.String(),
		ResourceID:           nil,
		Received:             false,
		Given:                false,
	}
}

func NewResourceOfferItem(offerKey model.OfferKey, key model.OfferItemKey, fromUser model.UserKey, toUser model.UserKey, offeredResource model.ResourceKey) OfferItem {
	return OfferItem{
		ID:                   key.ID,
		ItemType:             ResourceItem,
		OfferID:              offerKey.ID,
		OfferedTimeInSeconds: nil,
		FromUserID:           fromUser.String(),
		ToUserID:             toUser.String(),
		ResourceID:           &offeredResource.ID,
		Received:             false,
		Given:                false,
	}
}

func (i *OfferItem) IsTimeExchangeItem() bool {
	return i.ItemType == TimeItem
}

func (i *OfferItem) IsResourceExchangeItem() bool {
	return i.ItemType == ResourceItem
}

func (i *OfferItem) GetOfferKey() model.OfferKey {
	return model.NewOfferKey(i.OfferID)
}

func (i *OfferItem) GetKey() model.OfferItemKey {
	return model.NewOfferItemKey(i.ID)
}

func (i *OfferItem) GetFromUserKey() model.UserKey {
	return model.NewUserKey(i.FromUserID)
}

func (i *OfferItem) GetToUserKey() model.UserKey {
	return model.NewUserKey(i.ToUserID)
}

func (i *OfferItem) IsReceivedBy(userKey model.UserKey) bool {
	return i.GetToUserKey() == userKey
}

func (i *OfferItem) IsGivenBy(userKey model.UserKey) bool {
	return i.GetFromUserKey() == userKey
}

func (i *OfferItem) GetResourceKey() model.ResourceKey {
	return model.NewResourceKey(*i.ResourceID)
}

type OfferItems struct {
	Items        []OfferItem
	Users        []model.UserKey
	ItemsPerUser map[model.UserKey][]OfferItem
}

func (i *OfferItems) AllResourceItemsReceivedAndGiven() bool {
	for _, item := range i.Items {
		if item.IsTimeExchangeItem() {
			continue
		}
		if !item.Received || !item.Given {
			return false
		}
	}
	return true
}

func NewOfferItems(offerItems []OfferItem) *OfferItems {

	itemsMap := map[model.OfferItemKey]bool{}
	itemsPerUser := map[model.UserKey][]OfferItem{}
	var items []OfferItem
	var userKeys []model.UserKey

	for _, item := range offerItems {
		if _, ok := itemsMap[item.GetKey()]; ok {
			continue
		}
		itemsMap[item.GetKey()] = true

		toUser := item.GetToUserKey()
		fromUser := item.GetFromUserKey()
		if _, ok := itemsPerUser[toUser]; !ok {
			itemsPerUser[toUser] = []OfferItem{}
			userKeys = append(userKeys, toUser)
		}
		if _, ok := itemsPerUser[fromUser]; !ok {
			itemsPerUser[fromUser] = []OfferItem{}
			userKeys = append(userKeys, fromUser)
		}
		itemsPerUser[toUser] = append(itemsPerUser[toUser], item)
		itemsPerUser[fromUser] = append(itemsPerUser[fromUser], item)
		items = append(items, item)
	}

	return &OfferItems{
		Items:        items,
		Users:        userKeys,
		ItemsPerUser: itemsPerUser,
	}

}

func (i *OfferItems) GetOfferItemsForUser(userKey model.UserKey) *OfferItems {
	itemsForUser, ok := i.ItemsPerUser[userKey]
	if !ok {
		return NewOfferItems([]OfferItem{})
	}
	return NewOfferItems(itemsForUser)
}

func (i *OfferItems) Append(offerItem OfferItem) *OfferItems {
	newOfferItems := append(i.Items, offerItem)
	return NewOfferItems(newOfferItems)
}

func (i *OfferItems) ItemCount() int {
	return len(i.Items)
}

func (i *OfferItems) HasItemsForUser(userKey model.UserKey) bool {
	itemsForUser, ok := i.ItemsPerUser[userKey]
	if !ok {
		return false
	}
	if len(itemsForUser) == 0 {
		return false
	}
	return true
}

func (i *OfferItems) GetUserKeys() *model.UserKeys {
	var userKeys []model.UserKey
	for _, user := range i.Users {
		userKeys = append(userKeys, user)
	}
	return model.NewUserKeys(userKeys)
}

func (i *OfferItems) GetResourceKeys() []model.ResourceKey {
	var resourceKeys []model.ResourceKey
	for _, item := range i.Items {
		if item.IsTimeExchangeItem() {
			continue
		}
		resourceKeys = append(resourceKeys, item.GetResourceKey())
	}
	if resourceKeys == nil {
		resourceKeys = []model.ResourceKey{}
	}
	return resourceKeys
}
