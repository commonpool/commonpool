package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

type ConfirmServiceGivenPayload struct {
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type ConfirmServiceGiven struct {
	commands.CommandEnvelope
	Payload ConfirmServiceGivenPayload `json:"payload"`
}

func NewConfirmServiceGiven(
	ctx context.Context,
	offerKey keys.OfferKey,
	offerItemKey keys.OfferItemKey) ConfirmServiceGiven {
	return ConfirmServiceGiven{
		commands.NewCommandEnvelope(ctx, ConfirmServiceGivenCommand, "offer", offerKey.String()),
		ConfirmServiceGivenPayload{offerItemKey},
	}
}
