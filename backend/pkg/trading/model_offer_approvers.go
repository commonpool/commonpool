package trading

import (
	"fmt"
	"github.com/commonpool/backend/model"
)

type OfferApprovers struct {
	OfferKey                  model.OfferKey
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

func (o *OfferApprovers) AllUserKeys() *model.UserKeys {
	userKeyMap := map[model.UserKey]bool{}
	for userKey := range o.OfferItemsUsersCanGive {
		userKeyMap[userKey] = true
	}
	for userKey := range o.OfferItemsUsersCanReceive {
		userKeyMap[userKey] = true
	}
	var userKeys []model.UserKey
	for userKey := range userKeyMap {
		userKeys = append(userKeys, userKey)
	}
	return model.NewUserKeys(userKeys)
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

func (a *OffersApprovers) GetApproversForOffer(offerKey model.OfferKey) (*OfferApprovers, error) {
	for _, offerApprovers := range a.Items {
		if offerApprovers.OfferKey == offerKey {
			return offerApprovers, nil
		}
	}
	return nil, fmt.Errorf("could not find approvers for offer")
}
