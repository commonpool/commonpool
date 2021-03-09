package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type ConfirmResourceReturnedHandler struct {
	offerRepo domain.OfferRepository
}

func NewConfirmResourceReturnedHandler(offerRepo domain.OfferRepository) *ConfirmResourceReturnedHandler {
	return &ConfirmResourceReturnedHandler{offerRepo: offerRepo}
}

func (c *ConfirmResourceReturnedHandler) Execute(ctx context.Context, command domain.ConfirmResourceReturned) error {
	offerKey, err := keys.ParseOfferKey(command.AggregateID)
	if err != nil {
		return err
	}
	return doWithOffer(
		ctx,
		offerKey,
		c.offerRepo,
		c.confirmResourceReturned(ctx, command))
}

func (c *ConfirmResourceReturnedHandler) confirmResourceReturned(ctx context.Context, command domain.ConfirmResourceReturned) func(offer *domain.Offer) error {
	return func(offer *domain.Offer) error {
		loggedInUser, err := oidc.GetLoggedInUser(ctx)
		if err != nil {
			return err
		}
		return offer.NotifyBorrowerReturnedResource(loggedInUser.GetUserKey(), command.Payload.OfferItemKey)
	}
}
