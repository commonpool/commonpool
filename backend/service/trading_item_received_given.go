package service

import (
	"context"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/trading"
	"go.uber.org/zap"
)

// ConfirmItemReceivedOrGiven will send a notification to concerned users that an offer item was either accepted or given
// It will also complete the offer if all offer items have been given and received, moving time credits around
func (t TradingService) ConfirmItemReceivedOrGiven(ctx context.Context, confirmedItemKey model.OfferItemKey) error {

	ctx, l := GetCtx(ctx, "ChatService", "ConfirmItemReceived")

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		return err
	}

	l.Debug("retrieving item")

	// retrieving item
	offerItem, err := t.tradingStore.GetItem(nil, confirmedItemKey)
	if err != nil {
		l.Error("could not get offer item", zap.Error(err))
		return err
	}

	if offerItem.GetFromUserKey() != userSession.GetUserKey() && offerItem.GetToUserKey() != userSession.GetUserKey() {
		return trading.ErrUserNotPartOfOfferItem
	}

	l.Debug("getting parent offer")

	offer, err := t.tradingStore.GetOffer(offerItem.GetOfferKey())
	if err != nil {
		l.Error("could not get parent offer", zap.Error(err))
		return err
	}

	l.Debug("make sure offer is not already completed")

	// cannot confirm an item for an offer that's already completed
	if offer.Status == trading.CompletedOffer {
		l.Warn("cannot complete offer: already completed")
		return err
	}

	l.Debug("make sure item being confirmed is either given or received by user making the request")

	// making sure that the item being confirmed is either given or received by the user making the request

	l.Debug("marking the item as either given or received")

	// mark the item as either received or given
	if userSession.GetUserKey() == offerItem.GetToUserKey() {

		l.Debug("item is being marked as received")

		// skip altogether when item is already marked as received
		if offerItem.IsReceived() {
			l.Warn("could not complete item: item already received")
			return err
		}

		err := t.tradingStore.ConfirmItemReceived(ctx, confirmedItemKey)
		if err != nil {
			l.Error("could not confirm having received offer item", zap.Error(err))
			return err
		}

	} else {

		l.Debug("item is being marked as given")

		// skip altogether when item is already marked as given
		if offerItem.IsGiven() {
			l.Warn("aborting: item already given")
			return err
		}

		err := t.tradingStore.ConfirmItemGiven(ctx, confirmedItemKey)
		if err != nil {
			l.Error("could not confirm having given the item", zap.Error(err))
			return err
		}

	}

	l.Debug("getting all offer items for offer")

	offerItems, err := t.tradingStore.GetItems(offerItem.GetOfferKey())
	if err != nil {
		l.Error("could not get all offer items", zap.Error(err))
		return err
	}

	// retrieving the item's from and to users
	confirmedItemFromUserKey := offerItem.GetFromUserKey()
	confirmedItemToUserKey := offerItem.GetToUserKey()

	l.Debug("getting all users involved in offer")

	offerUsers, err := t.us.GetByKeys(nil, []model.UserKey{confirmedItemFromUserKey, confirmedItemToUserKey})
	if err != nil {
		l.Error("could not get users", zap.Error(err))
		return err
	}

	confirmingUser, err := offerUsers.GetUser(userSession.GetUserKey())
	if err != nil {
		l.Error("could not get confirming user from user list")
		return err
	}

	l.Debug("notifying concerned users that the offer item was either given or received")

	err = t.notifyItemGivenOrReceived(ctx, offerItem, confirmingUser, offerUsers)
	if err != nil {
		l.Error("could not notify users the item was given/received", zap.Error(err))
		return err
	}

	err = t.checkOfferCompleted(ctx, offerItem.GetOfferKey(), offerItems, confirmingUser, offerUsers)
	if err != nil {
		return err
	}

	return nil

}
