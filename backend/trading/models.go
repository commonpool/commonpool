package trading

import (
	"github.com/commonpool/backend/model"
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

type OfferItemTargetType string

const (
	UserTarget  OfferItemTargetType = "user"
	GroupTarget OfferItemTargetType = "group"
)

type OfferItemTarget struct {
	UserKey  *model.UserKey
	GroupKey *model.GroupKey
	Type     OfferItemTargetType
}

func (t OfferItemTarget) IsForGroup() bool {
	return t.Type == GroupTarget
}

func (t OfferItemTarget) IsForUser() bool {
	return t.Type == UserTarget
}

func (t OfferItemTarget) GetGroupKey() model.GroupKey {
	return *t.GroupKey
}
func (t OfferItemTarget) GetUserKey() model.UserKey {
	return *t.UserKey
}

type Offer struct {
	Key            model.OfferKey
	CreatedByKey   model.UserKey
	Status         OfferStatus
	CreatedAt      time.Time
	ExpirationTime *time.Time
	CompletedAt    *time.Time
	Message        string
}

type HistoryEntry struct {
	Timestamp         time.Time
	FromUserID        model.UserKey
	ToUserID          model.UserKey
	ResourceID        *model.ResourceKey
	TimeAmountSeconds *int64
}

func NewOffer(offerKey model.OfferKey, author model.UserKey, message string, expiration *time.Time) *Offer {

	return &Offer{
		Key:            offerKey,
		CreatedByKey:   author,
		Status:         PendingOffer,
		ExpirationTime: expiration,
		Message:        message,
		CreatedAt:      time.Now().UTC(),
	}
}

func (o *Offer) GetKey() model.OfferKey {
	return o.Key
}

func (o *Offer) GetAuthorKey() model.UserKey {
	return o.CreatedByKey
}

func (o *Offer) IsPending() bool {
	return o.Status == PendingOffer
}

type OfferItems struct {
	Items []OfferItem2
}

func (i *OfferItems) AllUserActionsCompleted() bool {
	for _, item := range i.Items {
		if !item.IsCompleted() {
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

func NewOfferItems(offerItems []OfferItem2) *OfferItems {
	copied := make([]OfferItem2, len(offerItems))
	copy(copied, offerItems)
	return &OfferItems{
		Items: copied,
	}
}

func (i *OfferItems) GetOfferItem(key model.OfferItemKey) OfferItem2 {
	for _, offerItem := range i.Items {
		if offerItem.GetKey() == key {
			return offerItem
		}
	}
	return nil
}

func (i *OfferItems) GetOfferItemsReceivedByUser(userKey model.UserKey) *OfferItems {
	var offerItems []OfferItem2
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

func (i *OfferItems) GetResourceKeys() *model.ResourceKeys {
	var resourceKeys []model.ResourceKey
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
		resourceKeys = []model.ResourceKey{}
	}
	return model.NewResourceKeys(resourceKeys)
}

func (i *OfferItems) GetUserKeys() *model.UserKeys {
	var userKeys []model.UserKey
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
	return model.NewUserKeys(userKeys)
}

func (i *OfferItems) GetGroupKeys() *model.GroupKeys {
	var groupKeys []model.GroupKey
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
	return model.NewGroupKeys(groupKeys)
}

type OfferApprovers struct {
	OfferItemsUsersCanGive    map[model.UserKey]*model.OfferItemKeys
	OfferItemsUsersCanReceive map[model.UserKey]*model.OfferItemKeys
	UsersAbleToGiveItem       map[model.OfferItemKey]*model.UserKeys
	UsersAbleToReceiveItem    map[model.OfferItemKey]*model.UserKeys
}

func (o OfferApprovers) IsUserAnApprover(userKey model.UserKey) bool {
	_, canApproveGive := o.OfferItemsUsersCanGive[userKey]
	_, canApproveReceive := o.OfferItemsUsersCanReceive[userKey]
	return canApproveGive || canApproveReceive
}
