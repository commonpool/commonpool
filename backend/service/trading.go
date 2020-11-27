package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/model"
	res "github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
	"go.uber.org/zap"
	"time"
)

type TradingService struct {
	tradingStore trading.Store
	rs           res.Store
	us           auth.Store
	cs           chat.Service
}

func (h TradingService) AcceptOffer(ctx context.Context, request *trading.AcceptOffer) (*trading.AcceptOfferResponse, error) {

	ctx, l := GetCtx(ctx, "TradingService", "AcceptOffer")

	l.Debug("getting user session")

	user, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session", zap.Error(err))
		return nil, err
	}
	loggedInUserKey := user.GetUserKey()

	l.Debug("retrieving offer")

	//  Retrieving offer
	offer, err := h.tradingStore.GetOffer(request.OfferKey)
	if err != nil {
		l.Error("could not get offer", zap.Error(err), zap.String("offerId", request.OfferKey.ID.String()))
		return nil, err
	}

	l.Debug("ensure offer is still pending")

	//  Ensure offer is still pending approval
	if offer.Status != trading.PendingOffer {
		err := fmt.Errorf("offer is not pending approval")
		l.Error("", zap.Error(err))
		return nil, err
	}

	l.Debug("getting offer decisions")

	//  Retrieve Offer decisions
	decisions, err := h.tradingStore.GetDecisions(request.OfferKey)
	if err != nil {
		l.Error("could not get offer decisions", zap.Error(err))
		return nil, err
	}

	var didAllOtherParticipantsAlreadyAccept = true
	var currentUserDecision *trading.OfferDecision

	//  Retrieving current user decision, and check if everyone else accepted the offer already

	l.Debug("retrieving current user decision")

	// todo: wrap in collection struct

	for _, decision := range decisions {
		decisionKey := decision.GetKey()
		decisionUserKey := decisionKey.GetUserKey()
		if decisionUserKey != loggedInUserKey {
			if decision.Decision != trading.AcceptedDecision {
				didAllOtherParticipantsAlreadyAccept = false
			}
		} else {
			currentUserDecision = &decision
		}
	}

	// in the case that a user is approving an offer he's not part of
	if currentUserDecision == nil {
		err := fmt.Errorf("could not find current user decision")
		l.Error("", zap.Error(err))
		return nil, err
	}

	l.Debug("persisting user decision")

	//  Persisting the decision
	err = h.tradingStore.SaveDecision(request.OfferKey, loggedInUserKey, trading.AcceptedDecision)
	if err != nil {
		return nil, err
	}

	l.Debug("fetching user decisions anew")

	decisions, err = h.tradingStore.GetDecisions(request.OfferKey)
	if err != nil {
		return nil, err
	}

	// todo: wrap in collection object
	var userKeyLst []model.UserKey
	for _, decision := range decisions {
		userKeyLst = append(userKeyLst, decision.GetUserKey())
	}
	userKeys := model.NewUserKeys(userKeyLst)
	users, err := h.us.GetByKeys(ctx, userKeyLst)
	if err != nil {
		return nil, err
	}

	//  Complete offer if everyone accepted already
	var currentUserLastOneToDecide = didAllOtherParticipantsAlreadyAccept
	if currentUserLastOneToDecide {

		l.Debug("user is last one to decide. mark offer as accepted")
		err = h.tradingStore.SaveOfferStatus(request.OfferKey, trading.AcceptedOffer)
		if err != nil {
			l.Error("could not save offer status", zap.Error(err))
			return nil, err
		}

	}

	l.Debug("building chat messages")

	var blocks []chat.Block

	blocks = append(blocks, *chat.NewHeaderBlock(
		chat.NewMarkdownObject(fmt.Sprintf(":+1: Good news! [_%s_](/users/%s) has accepted the offer :)", user.GetUsername(), loggedInUserKey.String())),
		nil))

	for _, user := range users.Items {
		var userDecision *trading.OfferDecision
		for _, decision := range decisions {
			if decision.GetUserKey() == user.GetUserKey() {
				userDecision = &decision
				break
			}
		}
		if userDecision == nil {
			err := fmt.Errorf("could not find user decision")
			l.Error("", zap.Error(err))
			return nil, err
		}

		if userDecision.Decision == trading.AcceptedDecision {
			blocks = append(blocks, *chat.NewSectionBlock(
				chat.NewMarkdownObject(fmt.Sprintf("[_%s_](/users/%s) has accepted the offer :relaxed:", user.Username, user.ID)),
				nil,
				nil,
				nil))
		} else if userDecision.Decision == trading.PendingDecision {
			blocks = append(blocks, *chat.NewSectionBlock(
				chat.NewMarkdownObject(fmt.Sprintf("[_%s_](/users/%s) still no answer :expressionless:", user.Username, user.ID)),
				nil,
				nil,
				nil))
		} else if userDecision.Decision == trading.DeclinedDecision {
			blocks = append(blocks, *chat.NewSectionBlock(
				chat.NewMarkdownObject(fmt.Sprintf("[_%s_](/users/%s) declined the offer :slightly_frowning_face:", user.Username, user.ID)),
				nil,
				nil,
				nil))
		}
	}

	sendMessage := chat.NewSendConversationMessage(
		loggedInUserKey,
		user.GetUsername(),
		userKeys,
		"",
		blocks,
		[]chat.Attachment{},
		nil,
	)
	_, err = h.cs.SendConversationMessage(ctx, sendMessage)
	if err != nil {
		l.Error("could not send message", zap.Error(err))
		return nil, err
	}

	if currentUserLastOneToDecide {

		var blocks []chat.Block
		blocks = append(blocks, *chat.NewHeaderBlock(
			chat.NewMarkdownObject(fmt.Sprintf(":champagne: Alright! Everyone accepted the offer!")),
			nil))

		linkBlock := chat.NewSectionBlock(
			chat.NewMarkdownObject(
				fmt.Sprintf("It's now time to do your thing! Once you've kept up with your side of the bargain, "+
					"just mark it as 'completed' by going into your [transactions](/transactions). Also, when other parties "+
					"give you what was agreed upon, you also have to confirm it, so we can exchange hours from your timebanks.")),
			nil,
			nil,
			nil)
		blocks = append(blocks, *linkBlock)

		sendMessage := chat.NewSendConversationMessage(
			loggedInUserKey,
			user.GetUsername(),
			userKeys,
			"",
			blocks,
			[]chat.Attachment{},
			nil,
		)
		_, err = h.cs.SendConversationMessage(ctx, sendMessage)
		if err != nil {
			l.Error("could not send conversation message", zap.Error(err))
			return nil, err
		}

	}

	l.Debug("getting offer items")

	offerItems, err := h.tradingStore.GetItems(request.OfferKey)
	if err != nil {
		l.Error("could not get offer items", zap.Error(err))
		return nil, err
	}

	l.Debug("getting resources in offer")

	var resources *res.Resources = res.NewResources([]res.Resource{})
	if len(offerItems.Items) > 0 {
		getByKeys := res.NewGetResourceByKeysQuery(offerItems.GetResourceKeys())
		getResourcesByKeysResponse, err := h.rs.GetByKeys(getByKeys)
		if err != nil {
			l.Error("could not get resources by keys", zap.Error(err))
			return nil, err
		}
		resources = getResourcesByKeysResponse.Items
	}

	if currentUserLastOneToDecide {

		for _, user := range users.Items {
			userKey := user.GetUserKey()

			var blocks []chat.Block

			userItems := offerItems.GetOfferItemsForUser(user.GetUserKey())
			for _, userItem := range userItems.Items {
				userItemKey := userItem.GetKey()
				userItemIdValue := userItemKey.ID.String()
				if userItem.IsReceivedBy(user.GetUserKey()) {

					if !userItem.IsResourceExchangeItem() {
						continue
					}

					actionId := "confirm_item_received"
					resource, _ := resources.GetResource(userItem.GetResourceKey())
					fromUser, _ := users.GetUser(userItem.GetFromUserKey())
					block := chat.NewSectionBlock(
						chat.NewMarkdownObject(fmt.Sprintf("**%s**", resource.Summary)),
						[]chat.BlockElement{},
						chat.NewButtonElement(
							chat.NewMarkdownObject(fmt.Sprintf("I received it from **%s**", fromUser.Username)),
							nil,
							&actionId,
							nil,
							&userItemIdValue,
							nil),
						nil)
					blocks = append(blocks, *block)

				} else if userItem.IsGivenBy(user.GetUserKey()) {
					if !userItem.IsResourceExchangeItem() {
						continue
					}

					actionId := "confirm_item_given"
					resource, _ := resources.GetResource(userItem.GetResourceKey())
					toUser, _ := users.GetUser(userItem.GetToUserKey())
					block := chat.NewSectionBlock(
						chat.NewMarkdownObject(fmt.Sprintf("**%s**", resource.Summary)),
						[]chat.BlockElement{},
						chat.NewButtonElement(
							chat.NewMarkdownObject(fmt.Sprintf("I've given it to **%s**", toUser.Username)),
							nil,
							&actionId,
							nil,
							&userItemIdValue,
							nil),
						nil)
					blocks = append(blocks, *block)

				}

				blocks = append(blocks, *chat.NewDividerBlock())
			}

			if blocks == nil {
				blocks = []chat.Block{}
			}

			if blocks[len(blocks)-1].Type == chat.Divider {
				blocks = blocks[:len(blocks)-1]
			}

			sendMessage := chat.NewSendConversationMessage(
				user.GetUserKey(),
				user.GetUsername(),
				userKeys,
				"",
				blocks,
				[]chat.Attachment{},
				&userKey,
			)
			_, err = h.cs.SendConversationMessage(nil, sendMessage)
			if err != nil {
				l.Error("could not send conversation message", zap.Error(err))
				return nil, err
			}
		}
	}

	return &trading.AcceptOfferResponse{}, nil
}

