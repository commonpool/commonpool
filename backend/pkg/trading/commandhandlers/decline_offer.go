package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type DeclineOfferHandler struct {
	offerRepo domain.OfferRepository
}

func NewDeclineOfferHandler(offerRepo domain.OfferRepository) *DeclineOfferHandler {
	return &DeclineOfferHandler{offerRepo: offerRepo}
}

func (c *DeclineOfferHandler) Execute(ctx context.Context, command domain.DeclineOffer) error {
	offerKey, err := keys.ParseOfferKey(command.AggregateID)
	if err != nil {
		return err
	}
	return doWithOffer(
		ctx,
		offerKey,
		c.offerRepo,
		c.declineOffer(ctx))
}

func (c *DeclineOfferHandler) declineOffer(ctx context.Context) func(offer *domain.Offer) error {
	return func(offer *domain.Offer) error {
		loggedInUser, err := oidc.GetLoggedInUser(ctx)
		if err != nil {
			return err
		}
		return offer.DeclineOffer(loggedInUser.GetUserKey())
	}
}
