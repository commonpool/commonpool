package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

type ConfirmResourceBorrowedPayload struct {
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}

type ConfirmResourceBorrowed struct {
	commands.CommandEnvelope
	Payload ConfirmResourceBorrowedPayload `json:"payload"`
}

func NewConfirmResourceBorrowed(
	ctx context.Context,
	offerKey keys.OfferKey,
	offerItemKey keys.OfferItemKey) ConfirmResourceBorrowed {
	return ConfirmResourceBorrowed{
		commands.NewCommandEnvelope(ctx, ConfirmResourceBorrowedCommand, "offer", offerKey.String()),
		ConfirmResourceBorrowedPayload{offerItemKey},
	}
}