var _ trading.Service = &TradingService{}

func NewTradingService(tradingStore trading.Store, resourceStore res.Store, authStore auth.Store, chatService chat.Service) *TradingService {
	return &TradingService{
		tradingStore: tradingStore,
		rs:           resourceStore,
		us:           authStore,
		cs:           chatService,
	}
}

// ConfirmItemReceived marks an offerItem as being received by a user
func (t TradingService) ConfirmItemReceived(ctx context.Context, request *trading.ConfirmItemReceived) (*trading.ConfirmItemReceivedResponse, error) {
	ctx, l := GetCtx(ctx, "ChatService", "ConfirmItemReceived")

	ctx = logging.NewContext(
		ctx,
		zap.String("userId", request.ReceivedByUser.String()),
		zap.String("offerItemId", request.OfferItemKey.ID.String()),
	)

	l.Debug("confirming item was sent")

	err := t.ConfirmItemReceivedOrGiven(ctx, request.OfferItemKey, request.ReceivedByUser, trading.OfferItemReceiving)
	if err != nil {
		return nil, err
	}
	return &trading.ConfirmItemReceivedResponse{}, nil
}

// ConfirmItemReceived marks an offerItem as being given by a user
func (t TradingService) ConfirmItemGiven(ctx context.Context, request *trading.ConfirmItemGiven) (*trading.ConfirmItemGivenResponse, error) {

	ctx, l := GetCtx(ctx, "ChatService", "ConfirmItemReceived")

	ctx = logging.NewContext(
		ctx,
		zap.String("userId", request.GivenByUser.String()),
		zap.String("offerItemId", request.OfferItemKey.ID.String()),
	)

	l.Debug("confirming item was given")

	err := t.ConfirmItemReceivedOrGiven(ctx, request.OfferItemKey, request.GivenByUser, trading.OfferItemGiving)
	if err != nil {
		return nil, err
	}
	return &trading.ConfirmItemGivenResponse{}, nil
}

