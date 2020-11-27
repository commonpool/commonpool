package trading

import (
	ctx "context"
	"github.com/commonpool/backend/model"
)

type Service interface {
	ConfirmItemReceived(ctx ctx.Context, request *ConfirmItemReceived) (*ConfirmItemReceivedResponse, error)
	ConfirmItemGiven(ctx ctx.Context, request *ConfirmItemGiven) (*ConfirmItemGivenResponse, error)
	AcceptOffer(ctx ctx.Context, request *AcceptOffer) (*AcceptOfferResponse, error)
}

type ConfirmItemReceived struct {
	ReceivedByUser model.UserKey
	OfferItemKey   model.OfferItemKey
}

type ConfirmItemReceivedResponse struct {
}

func NewConfirmItemReceived(receivedByUser model.UserKey, offerItemKey model.OfferItemKey) *ConfirmItemReceived {
	return &ConfirmItemReceived{
		ReceivedByUser: receivedByUser,
		OfferItemKey:   offerItemKey,
	}
}

type ConfirmItemGiven struct {
	GivenByUser  model.UserKey
	OfferItemKey model.OfferItemKey
}

type ConfirmItemGivenResponse struct {
}

func NewConfirmItemGiven(GivenByUser model.UserKey, offerItemKey model.OfferItemKey) *ConfirmItemGiven {
	return &ConfirmItemGiven{
		GivenByUser:  GivenByUser,
		OfferItemKey: offerItemKey,
	}
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
