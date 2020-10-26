package trading

import "github.com/commonpool/backend/model"

type Store interface {
	SaveOffer(offer model.Offer, items []model.OfferItem) error
	GetOffer(key model.OfferKey) (model.Offer, error)
	GetItems(key model.OfferKey) ([]model.OfferItem, error)
	GetOffers(qry GetOffersQuery) (GetOffersResult, error)
	GetDecisions(key model.OfferKey) ([]model.OfferDecision, error)
	SaveDecision(key model.OfferKey, user model.UserKey, decision model.Decision) error
	CompleteOffer(key model.OfferKey, status model.OfferStatus) error
}

type GetOffersQuery struct {
	ResourceKey *model.ResourceKey
	Status      *model.OfferStatus
	UserKeys    []model.UserKey
}

type GetOffersResult struct {
	Items []GetOffersResultItem
}

type GetOffersResultItem struct {
	Offer          model.Offer
	OfferItems     []model.OfferItem
	OfferDecisions []model.OfferDecision
}
