package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) ConfirmServiceProvided(ctx context.Context, confirmedItemKey trading.OfferItemKey) error {

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	offerItem, err := t.tradingStore.GetOfferItem(nil, confirmedItemKey)
	if err != nil {
		return err
	}

	if !offerItem.IsServiceProviding() {
		return exceptions.ErrWrongOfferItemType
	}

	serviceProvided := offerItem.(*trading.ProvideServiceItem)

	receivingApprovers, err := t.tradingStore.FindReceivingApproversForOfferItem(offerItem.GetKey())
	if err != nil {
		return err
	}

	givingApprovers, err := t.tradingStore.FindGivingApproversForOfferItem(offerItem.GetKey())
	if err != nil {
		return err
	}

	if receivingApprovers.Contains(loggedInUserKey) {
		serviceProvided.ServiceReceivedConfirmation = true
	}
	if givingApprovers.Contains(loggedInUserKey) {
		serviceProvided.ServiceGivenConfirmation = true
	}

	err = t.tradingStore.UpdateOfferItem(ctx, serviceProvided)
	if err != nil {
		return err
	}

	return t.checkIfAllItemsCompleted(ctx, loggedInUser, offerItem)

}
