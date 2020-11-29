package trading

import (
	"context"
	"github.com/commonpool/backend/model"
)

type Store interface {
	SaveOffer(offer Offer, items *OfferItems) error
	GetOffer(key model.OfferKey) (Offer, error)
	GetItems(key model.OfferKey) (*OfferItems, error)
	GetItem(ctx context.Context, key model.OfferItemKey) (*OfferItem, error)
	GetOffers(qry GetOffersQuery) (GetOffersResult, error)
	GetDecisions(key model.OfferKey) ([]OfferDecision, error)
	SaveDecision(key model.OfferKey, user model.UserKey, decision Decision) error
	ConfirmItemReceived(ctx context.Context, key model.OfferItemKey) error
	ConfirmItemGiven(ctx context.Context, key model.OfferItemKey) error
	SaveOfferStatus(key model.OfferKey, offer OfferStatus) error
	GetTradingHistory(ctx context.Context, ids *model.UserKeys) ([]TradingHistoryEntry, error)
}

type GetOffersQuery struct {
	ResourceKey *model.ResourceKey
	Status      *OfferStatus
	UserKeys    []model.UserKey
}

type GetOffersResult struct {
	Items []GetOffersResultItem
}

type GetOffersResultItem struct {
	Offer          Offer
	OfferItems     []OfferItem
	OfferDecisions []OfferDecision
}
