package trading

import (
	ctx "context"
	"github.com/commonpool/backend/model"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"golang.org/x/net/context"
)

type Service interface {
	GetOfferItem(ctx context.Context, offerItemKey tradingmodel.OfferItemKey) (tradingmodel.OfferItem, error)
	ConfirmServiceProvided(ctx context.Context, offerItemKey tradingmodel.OfferItemKey) error
	ConfirmResourceTransferred(ctx context.Context, confirmedItemKey tradingmodel.OfferItemKey) error
	ConfirmResourceBorrowed(ctx context.Context, confirmedItemKey tradingmodel.OfferItemKey) error
	ConfirmBorrowedResourceReturned(ctx context.Context, confirmedItemKey tradingmodel.OfferItemKey) error
	AcceptOffer(ctx ctx.Context, offerKey tradingmodel.OfferKey) error
	DeclineOffer(ctx ctx.Context, offerKey tradingmodel.OfferKey) error
	GetTradingHistory(ctx context.Context, userIDs *usermodel.UserKeys) ([]tradingmodel.HistoryEntry, error)
	SendOffer(ctx context.Context, groupKey groupmodel.GroupKey, offerItems *tradingmodel.OfferItems, message string) (*tradingmodel.Offer, *tradingmodel.OfferItems, error)
	FindTargetsForOfferItem(ctx ctx.Context, groupKey groupmodel.GroupKey, itemType tradingmodel.OfferItemType, from *model.Target, to *model.Target) (*model.Targets, error)
}
