package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/pkg/trading/queries"
)

type AcceptOfferHandler struct {
	offerRepo     domain.OfferRepository
	getPermission *queries.GetOfferPermissions
}

func NewAcceptOfferHandler(
	offerRepo domain.OfferRepository,
	getPermission *queries.GetOfferPermissions) *AcceptOfferHandler {
	return &AcceptOfferHandler{offerRepo: offerRepo, getPermission: getPermission}
}

func (c *AcceptOfferHandler) Execute(ctx context.Context, command domain.AcceptOffer) error {

	offerKey, err := keys.ParseOfferKey(command.AggregateID)
	if err != nil {
		return err
	}

	return doWithOffer(ctx, offerKey, c.offerRepo, func(offer *domain.Offer) error {
		return c.approve(ctx, offer)
	})

}

func (c *AcceptOfferHandler) approve(ctx context.Context, offer *domain.Offer) error {

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	offerPermissions, err := c.getPermission.Get(ctx, offer.GetKey())
	if err != nil {
		return err
	}

	return offer.ApproveAll(loggedInUser.GetUserKey(), offerPermissions)

}
