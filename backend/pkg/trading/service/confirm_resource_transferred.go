package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/trading/model"
)

func (t TradingService) ConfirmResourceTransferred(ctx context.Context, confirmedItemKey model.OfferItemKey) error {

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	offerItem, err := t.tradingStore.GetOfferItem(nil, confirmedItemKey)
	if err != nil {
		return err
	}

	if !offerItem.IsResourceTransfer() {
		return exceptions.ErrWrongOfferItemType
	}

	resourceTransfer := offerItem.(*model.ResourceTransferItem)

	receivingApprovers, err := t.tradingStore.FindReceivingApproversForOfferItem(offerItem.GetKey())
	if err != nil {
		return err
	}

	givingApprovers, err := t.tradingStore.FindGivingApproversForOfferItem(offerItem.GetKey())
	if err != nil {
		return err
	}

	if receivingApprovers.Contains(loggedInUserKey) {
		resourceTransfer.ItemReceived = true
	}
	if givingApprovers.Contains(loggedInUserKey) {
		resourceTransfer.ItemGiven = true
	}

	err = t.tradingStore.UpdateOfferItem(ctx, resourceTransfer)
	if err != nil {
		return err
	}

	return t.checkIfAllItemsCompleted(ctx, loggedInUser, offerItem)

}
