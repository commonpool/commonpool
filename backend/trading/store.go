package trading

import (
	"context"
	"github.com/commonpool/backend/model"
)

type Store interface {
	SaveOffer(offer *Offer, offerItems *OfferItems) error
	GetOffer(key model.OfferKey) (*Offer, error)
	GetOfferItemsForOffer(key model.OfferKey) (*OfferItems, error)
	GetOfferItem(ctx context.Context, key model.OfferItemKey) (OfferItem2, error)
	GetOffersForUser(userKey model.UserKey) (*GetOffersResult, error)
	UpdateOfferItem(ctx context.Context, offerItem OfferItem2) error
	SaveOfferStatus(key model.OfferKey, offer OfferStatus) error
	GetTradingHistory(ctx context.Context, ids *model.UserKeys) ([]HistoryEntry, error)
	FindApproversForOffer(offerKey model.OfferKey) (*OfferApprovers, error)
	FindReceivingApproversForOfferItem(offerItemKey model.OfferItemKey) (*model.UserKeys, error)
	FindGivingApproversForOfferItem(offerItemKey model.OfferItemKey) (*model.UserKeys, error)
	MarkOfferItemsAsAccepted(ctx context.Context, approvedByGiver *model.OfferItemKeys, approvedByReceiver *model.OfferItemKeys) error
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
