package trading

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
)

type Approvers interface {
	GetInboundOfferItems(userKey keys.UserKey) *keys.OfferItemKeys
	GetOutboundOfferItems(userKey keys.UserKey) *keys.OfferItemKeys
	HasAnyOfferItemsToApprove(userKey keys.UserKey) bool
	GetInboundApprovers(offerItemKey keys.OfferItemKey) *keys.UserKeys
	GetOutboundApprovers(offerItemKey keys.OfferItemKey) *keys.UserKeys
	AllUserKeys() *keys.UserKeys
}

type OfferApprovers struct {
	OfferKey                  keys.OfferKey
	OfferItemsUsersCanGive    map[keys.UserKey]*keys.OfferItemKeys
	OfferItemsUsersCanReceive map[keys.UserKey]*keys.OfferItemKeys
	UsersAbleToGiveItem       map[keys.OfferItemKey]*keys.UserKeys
	UsersAbleToReceiveItem    map[keys.OfferItemKey]*keys.UserKeys
}

var _ Approvers = &OfferApprovers{}

func (o OfferApprovers) GetInboundOfferItems(userKey keys.UserKey) *keys.OfferItemKeys {
	if oi, ok := o.OfferItemsUsersCanReceive[userKey]; ok {
		return oi
	}
	return keys.NewOfferItemKeys([]keys.OfferItemKey{})
}

func (o OfferApprovers) GetOutboundOfferItems(userKey keys.UserKey) *keys.OfferItemKeys {
	if oi, ok := o.OfferItemsUsersCanGive[userKey]; ok {
		return oi
	}
	return keys.NewOfferItemKeys([]keys.OfferItemKey{})
}

func (o OfferApprovers) HasAnyOfferItemsToApprove(userKey keys.UserKey) bool {
	outboundOfferItems, hasOutbound := o.OfferItemsUsersCanGive[userKey]
	inboundOfferItems, hasInbound := o.OfferItemsUsersCanReceive[userKey]

	if !hasOutbound && !hasInbound {
		return false
	}

	if outboundOfferItems.IsEmpty() && inboundOfferItems.IsEmpty() {
		return false
	}

	return true
}

func (o OfferApprovers) GetInboundApprovers(offerItemKey keys.OfferItemKey) *keys.UserKeys {
	userKeys, ok := o.UsersAbleToReceiveItem[offerItemKey]
	if !ok {
		return keys.NewEmptyUserKeys()
	}
	return userKeys
}

func (o OfferApprovers) GetOutboundApprovers(offerItemKey keys.OfferItemKey) *keys.UserKeys {
	userKeys, ok := o.UsersAbleToGiveItem[offerItemKey]
	if !ok {
		return keys.NewEmptyUserKeys()
	}
	return userKeys
}

func (o OfferApprovers) IsUserAnApprover(userKey keys.UserKey) bool {
	_, canApproveGive := o.OfferItemsUsersCanGive[userKey]
	_, canApproveReceive := o.OfferItemsUsersCanReceive[userKey]
	return canApproveGive || canApproveReceive
}

func (o *OfferApprovers) AllUserKeys() *keys.UserKeys {
	userKeyMap := map[keys.UserKey]bool{}
	for userKey := range o.OfferItemsUsersCanGive {
		userKeyMap[userKey] = true
	}
	for userKey := range o.OfferItemsUsersCanReceive {
		userKeyMap[userKey] = true
	}
	var userKeys []keys.UserKey
	for userKey := range userKeyMap {
		userKeys = append(userKeys, userKey)
	}
	return keys.NewUserKeys(userKeys)
}

type OffersApprovers struct {
	Items []*OfferApprovers
}

func NewOffersApprovers(items []*OfferApprovers) *OffersApprovers {
	copied := make([]*OfferApprovers, len(items))
	copy(copied, items)
	return &OffersApprovers{
		Items: copied,
	}
}

func (a *OffersApprovers) GetApproversForOffer(offerKey keys.OfferKey) (*OfferApprovers, error) {
	for _, offerApprovers := range a.Items {
		if offerApprovers.OfferKey == offerKey {
			return offerApprovers, nil
		}
	}
	return nil, fmt.Errorf("could not find approvers for offer")
}

type Approvers2 struct {
	Items []*Approver
}

func NewApprovers(approvers ...*Approver) *Approvers2 {
	return &Approvers2{
		Items: approvers,
	}
}

func (a Approvers2) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Items)
}

func (a *Approvers2) UnmarshalJSON(bytes []byte) error {
	var items []*Approver
	if err := json.Unmarshal(bytes, &items); err != nil {
		return err
	}
	a.Items = items
	return nil
}
