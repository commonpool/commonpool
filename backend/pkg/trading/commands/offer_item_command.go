package commands

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	keys2 "github.com/commonpool/backend/pkg/trading/readmodels"
)

func doWithOfferItem(
	ctx context.Context,
	offerItemKey keys.OfferItemKey,
	do func(offer *domain.Offer) error) (*keys2.OfferItemReadModel2, error) {

	offerKey, err := h.getOfferKeyForOfferItem.Get(ctx, offerItemKey)
	if err != nil {
		return nil, err
	}

	offer, err := h.offerRepo.Load(ctx, offerKey)
	if err != nil {
		return nil, err
	}

	if offer.GetVersion() == 0 {
		return nil, exceptions.ErrOfferNotFound
	}

	if err := do(offer); err != nil {
		return nil, err
	}

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	var offerItem *keys2.OfferItemReadModel2
	err = retry.Do(func() error {
		var err error
		offerItem, err = h.getOfferItem.Get(ctx, offerItemKey)
		if err != nil {
			return err
		}
		if offerItem.Version != offer.GetVersion() {
			return fmt.Errorf("unexpected version: %d, expected: %d", offerItem.Version, offer.GetVersion())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return offerItem, nil

}
