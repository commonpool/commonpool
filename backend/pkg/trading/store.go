package trading

import (
	"context"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

type Store interface {
	SaveOffer(offer *Offer, offerItems *OfferItems) error
	GetOffer(key OfferKey) (*Offer, error)
	GetOfferItemsForOffer(key OfferKey) (*OfferItems, error)
	GetOfferItem(ctx context.Context, key OfferItemKey) (OfferItem, error)
	GetOffersForUser(userKey usermodel.UserKey) (*GetOffersResult, error)
	UpdateOfferItem(ctx context.Context, offerItem OfferItem) error
	UpdateOfferStatus(key OfferKey, offer OfferStatus) error
	GetTradingHistory(ctx context.Context, ids *usermodel.UserKeys) ([]HistoryEntry, error)
	FindApproversForOffer(offerKey OfferKey) (*OfferApprovers, error)
	FindApproversForOffers(offerKeys *OfferKeys) (*OffersApprovers, error)
	FindApproversForCandidateOffer(offer *Offer, offerItems *OfferItems) (*usermodel.UserKeys, error)
	FindReceivingApproversForOfferItem(offerItemKey OfferItemKey) (*usermodel.UserKeys, error)
	FindGivingApproversForOfferItem(offerItemKey OfferItemKey) (*usermodel.UserKeys, error)
	MarkOfferItemsAsAccepted(ctx context.Context, approvedBy usermodel.UserKey, approvedByGiver *OfferItemKeys, approvedByReceiver *OfferItemKeys) error
}

type GetOffersQuery struct {
	ResourceKey *resourcemodel.ResourceKey
	Status      *OfferStatus
	UserKeys    []usermodel.UserKey
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
