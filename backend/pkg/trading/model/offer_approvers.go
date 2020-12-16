package model

import (
	"fmt"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

type OfferApprovers struct {
	OfferKey                  OfferKey
	OfferItemsUsersCanGive    map[usermodel.UserKey]*OfferItemKeys
	OfferItemsUsersCanReceive map[usermodel.UserKey]*OfferItemKeys
	UsersAbleToGiveItem       map[OfferItemKey]*usermodel.UserKeys
	UsersAbleToReceiveItem    map[OfferItemKey]*usermodel.UserKeys
}

func (o OfferApprovers) IsUserAnApprover(userKey usermodel.UserKey) bool {
	_, canApproveGive := o.OfferItemsUsersCanGive[userKey]
	_, canApproveReceive := o.OfferItemsUsersCanReceive[userKey]
	return canApproveGive || canApproveReceive
}

func (o *OfferApprovers) AllUserKeys() *usermodel.UserKeys {
	userKeyMap := map[usermodel.UserKey]bool{}
	for userKey := range o.OfferItemsUsersCanGive {
		userKeyMap[userKey] = true
	}
	for userKey := range o.OfferItemsUsersCanReceive {
		userKeyMap[userKey] = true
	}
	var userKeys []usermodel.UserKey
	for userKey := range userKeyMap {
		userKeys = append(userKeys, userKey)
	}
	return usermodel.NewUserKeys(userKeys)
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
