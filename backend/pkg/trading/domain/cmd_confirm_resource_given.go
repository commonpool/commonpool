package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

type ConfirmResourceGivenPayload struct {
	OfferItemKey keys.OfferItemKey `json:"offerItemKey"`
	ConfirmedBy  keys.UserKey      `json:"confirmedBy"`
}

type ConfirmResourceGiven struct {
	commands.CommandEnvelope
	Payload ConfirmResourceGivenPayload `json:"payload"`
}

func NewConfirmResourceGiven(ctx context.Context, offerKey keys.OfferKey, offerItemKey keys.OfferItemKey, by keys.UserKey) ConfirmResourceGiven {
	return ConfirmResourceGiven{
		commands.NewCommandEnvelope(ctx, ConfirmResourceGivenCommand, "offer", offerKey.String()),
		ConfirmResourceGivenPayload{offerItemKey, by},
	}
}
