package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

type ConfirmResourceReturnedPayload struct {
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type ConfirmResourceReturned struct {
	commands.CommandEnvelope
	Payload ConfirmResourceReturnedPayload `json:"payload"`
}

func NewConfirmResourceReturned(
	ctx context.Context,
	offerKey keys.OfferKey,
	offerItemKey keys.OfferItemKey) ConfirmResourceReturned {
	return ConfirmResourceReturned{
		commands.NewCommandEnvelope(ctx, ConfirmResourceReturnedCommand, "offer", offerKey.String()),
		ConfirmResourceReturnedPayload{offerItemKey},
	}
}
