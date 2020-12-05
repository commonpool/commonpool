package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
	"go.uber.org/zap"
)

func (t TradingService) AcceptOffer(ctx context.Context, request *trading.AcceptOffer) (*trading.AcceptOfferResponse, error) {

	ctx, l := GetCtx(ctx, "TradingService", "AcceptOffer")

	user, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session", zap.Error(err))
		return nil, err
	}
	loggedInUserKey := user.GetUserKey()

	//  Retrieving offer
	offer, err := t.tradingStore.GetOffer(request.OfferKey)
	if err != nil {
		l.Error("could not get offer", zap.Error(err), zap.String("offerId", request.OfferKey.ID.String()))
		return nil, err
	}

	//  Ensure offer is still pending approval
	if offer.Status != trading.PendingOffer {
		err := fmt.Errorf("offer is not pending approval")
		l.Error("", zap.Error(err))
		return nil, err
	}

	//  Retrieve Offer decisions
	decisions, err := t.tradingStore.GetDecisions(request.OfferKey)
	if err != nil {
		l.Error("could not get offer decisions", zap.Error(err))
		return nil, err
	}

	var didAllOtherParticipantsAlreadyAccept = true
	var currentUserDecision *trading.OfferDecision

	//  Retrieving current user decision, and check if everyone else accepted the offer already

	// todo: wrap in collection struct

	for _, decision := range decisions.Items {
		decisionKey := decision.GetKey()
		decisionUserKey := decisionKey.GetUserKey()
		if decisionUserKey != loggedInUserKey {
			if decision.Decision != trading.AcceptedDecision {
				didAllOtherParticipantsAlreadyAccept = false
			}
		} else {
			currentUserDecision = decision
		}
	}

	// in the case that a user is approving an offer he's not part of
	if currentUserDecision == nil {
		err := fmt.Errorf("could not find current user decision")
		l.Error("", zap.Error(err))
		return nil, err
	}

	//  Persisting the decision
	err = t.tradingStore.SaveDecision(request.OfferKey, loggedInUserKey, trading.AcceptedDecision)
	if err != nil {
		return nil, err
	}

	decisions, err = t.tradingStore.GetDecisions(request.OfferKey)
	if err != nil {
		return nil, err
	}

	// todo: wrap in collection object
	var userKeyLst []model.UserKey
	for _, decision := range decisions.Items {
		userKeyLst = append(userKeyLst, decision.GetUserKey())
	}
	userKeys := model.NewUserKeys(userKeyLst)
	users, err := t.us.GetByKeys(ctx, userKeyLst)
	if err != nil {
		return nil, err
	}

	//  Complete offer if everyone accepted already
	var currentUserLastOneToDecide = didAllOtherParticipantsAlreadyAccept
	if currentUserLastOneToDecide {

		l.Debug("user is last one to decide. mark offer as accepted")
		err = t.tradingStore.SaveOfferStatus(request.OfferKey, trading.AcceptedOffer)
		if err != nil {
			l.Error("could not save offer status", zap.Error(err))
			return nil, err
		}

	}

	var blocks []chat.Block

	blocks = append(blocks, *chat.NewHeaderBlock(
		chat.NewMarkdownObject(fmt.Sprintf(":+1: Good news! [_%s_](/users/%s) has accepted the offer :)", user.GetUsername(), loggedInUserKey.String())),
		nil))

	for _, user := range users.Items {
		var userDecision *trading.OfferDecision
		for _, decision := range decisions.Items {
			if decision.GetUserKey() == user.GetUserKey() {
				userDecision = decision
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
	_, err = t.chatService.SendConversationMessage(ctx, sendMessage)
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
		_, err = t.chatService.SendConversationMessage(ctx, sendMessage)
		if err != nil {
			l.Error("could not send conversation message", zap.Error(err))
			return nil, err
		}

	}

	offerItems, err := t.tradingStore.GetItems(request.OfferKey)
	if err != nil {
		l.Error("could not get offer items", zap.Error(err))
		return nil, err
	}

	var resources = resource.NewEmptyResources()
	if len(offerItems.Items) > 0 {
		getResourcesByKeysResponse, err := t.rs.GetByKeys(ctx, offerItems.GetResourceKeys())
		if err != nil {
			l.Error("could not get resources by keys", zap.Error(err))
			return nil, err
		}
		resources = getResourcesByKeysResponse
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
					res, _ := resources.GetResource(userItem.GetResourceKey())
					fromUser, _ := users.GetUser(userItem.GetFromUserKey())
					block := chat.NewSectionBlock(
						chat.NewMarkdownObject(fmt.Sprintf("**%s**", res.Summary)),
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
					res, _ := resources.GetResource(userItem.GetResourceKey())
					toUser, _ := users.GetUser(userItem.GetToUserKey())
					block := chat.NewSectionBlock(
						chat.NewMarkdownObject(fmt.Sprintf("**%s**", res.Summary)),
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
			_, err = t.chatService.SendConversationMessage(nil, sendMessage)
			if err != nil {
				l.Error("could not send conversation message", zap.Error(err))
				return nil, err
			}
		}
	}

	return &trading.AcceptOfferResponse{}, nil
}
