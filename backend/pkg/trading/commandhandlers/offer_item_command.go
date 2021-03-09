package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

func doWithOffer(
	ctx context.Context,
	offerKey keys.OfferKey,
	offerRepo domain.OfferRepository,
	do func(offer *domain.Offer) error,
) error {

	offer, err := offerRepo.Load(ctx, offerKey)
	if err != nil {
		return err
	}

	if err := do(offer); err != nil {
		return err
	}

	if err := offerRepo.Save(ctx, offer); err != nil {
		return err
	}

	return nil

}
