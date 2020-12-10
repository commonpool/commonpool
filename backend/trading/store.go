package trading

import (
	"context"
	"github.com/commonpool/backend/model"
)

type Store interface {
	SaveOffer(offer *Offer, offerItems *OfferItems) error
	GetOffer(key model.OfferKey) (*Offer, error)
	GetOfferItemsForOffer(key model.OfferKey) (*OfferItems, error)
	GetOfferItem(ctx context.Context, key model.OfferItemKey) (OfferItem, error)
	GetOffersForUser(userKey model.UserKey) (*GetOffersResult, error)
	UpdateOfferItem(ctx context.Context, offerItem OfferItem) error
	UpdateOfferStatus(key model.OfferKey, offer OfferStatus) error
	GetTradingHistory(ctx context.Context, ids *model.UserKeys) ([]HistoryEntry, error)
	FindApproversForOffer(offerKey model.OfferKey) (*OfferApprovers, error)
	FindApproversForOffers(offerKeys *model.OfferKeys) (*OffersApprovers, error)
	FindApproversForCandidateOffer(offer *Offer, offerItems *OfferItems) (*model.UserKeys, error)
	FindReceivingApproversForOfferItem(offerItemKey model.OfferItemKey) (*model.UserKeys, error)
	FindGivingApproversForOfferItem(offerItemKey model.OfferItemKey) (*model.UserKeys, error)
	MarkOfferItemsAsAccepted(ctx context.Context, approvedBy model.UserKey, approvedByGiver *model.OfferItemKeys, approvedByReceiver *model.OfferItemKeys) error
}

type GetOffersQuery struct {
	ResourceKey *model.ResourceKey
	Status      *OfferStatus
	UserKeys    []model.UserKey
}

type GetOffersResult struct {
	Items []*GetOffersResultItem
}

type GetOffersResultItem struct {
	Offer      *Offer
	OfferItems *OfferItems
}

func (r *GetOffersResult) GetOfferKeys() *model.OfferKeys {
	var offerKeys []model.OfferKey
	for _, item := range r.Items {
		offerKeys = append(offerKeys, item.Offer.Key)
	}
	if offerKeys == nil {
		offerKeys = []model.OfferKey{}
	}
	return model.NewOfferKeys(offerKeys)
}
