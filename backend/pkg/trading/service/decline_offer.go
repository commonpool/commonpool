package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) DeclineOffer(ctx context.Context, offerKey keys.OfferKey) error {

	user, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := user.GetUserKey()

	offer, err := t.tradingStore.GetOffer(offerKey)
	if err != nil {
		return err
	}

	if offer.Status != trading.PendingOffer {
		return fmt.Errorf("could not decline a offer that is not pending")
	}

	approvers, err := t.tradingStore.FindApproversForOffer(offerKey)
	if err != nil {
		return err
	}

	if !approvers.HasAnyOfferItemsToApprove(loggedInUserKey) {
		return exceptions.ErrForbidden
	}

	err = t.tradingStore.UpdateOfferStatus(offerKey, trading.DeclinedOffer)
	if err != nil {
		return err
	}

	return nil

}
