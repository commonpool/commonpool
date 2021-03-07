package trading

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	tradingdomain "github.com/commonpool/backend/pkg/trading/domain"
)

type Store interface {
	GetOffer(key keys.OfferKey) (*Offer, error)
	GetOfferItemsForOffer(key keys.OfferKey) (*tradingdomain.OfferItems, error)
	GetOfferItem(ctx context.Context, key keys.OfferItemKey) (tradingdomain.OfferItem, error)
	GetOffersForUser(userKey keys.UserKey) (*GetOffersResult, error)
	UpdateOfferStatus(key keys.OfferKey, offer tradingdomain.OfferStatus) error
	GetTradingHistory(ctx context.Context, ids *keys.UserKeys) ([]HistoryEntry, error)
	FindApproversForOffer(offerKey keys.OfferKey) (Approvers, error)
	FindApproversForOffers(offerKeys *keys.OfferKeys) (*OffersApprovers, error)
	FindApproversForCandidateOffer(offer *Offer, offerItems *tradingdomain.OfferItems) (*keys.UserKeys, error)
	FindReceivingApproversForOfferItem(offerItemKey keys.OfferItemKey) (*keys.UserKeys, error)
	FindGivingApproversForOfferItem(offerItemKey keys.OfferItemKey) (*keys.UserKeys, error)
}

type GetOffersQuery struct {
	ResourceKey *keys.ResourceKey
	Status      *tradingdomain.OfferStatus
	UserKeys    []keys.UserKey
}

type GetOffersResult struct {
	Items []*GetOffersResultItem
}

type GetOffersResultItem struct {
	Offer      *Offer
	OfferItems *tradingdomain.OfferItems
}

func (r *GetOffersResult) GetOfferKeys() *keys.OfferKeys {
	var offerKeys []keys.OfferKey
	for _, item := range r.Items {
		offerKeys = append(offerKeys, item.Offer.Key)
	}
	if offerKeys == nil {
		offerKeys = []keys.OfferKey{}
	}
	return keys.NewOfferKeys(offerKeys)
}
