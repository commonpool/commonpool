package service

import (
	"context"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/exceptions"
	trading2 "github.com/commonpool/backend/pkg/trading"
)

/**
CreditTransfer   OfferItemType = "transfer_credits"
ProvideService   OfferItemType = "provide_service"
BorrowResource   OfferItemType = "borrow_resource"
ResourceTransfer OfferItemType = "transfer_resource"
*/

func (t TradingService) ConfirmServiceProvided(ctx context.Context, confirmedItemKey model.OfferItemKey) error {

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	// retrieving item
	offerItem, err := t.tradingStore.GetOfferItem(nil, confirmedItemKey)
	if err != nil {
		return err
	}

	if !offerItem.IsServiceProviding() {
		return exceptions.ErrWrongOfferItemType
	}

	serviceProvided := offerItem.(*trading2.ProvideServiceItem)

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

func (t TradingService) ConfirmResourceTransferred(ctx context.Context, confirmedItemKey model.OfferItemKey) error {

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	// retrieving item
	offerItem, err := t.tradingStore.GetOfferItem(nil, confirmedItemKey)
	if err != nil {
		return err
	}

	if !offerItem.IsResourceTransfer() {
		return exceptions.ErrWrongOfferItemType
	}

	resourceTransfer := offerItem.(*trading2.ResourceTransferItem)

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

func (t TradingService) ConfirmResourceBorrowed(ctx context.Context, confirmedItemKey model.OfferItemKey) error {

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	// retrieving item
	offerItem, err := t.tradingStore.GetOfferItem(nil, confirmedItemKey)
	if err != nil {
		return err
	}

	if !offerItem.IsBorrowingResource() {
		return exceptions.ErrWrongOfferItemType
	}

	resourceTransfer := offerItem.(*trading2.BorrowResourceItem)

	if resourceTransfer.ItemGiven && resourceTransfer.ItemTaken {
		return nil
	}

	receivingApprovers, err := t.tradingStore.FindReceivingApproversForOfferItem(offerItem.GetKey())
	if err != nil {
		return err
	}

	givingApprovers, err := t.tradingStore.FindGivingApproversForOfferItem(offerItem.GetKey())
	if err != nil {
		return err
	}

	if receivingApprovers.Contains(loggedInUserKey) {
		resourceTransfer.ItemTaken = true
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

func (t TradingService) ConfirmBorrowedResourceReturned(ctx context.Context, confirmedItemKey model.OfferItemKey) error {

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	// retrieving item
	offerItem, err := t.tradingStore.GetOfferItem(nil, confirmedItemKey)
	if err != nil {
		return err
	}

	if !offerItem.IsBorrowingResource() {
		return exceptions.ErrWrongOfferItemType
	}

	resourceTransfer := offerItem.(*trading2.BorrowResourceItem)

	receivingApprovers, err := t.tradingStore.FindReceivingApproversForOfferItem(offerItem.GetKey())
	if err != nil {
		return err
	}

	givingApprovers, err := t.tradingStore.FindGivingApproversForOfferItem(offerItem.GetKey())
	if err != nil {
		return err
	}

	if receivingApprovers.Contains(loggedInUserKey) {
		resourceTransfer.ItemTaken = true
		resourceTransfer.ItemReturnedBack = true
	}
	if givingApprovers.Contains(loggedInUserKey) {
		resourceTransfer.ItemGiven = true
		resourceTransfer.ItemReceivedBack = true
	}

	err = t.tradingStore.UpdateOfferItem(ctx, resourceTransfer)
	if err != nil {
		return err
	}

	return t.checkIfAllItemsCompleted(ctx, loggedInUser, offerItem)

}

func (t TradingService) checkIfAllItemsCompleted(ctx context.Context, loggerInUser model.UserReference, offerItem trading2.OfferItem) error {

	offer, err := t.tradingStore.GetOffer(offerItem.GetOfferKey())
	if err != nil {
		return err
	}

	offerItems, err := t.tradingStore.GetOfferItemsForOffer(offer.Key)
	if err != nil {
		return err
	}

	approvers, err := t.tradingStore.FindApproversForOffer(offer.Key)
	if err != nil {
		return err
	}

	allUsersInOffer, err := t.us.GetByKeys(ctx, approvers.AllUserKeys().Items)
	if err != nil {
		return err
	}

	return t.checkOfferCompleted(ctx, offer.GroupKey, offer.Key, offerItems, loggerInUser, allUsersInOffer)

}
