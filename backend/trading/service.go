package trading

import (
	ctx "context"
	"github.com/commonpool/backend/model"
	"golang.org/x/net/context"
)

type Service interface {
	GetOfferItem(ctx context.Context, offerItemKey model.OfferItemKey) (*OfferItem, error)
	ConfirmItemReceivedOrGiven(ctx context.Context, offerItemKey model.OfferItemKey) error
	AcceptOffer(ctx ctx.Context, request *AcceptOffer) (*AcceptOfferResponse, error)
	GetTradingHistory(ctx context.Context, userIDs *model.UserKeys) ([]HistoryEntry, error)
	SendOffer(ctx context.Context, offerItems *OfferItems, message string) (*Offer, *OfferItems, *OfferDecisions, error)
}

type AcceptOffer struct {
	OfferKey model.OfferKey
}

type AcceptOfferResponse struct {
}

func NewAcceptOffer(offerKey model.OfferKey) *AcceptOffer {
	return &AcceptOffer{
		OfferKey: offerKey,
	}
}
