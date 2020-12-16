package trading

import (
	ctx "context"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/resource/model"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
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
	GetTradingHistory(ctx context.Context, userIDs *usermodel.UserKeys) ([]HistoryEntry, error)
	SendOffer(ctx context.Context, groupKey group.GroupKey, offerItems *OfferItems, message string) (*Offer, *OfferItems, error)
	FindTargetsForOfferItem(ctx ctx.Context, groupKey group.GroupKey, itemType OfferItemType, from *model.Target, to *model.Target) (*model.Targets, error)
	GetOffersForUser(key usermodel.UserKey) (*GetOffersResult, error)
	FindApproversForOffers(offers *OfferKeys) (*OffersApprovers, error)
	GetOffer(offerKey OfferKey) (*Offer, error)
	GetOfferItemsForOffer(offerKey OfferKey) (*OfferItems, error)
	FindApproversForOffer(offerKey OfferKey) (*OfferApprovers, error)
}
