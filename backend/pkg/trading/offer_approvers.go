package trading

import (
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
)

type OfferApprovers struct {
	OfferKey                  OfferKey
	OfferItemsUsersCanGive    map[keys.UserKey]*OfferItemKeys
	OfferItemsUsersCanReceive map[keys.UserKey]*OfferItemKeys
	UsersAbleToGiveItem       map[OfferItemKey]*keys.UserKeys
	UsersAbleToReceiveItem    map[OfferItemKey]*keys.UserKeys
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

func (a *OffersApprovers) GetApproversForOffer(offerKey OfferKey) (*OfferApprovers, error) {
	for _, offerApprovers := range a.Items {
		if offerApprovers.OfferKey == offerKey {
			return offerApprovers, nil
		}
	}
	return nil, fmt.Errorf("could not find approvers for offer")
}
