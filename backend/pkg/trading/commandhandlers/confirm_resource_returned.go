package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/pkg/trading/queries"
)

type ConfirmResourceReturnedHandler struct {
	offerRepo           domain.OfferRepository
	getOfferPermissions *queries.GetOfferPermissions
}

func NewConfirmResourceReturnedHandler(offerRepo domain.OfferRepository, getOfferPermissions *queries.GetOfferPermissions) *ConfirmResourceReturnedHandler {
	return &ConfirmResourceReturnedHandler{
		offerRepo:           offerRepo,
		getOfferPermissions: getOfferPermissions,
	}
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
		c.confirmResourceReturned(ctx, command.Payload.OfferItemKey, command.Payload.ConfirmedBy))
}

func (c *ConfirmResourceReturnedHandler) confirmResourceReturned(ctx context.Context, offerItemKey keys.OfferItemKey, confirmedBy keys.UserKey) func(offer *domain.Offer) error {
	return func(offer *domain.Offer) error {
		offerPermissions, err := c.getOfferPermissions.Get(ctx, offer.GetKey())
		if err != nil {
			return err
		}
		return offer.ConfirmResourceReturned(confirmedBy, offerItemKey, offerPermissions)
	}
}
