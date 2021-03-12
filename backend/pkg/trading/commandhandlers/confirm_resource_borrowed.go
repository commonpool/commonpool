package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/pkg/trading/queries"
)

type ConfirmResourceBorrowedHandler struct {
	offerRepo           domain.OfferRepository
	getOfferPermissions *queries.GetOfferPermissions
}

func NewConfirmResourceBorrowedHandler(offerRepo domain.OfferRepository, getOfferPermissions *queries.GetOfferPermissions) *ConfirmResourceBorrowedHandler {
	return &ConfirmResourceBorrowedHandler{
		offerRepo:           offerRepo,
		getOfferPermissions: getOfferPermissions,
	}
}

func (c *ConfirmResourceBorrowedHandler) Execute(ctx context.Context, command domain.ConfirmResourceBorrowed) error {
	offerKey, err := keys.ParseOfferKey(command.AggregateID)
	if err != nil {
		return err
	}

	return doWithOffer(ctx, offerKey, c.offerRepo, func(offer *domain.Offer) error {
		return c.confirmResourceReturned(ctx, offer, command.Payload.OfferItemKey, command.Payload.ConfirmedBy)
	})
}

func (c *ConfirmResourceBorrowedHandler) confirmResourceReturned(ctx context.Context, offer *domain.Offer, offerItemKey keys.OfferItemKey, by keys.UserKey) error {
	offerPermissions, err := c.getOfferPermissions.Get(ctx, offer.GetKey())
	if err != nil {
		return err
	}
	return offer.ConfirmResourceBorrowed(by, offerItemKey, offerPermissions)
}
