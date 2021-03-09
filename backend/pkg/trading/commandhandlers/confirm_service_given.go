package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type ConfirmServiceGivenHandler struct {
	offerRepo domain.OfferRepository
}

func NewConfirmServiceGivenHandler(offerRepo domain.OfferRepository) *ConfirmServiceGivenHandler {
	return &ConfirmServiceGivenHandler{offerRepo: offerRepo}
}

func (c *ConfirmServiceGivenHandler) Execute(ctx context.Context, command domain.ConfirmServiceGiven) error {
	offerKey, err := keys.ParseOfferKey(command.AggregateID)
	if err != nil {
		return err
	}
	return doWithOffer(
		ctx,
		offerKey,
		c.offerRepo,
		c.confirmServiceGiven(ctx, command))
}

func (c *ConfirmServiceGivenHandler) confirmServiceGiven(ctx context.Context, command domain.ConfirmServiceGiven) func(offer *domain.Offer) error {
	return func(offer *domain.Offer) error {
		loggedInUser, err := oidc.GetLoggedInUser(ctx)
		if err != nil {
			return err
		}
		return offer.NotifyBorrowerReturnedResource(loggedInUser.GetUserKey(), command.Payload.OfferItemKey)
	}
}
