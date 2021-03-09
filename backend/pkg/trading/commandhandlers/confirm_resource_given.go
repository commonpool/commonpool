package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type ConfirmResourceGivenHandler struct {
	offerRepo domain.OfferRepository
}

func NewConfirmResourceGivenHandler(offerRepo domain.OfferRepository) *ConfirmResourceGivenHandler {
	return &ConfirmResourceGivenHandler{offerRepo: offerRepo}
}

func (c *ConfirmResourceGivenHandler) Execute(ctx context.Context, command domain.ConfirmResourceGiven) error {
	offerKey, err := keys.ParseOfferKey(command.AggregateID)
	if err != nil {
		return err
	}
	return doWithOffer(
		ctx,
		offerKey,
		c.offerRepo,
		c.confirmResourceGiven(ctx, command))
}

func (c *ConfirmResourceGivenHandler) confirmResourceGiven(ctx context.Context, command domain.ConfirmResourceGiven) func(offer *domain.Offer) error {
	return func(offer *domain.Offer) error {
		loggedInUser, err := oidc.GetLoggedInUser(ctx)
		if err != nil {
			return err
		}
		return offer.NotifyBorrowerReturnedResource(loggedInUser.GetUserKey(), command.Payload.OfferItemKey)
	}
}
