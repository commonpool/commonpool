package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/pkg/trading/queries"
)

type ConfirmServiceGivenHandler struct {
	offerRepo           domain.OfferRepository
	getOfferPermissions *queries.GetOfferPermissions
}

func NewConfirmServiceGivenHandler(offerRepo domain.OfferRepository, getOfferPermissions *queries.GetOfferPermissions) *ConfirmServiceGivenHandler {
	return &ConfirmServiceGivenHandler{
		offerRepo:           offerRepo,
		getOfferPermissions: getOfferPermissions,
	}
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
		c.confirmServiceGiven(ctx, command.Payload.OfferItemKey, command.Payload.ConfirmedBy))
}

func (c *ConfirmServiceGivenHandler) confirmServiceGiven(ctx context.Context, offerItemKey keys.OfferItemKey, confirmedBy keys.UserKey) func(offer *domain.Offer) error {
	return func(offer *domain.Offer) error {
		offerPermissions, err := c.getOfferPermissions.Get(ctx, offer.GetKey())
		if err != nil {
			return err
		}
		return offer.ConfirmServiceGiven(confirmedBy, offerItemKey, offerPermissions)
	}
}
