package trading

import (
	ctx "context"
	"github.com/commonpool/backend/model"
	"golang.org/x/net/context"
)

type Service interface {
	GetOfferItem(ctx context.Context, offerItemKey model.OfferItemKey) (OfferItem, error)
	ConfirmServiceProvided(ctx context.Context, offerItemKey model.OfferItemKey) error
	ConfirmResourceTransferred(ctx context.Context, confirmedItemKey model.OfferItemKey) error
	ConfirmResourceBorrowed(ctx context.Context, confirmedItemKey model.OfferItemKey) error
	ConfirmBorrowedResourceReturned(ctx context.Context, confirmedItemKey model.OfferItemKey) error
	AcceptOffer(ctx ctx.Context, offerKey model.OfferKey) error
	DeclineOffer(ctx ctx.Context, offerKey model.OfferKey) error
	GetTradingHistory(ctx context.Context, userIDs *model.UserKeys) ([]HistoryEntry, error)
	SendOffer(ctx context.Context, groupKey model.GroupKey, offerItems *OfferItems, message string) (*Offer, *OfferItems, error)
	FindTargetsForOfferItem(ctx ctx.Context, groupKey model.GroupKey, itemType OfferItemType, from *model.Target, to *model.Target) (*model.Targets, error)
}
