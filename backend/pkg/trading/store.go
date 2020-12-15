package trading

import (
	"context"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
)

type Store interface {
	SaveOffer(offer *tradingmodel.Offer, offerItems *tradingmodel.OfferItems) error
	GetOffer(key tradingmodel.OfferKey) (*tradingmodel.Offer, error)
	GetOfferItemsForOffer(key tradingmodel.OfferKey) (*tradingmodel.OfferItems, error)
	GetOfferItem(ctx context.Context, key tradingmodel.OfferItemKey) (tradingmodel.OfferItem, error)
	GetOffersForUser(userKey usermodel.UserKey) (*GetOffersResult, error)
	UpdateOfferItem(ctx context.Context, offerItem tradingmodel.OfferItem) error
	UpdateOfferStatus(key tradingmodel.OfferKey, offer tradingmodel.OfferStatus) error
	GetTradingHistory(ctx context.Context, ids *usermodel.UserKeys) ([]tradingmodel.HistoryEntry, error)
	FindApproversForOffer(offerKey tradingmodel.OfferKey) (*tradingmodel.OfferApprovers, error)
	FindApproversForOffers(offerKeys *tradingmodel.OfferKeys) (*tradingmodel.OffersApprovers, error)
	FindApproversForCandidateOffer(offer *tradingmodel.Offer, offerItems *tradingmodel.OfferItems) (*usermodel.UserKeys, error)
	FindReceivingApproversForOfferItem(offerItemKey tradingmodel.OfferItemKey) (*usermodel.UserKeys, error)
	FindGivingApproversForOfferItem(offerItemKey tradingmodel.OfferItemKey) (*usermodel.UserKeys, error)
	MarkOfferItemsAsAccepted(ctx context.Context, approvedBy usermodel.UserKey, approvedByGiver *tradingmodel.OfferItemKeys, approvedByReceiver *tradingmodel.OfferItemKeys) error
}

type GetOffersQuery struct {
	ResourceKey *resourcemodel.ResourceKey
	Status      *tradingmodel.OfferStatus
	UserKeys    []usermodel.UserKey
}

type GetOffersResult struct {
	Items []*GetOffersResultItem
}

type GetOffersResultItem struct {
	Offer      *tradingmodel.Offer
	OfferItems *tradingmodel.OfferItems
}

func (r *GetOffersResult) GetOfferKeys() *tradingmodel.OfferKeys {
	var offerKeys []tradingmodel.OfferKey
	for _, item := range r.Items {
		offerKeys = append(offerKeys, item.Offer.Key)
	}
	if offerKeys == nil {
		offerKeys = []tradingmodel.OfferKey{}
	}
	return tradingmodel.NewOfferKeys(offerKeys)
}