// ConfirmItemReceivedOrGiven will send a notification to concerned users that an offer item was either accepted or given
// It will also complete the offer if all offer items have been given and received, moving time credits around
func (t TradingService) ConfirmItemReceivedOrGiven(ctx context.Context, confirmedItemKey model.OfferItemKey, confirmingUserKey model.UserKey, expectedOfferItemSide trading.OfferItemBond) error {

	ctx, l := GetCtx(ctx, "ChatService", "ConfirmItemReceived")

	l.Debug("retrieving item")

	// retrieving item
	offerItem, err := t.tradingStore.GetItem(confirmedItemKey)
	if err != nil {
		l.Error("could not get offer item", zap.Error(err))
		return err
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

	userItemDirection := offerItem.GetUserBondDirection(confirmingUserKey)
	if userItemDirection != expectedOfferItemSide {
		err := fmt.Errorf("cannot accept this offer. user is not confirming the right side of the item")
		l.Error("", zap.Error(err))
		return err
	}

	if userItemDirection == trading.OfferItemNeither {
		err := fmt.Errorf("cannot accept this offer. neither at receiving or giving end")
		l.Error("cannot accept offer", zap.Error(err))
		return err
	}

	l.Debug("marking the item as either given or received")

	// mark the item as either received or given
	if userItemDirection == trading.OfferItemReceiving {

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

	} else if userItemDirection == trading.OfferItemGiving {

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

	} else {
		msg := "unexpected item bond direction"
		l.Error(msg)
		return fmt.Errorf(msg)
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

	confirmingUser, err := offerUsers.GetUser(confirmingUserKey)
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

func (t TradingService) notifyItemGivenOrReceived(ctx context.Context, offerItemBeingConfirmed *trading.OfferItem, confirmingUser auth.User, concernedOfferUsers auth.Users) error {

	l := logging.WithContext(ctx)

	l.Debug("getting offer item resource")

	// confirming items is only for "Resource" offer items, so it's safe to assume that
	// the item.resourceKey is not going to be nil
	getResource := t.rs.GetByKey(ctx, res.NewGetResourceByKeyQuery(offerItemBeingConfirmed.GetResourceKey()))
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

	_, err := t.cs.SendConversationMessage(ctx, chat.NewSendConversationMessage(
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

func (t TradingService) checkOfferCompleted(ctx context.Context, offerKey model.OfferKey, offerItems *trading.OfferItems, userConfirmingItem auth.User, usersInOffer auth.Users) error {

	ctx, l := GetCtx(ctx, "TradingService", "checkOfferCompleted")

	if offerItems.AllResourceItemsReceivedAndGiven() {

		l.Debug("all items have been given and received. Marking offer as completed")

		err := t.tradingStore.SaveOfferStatus(offerKey, trading.CompletedOffer)
		if err != nil {
			l.Error("could not mark offer as completed", zap.Error(err))
			return err
		}

		l.Debug("building message to send these users")

		blocks, mainText, err := t.buildOfferCompletedMessage(ctx, offerItems, usersInOffer)
		if err != nil {
			l.Debug("could not build offer completion message", zap.Error(err))
			return err
		}

		_, err = t.cs.SendConversationMessage(ctx, chat.NewSendConversationMessage(
			userConfirmingItem.GetUserKey(),
			userConfirmingItem.Username,
			usersInOffer.GetUserKeys(),
			fmt.Sprintf(mainText),
			blocks,
			[]chat.Attachment{},
			nil,
		))
	}
	return nil
}

func (t TradingService) buildOfferCompletedMessage(ctx context.Context, items *trading.OfferItems, users auth.Users) ([]chat.Block, string, error) {

	ctx, l := GetCtx(ctx, "TradingService", "buildOfferCompletedMessage")

	var blocks []chat.Block

	mainText := ":champagne: Alright! everybody confirmed having received and given their stuff."
	blocks = append(blocks, *chat.NewHeaderBlock(
		chat.NewMarkdownObject(mainText),
		nil,
	))

	for _, offerItem := range items.Items {

		fromUser, err := users.GetUser(offerItem.GetFromUserKey())
		if err != nil {
			l.Error("could not get 'fromUser'", zap.Error(err))
			return nil, "", err
		}

		toUser, err := users.GetUser(offerItem.GetToUserKey())
		if err != nil {
			l.Error("could not get 'toUser'", zap.Error(err))
			return nil, "", err
		}

		if offerItem.IsTimeExchangeItem() {
			blocks = append(blocks, *chat.NewSectionBlock(
				chat.NewMarkdownObject(fmt.Sprintf("**%s** sent **%s** `%s` timebank credits",
					fromUser.Username,
					toUser.Username,
					time.Duration(int64(time.Second)**offerItem.OfferedTimeInSeconds).Truncate(time.Minute*1).String(),
				)),
				nil,
				nil,
				nil,
			))
		}
	}

	return blocks, mainText, nil

}
