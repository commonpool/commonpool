package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

type AcceptOfferPayload struct {
}

type AcceptOffer struct {
	commands.CommandEnvelope
	Payload AcceptOfferPayload `json:"payload"`
}

func NewAcceptOffer(
	ctx context.Context,
	offerKey keys.OfferKey) AcceptOffer {
	return AcceptOffer{
		commands.NewCommandEnvelope(ctx, AcceptOfferCommand, "offer", offerKey.String()),
		AcceptOfferPayload{},
	}
}
