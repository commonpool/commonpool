package commands

import (
	"context"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

const (
	AcceptOfferCommand = "accept_offer"
)

type AcceptOfferCommandPayload struct {
}

type AcceptOffer struct {
	commands.CommandEnvelope
	AcceptOfferCommandPayload `json:"payload"`
}

func NewAcceptOffer(ctx context.Context, offerKey keys.OfferKey) AcceptOffer {
	return AcceptOffer{
		commands.NewCommandEnvelope(ctx, AcceptOfferCommand, "offer", offerKey.String()),
		AcceptOfferCommandPayload{},
	}
}

type AcceptOfferCommandHandler struct {
}

func (h *AcceptOfferCommandHandler) Do(ctx context.Context, command *AcceptOffer) {

}
