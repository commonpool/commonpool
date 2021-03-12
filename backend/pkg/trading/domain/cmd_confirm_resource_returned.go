package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

type ConfirmResourceReturnedPayload struct {
	OfferItemKey keys.OfferItemKey `json:"offerItemKey"`
	ConfirmedBy  keys.UserKey      `json:"confirmedBy"`
}

type ConfirmResourceReturned struct {
	commands.CommandEnvelope
	Payload ConfirmResourceReturnedPayload `json:"payload"`
}

func NewConfirmResourceReturned(ctx context.Context, offerKey keys.OfferKey, offerItemKey keys.OfferItemKey, confirmedBy keys.UserKey) ConfirmResourceReturned {
	return ConfirmResourceReturned{
		commands.NewCommandEnvelope(ctx, ConfirmResourceReturnedCommand, "offer", offerKey.String()),
		ConfirmResourceReturnedPayload{offerItemKey, confirmedBy},
	}
}
