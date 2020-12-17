package trading

import (
	ctx "context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
	"golang.org/x/net/context"
)

type Service interface {
	GetOfferItem(ctx context.Context, offerItemKey OfferItemKey) (OfferItem, error)
	ConfirmServiceProvided(ctx context.Context, offerItemKey OfferItemKey) error
	ConfirmResourceTransferred(ctx context.Context, confirmedItemKey OfferItemKey) error
	ConfirmResourceBorrowed(ctx context.Context, confirmedItemKey OfferItemKey) error
	ConfirmBorrowedResourceReturned(ctx context.Context, confirmedItemKey OfferItemKey) error
	AcceptOffer(ctx ctx.Context, offerKey OfferKey) error
	DeclineOffer(ctx ctx.Context, offerKey OfferKey) error
	GetTradingHistory(ctx context.Context, userIDs *keys.UserKeys) ([]HistoryEntry, error)
	SendOffer(ctx context.Context, groupKey keys.GroupKey, offerItems *OfferItems, message string) (*Offer, *OfferItems, error)
	FindTargetsForOfferItem(ctx ctx.Context, groupKey keys.GroupKey, itemType OfferItemType, from *resource.Target, to *resource.Target) (*resource.Targets, error)
	GetOffersForUser(key keys.UserKey) (*GetOffersResult, error)
	FindApproversForOffers(offers *OfferKeys) (*OffersApprovers, error)
	GetOffer(offerKey OfferKey) (*Offer, error)
	GetOfferItemsForOffer(offerKey OfferKey) (*OfferItems, error)
	FindApproversForOffer(offerKey OfferKey) (*OfferApprovers, error)
}
