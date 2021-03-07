package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
)

func (t TradingService) ConfirmResourceTransferred(ctx context.Context, confirmedItemKey keys.OfferItemKey) error {

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	offerKey, err := t.getOfferKeyFromOfferItemKey.Get(nil, confirmedItemKey)
	if err != nil {
		return err
	}

	offer, err := t.offerRepo.Load(ctx, offerKey)
	if err != nil {
		return err
	}

	receivingApprovers, err := t.tradingStore.FindReceivingApproversForOfferItem(confirmedItemKey)
	if err != nil {
		return err
	}
	givingApprovers, err := t.tradingStore.FindGivingApproversForOfferItem(confirmedItemKey)
	if err != nil {
		return err
	}

	isReceiver := receivingApprovers.Contains(loggedInUserKey)
	isGiver := givingApprovers.Contains(loggedInUserKey)

	if !isGiver && !isReceiver {
		return exceptions.ErrForbidden
	}

	if isReceiver {
		if err := offer.NotifyResourceReceived(loggedInUserKey, confirmedItemKey); err != nil {
			return err
		}
	}

	if isGiver {
		if err := offer.NotifyResourceGiven(loggedInUserKey, confirmedItemKey); err != nil {
			return err
		}
	}

	return t.offerRepo.Save(ctx, offer)
}
