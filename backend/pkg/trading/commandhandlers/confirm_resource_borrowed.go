package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type ConfirmResourceBorrowedHandler struct {
	offerRepo domain.OfferRepository
}

func NewConfirmResourceBorrowedHandler(offerRepo domain.OfferRepository) *ConfirmResourceBorrowedHandler {
	return &ConfirmResourceBorrowedHandler{offerRepo: offerRepo}
}

func (c *ConfirmResourceBorrowedHandler) Execute(ctx context.Context, command domain.ConfirmResourceBorrowed) error {
	offerKey, err := keys.ParseOfferKey(command.AggregateID)
	if err != nil {
		return err
	}

	return doWithOffer(ctx, offerKey, c.offerRepo, func(offer *domain.Offer) error {
		return c.confirmResourceReturned(ctx, offer, command.Payload.OfferItemKey)
	})
}

func (c *ConfirmResourceBorrowedHandler) confirmResourceReturned(ctx context.Context, offer *domain.Offer, offerItemKey keys.OfferItemKey) error {

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	return offer.NotifyBorrowerReturnedResource(loggedInUser.GetUserKey(), offerItemKey)
}
