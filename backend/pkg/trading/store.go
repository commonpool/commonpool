package trading

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type Store interface {
	SaveOffer(offer *Offer, offerItems *OfferItems) error
	GetOffer(key OfferKey) (*Offer, error)
	GetOfferItemsForOffer(key OfferKey) (*OfferItems, error)
	GetOfferItem(ctx context.Context, key OfferItemKey) (OfferItem, error)
	GetOffersForUser(userKey keys.UserKey) (*GetOffersResult, error)
	UpdateOfferItem(ctx context.Context, offerItem OfferItem) error
	UpdateOfferStatus(key OfferKey, offer OfferStatus) error
	GetTradingHistory(ctx context.Context, ids *keys.UserKeys) ([]HistoryEntry, error)
	FindApproversForOffer(offerKey OfferKey) (*OfferApprovers, error)
	FindApproversForOffers(offerKeys *OfferKeys) (*OffersApprovers, error)
	FindApproversForCandidateOffer(offer *Offer, offerItems *OfferItems) (*keys.UserKeys, error)
	FindReceivingApproversForOfferItem(offerItemKey OfferItemKey) (*keys.UserKeys, error)
	FindGivingApproversForOfferItem(offerItemKey OfferItemKey) (*keys.UserKeys, error)
	MarkOfferItemsAsAccepted(ctx context.Context, approvedBy keys.UserKey, approvedByGiver *OfferItemKeys, approvedByReceiver *OfferItemKeys) error
}

type GetOffersQuery struct {
	ResourceKey *keys.ResourceKey
	Status      *OfferStatus
	UserKeys    []keys.UserKey
}

type GetOffersResult struct {
	Items []*GetOffersResultItem
}

type GetOffersResultItem struct {
	Offer      *Offer
	OfferItems *OfferItems
}

func (r *GetOffersResult) GetOfferKeys() *OfferKeys {
	var offerKeys []OfferKey
	for _, item := range r.Items {
		offerKeys = append(offerKeys, item.Offer.Key)
	}
	if offerKeys == nil {
		offerKeys = []OfferKey{}
	}
	return NewOfferKeys(offerKeys)
}
