package trading

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type Store interface {
	SaveOffer(offer *Offer, offerItems *OfferItems) error
	GetOffer(key keys.OfferKey) (*Offer, error)
	GetOfferItemsForOffer(key keys.OfferKey) (*OfferItems, error)
	GetOfferItem(ctx context.Context, key keys.OfferItemKey) (OfferItem, error)
	GetOffersForUser(userKey keys.UserKey) (*GetOffersResult, error)
	UpdateOfferItem(ctx context.Context, offerItem OfferItem) error
	UpdateOfferStatus(key keys.OfferKey, offer OfferStatus) error
	GetTradingHistory(ctx context.Context, ids *keys.UserKeys) ([]HistoryEntry, error)
	FindApproversForOffer(offerKey keys.OfferKey) (*OfferApprovers, error)
	FindApproversForOffers(offerKeys *keys.OfferKeys) (*OffersApprovers, error)
	FindApproversForCandidateOffer(offer *Offer, offerItems *OfferItems) (*keys.UserKeys, error)
	FindReceivingApproversForOfferItem(offerItemKey keys.OfferItemKey) (*keys.UserKeys, error)
	FindGivingApproversForOfferItem(offerItemKey keys.OfferItemKey) (*keys.UserKeys, error)
	MarkOfferItemsAsAccepted(ctx context.Context, approvedBy keys.UserKey, approvedByGiver *keys.OfferItemKeys, approvedByReceiver *keys.OfferItemKeys) error
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
