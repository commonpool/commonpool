package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
)

func (t TradingService) DeclineOffer(ctx context.Context, offerKey keys.OfferKey) error {

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	domainOffer, err := t.offerRepo.Load(ctx, offerKey)
	if err != nil {
		return err
	}

	approvers, err := t.tradingStore.FindApproversForOffer(offerKey)
	if err != nil {
		return err
	}

	if !approvers.AllUserKeys().Contains(loggedInUserKey) {
		return exceptions.ErrForbidden
	}

	if err = domainOffer.DeclineOffer(loggedInUserKey); err != nil {
		return err
	}

	if err := t.offerRepo.Save(ctx, domainOffer); err != nil {
		return err
	}

	return nil

}
