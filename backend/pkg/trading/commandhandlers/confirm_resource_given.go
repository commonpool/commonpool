package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/pkg/trading/queries"
)

type ConfirmResourceGivenHandler struct {
	offerRepo           domain.OfferRepository
	getOfferPermissions *queries.GetOfferPermissions
}

func NewConfirmResourceGivenHandler(offerRepo domain.OfferRepository, getOfferPermissions *queries.GetOfferPermissions) *ConfirmResourceGivenHandler {
	return &ConfirmResourceGivenHandler{
		offerRepo:           offerRepo,
		getOfferPermissions: getOfferPermissions,
	}
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
		c.confirmResourceGiven(ctx, command.Payload.OfferItemKey, command.Payload.ConfirmedBy))
}

func (c *ConfirmResourceGivenHandler) confirmResourceGiven(ctx context.Context, offerItemKey keys.OfferItemKey, by keys.UserKey) func(offer *domain.Offer) error {
	return func(offer *domain.Offer) error {
		offerPermissions, err := c.getOfferPermissions.Get(ctx, offer.GetKey())
		if err != nil {
			return err
		}
		return offer.ConfirmResourceReturned(by, offerItemKey, offerPermissions)
	}
}
