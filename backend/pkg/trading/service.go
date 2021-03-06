package trading

import (
	ctx "context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"golang.org/x/net/context"
)

type Service interface {
	GetOfferItem(ctx context.Context, offerItemKey keys.OfferItemKey) (domain.OfferItem, error)
	ConfirmServiceProvided(ctx context.Context, offerItemKey keys.OfferItemKey) error
	ConfirmResourceTransferred(ctx context.Context, confirmedItemKey keys.OfferItemKey) error
	ConfirmResourceBorrowed(ctx context.Context, confirmedItemKey keys.OfferItemKey) error
	ConfirmBorrowedResourceReturned(ctx context.Context, confirmedItemKey keys.OfferItemKey) error
	AcceptOffer(ctx ctx.Context, offerKey keys.OfferKey) error
	DeclineOffer(ctx ctx.Context, offerKey keys.OfferKey) error
	GetTradingHistory(ctx context.Context, userIDs *keys.UserKeys) ([]HistoryEntry, error)
	SendOffer(ctx context.Context, groupKey keys.GroupKey, offerItems *domain.OfferItems, message string) (*Offer, *domain.OfferItems, error)
	FindTargetsForOfferItem(ctx ctx.Context, groupKey keys.GroupKey, itemType domain.OfferItemType, from *domain.Target, to *domain.Target) (*domain.Targets, error)
	GetOffersForUser(key keys.UserKey) (*GetOffersResult, error)
	FindApproversForOffers(offers *keys.OfferKeys) (*OffersApprovers, error)
	GetOffer(offerKey keys.OfferKey) (*Offer, error)
	GetOfferItemsForOffer(offerKey keys.OfferKey) (*domain.OfferItems, error)
	FindApproversForOffer(offerKey keys.OfferKey) (Approvers, error)
}
