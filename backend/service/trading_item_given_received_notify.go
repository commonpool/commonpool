package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
	"go.uber.org/zap"
)

func (t TradingService) notifyItemGivenOrReceived(ctx context.Context, offerItemBeingConfirmed *trading.OfferItem, confirmingUser *auth.User, concernedOfferUsers *auth.Users) error {

	l := logging.WithContext(ctx)

	l.Debug("getting offer item resource")

	// confirming items is only for "Resource" offer items, so it's safe to assume that
	// the item.resourceKey is not going to be nil
	getResource := t.rs.GetByKey(ctx, resource.NewGetResourceByKeyQuery(offerItemBeingConfirmed.GetResourceKey()))
	if getResource.Error != nil {
		l.Error("could not get offer item resource", zap.Error(getResource.Error))
		return getResource.Error
	}
	resourceSummary := getResource.Resource.Summary

	offerItemFromUserKey := offerItemBeingConfirmed.GetFromUserKey()
	offerItemToUserKey := offerItemBeingConfirmed.GetToUserKey()

	// building sentence component for sending message
	var verb string
	var article string
	var otherUserName string
	if offerItemBeingConfirmed.IsGivenBy(confirmingUser.GetUserKey()) {
		verb = "given"
		article = "to"

		toUser, err := concernedOfferUsers.GetUser(offerItemToUserKey)
		if err != nil {
			l.Error("could not get 'to' user", zap.Error(err))
			return err
		}

		otherUserName = toUser.Username

	} else if offerItemBeingConfirmed.IsReceivedBy(confirmingUser.GetUserKey()) {
		verb = "received"
		article = "from"

		fromUser, err := concernedOfferUsers.GetUser(offerItemFromUserKey)
		if err != nil {
			l.Error("could not get 'from' user", zap.Error(err))
			return err
		}

		otherUserName = fromUser.Username

	}

	_, err := t.chatService.SendConversationMessage(ctx, chat.NewSendConversationMessage(
		confirmingUser.GetUserKey(),
		confirmingUser.Username,
		concernedOfferUsers.GetUserKeys(),
		"",
		[]chat.Block{
			*chat.NewHeaderBlock(chat.NewMarkdownObject(
				fmt.Sprintf(":heavy_check_mark: **%s** has confirmed having %s **%s** %s **%s**",
					confirmingUser.Username,
					verb,
					resourceSummary,
					article,
					otherUserName,
				),
			),
				nil),
		},
		[]chat.Attachment{},
		nil,
	))

	if err != nil {
		l.Error("could not send message to users")
		return err
	}

	return nil
}
