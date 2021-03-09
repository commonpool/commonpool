package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

type ConfirmResourceGivenPayload struct {
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type ConfirmResourceGiven struct {
	commands.CommandEnvelope
	Payload ConfirmResourceGivenPayload `json:"payload"`
}

func NewConfirmResourceGiven(
	ctx context.Context,
	offerKey keys.OfferKey,
	offerItemKey keys.OfferItemKey) ConfirmResourceGiven {
	return ConfirmResourceGiven{
		commands.NewCommandEnvelope(ctx, ConfirmResourceGivenCommand, "offer", offerKey.String()),
		ConfirmResourceGivenPayload{offerItemKey},
	}
}
