package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

type DeclineOfferPayload struct {
}

type DeclineOffer struct {
	commands.CommandEnvelope
	Payload DeclineOfferPayload `json:"payload"`
}

func NewDeclineOffer(
	ctx context.Context,
	offerKey keys.OfferKey) DeclineOffer {
	return DeclineOffer{
		commands.NewCommandEnvelope(ctx, DeclineOfferCommand, "offer", offerKey.String()),
		DeclineOfferPayload{},
	}
}
