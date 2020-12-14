package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	trading2 "github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (t TradingService) DeclineOffer(ctx context.Context, offerKey model.OfferKey) error {

	ctx, l := service.GetCtx(ctx, "TradingService", "AcceptOffer")
	l = l.With(zap.Object("offer", offerKey))

	user, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		l.Error("could not get user session", zap.Error(err))
		return err
	}
	loggedInUserKey := user.GetUserKey()

	offer, err := t.tradingStore.GetOffer(offerKey)
	if err != nil {
		l.Error("could not get offer", zap.Error(err))
		return err
	}

	if offer.Status != trading2.PendingOffer {
		return fmt.Errorf("could not decline a offer that is not pending")
	}

	approvers, err := t.tradingStore.FindApproversForOffer(offerKey)
	if err != nil {
		l.Error("could not get approvers for offer", zap.Error(err))
		return err
	}

	if !approvers.IsUserAnApprover(loggedInUserKey) {
		return errs.ErrForbidden
	}

	err = t.tradingStore.UpdateOfferStatus(offerKey, trading2.DeclinedOffer)
	if err != nil {
		return err
	}

	return nil

}

func (t TradingService) AcceptOffer(ctx context.Context, offerKey model.OfferKey) error {

	ctx, l := service.GetCtx(ctx, "TradingService", "AcceptOffer")
	l = l.With(zap.Object("offer", offerKey))

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		l.Error("could not get user session", zap.Error(err))
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	offer, err := t.tradingStore.GetOffer(offerKey)
	if err != nil {
		l.Error("could not get offer", zap.Error(err))
		return err
	}

	offerItems, err := t.tradingStore.GetOfferItemsForOffer(offerKey)
	if err != nil {
		l.Error("could not get offer items", zap.Error(err))
		return err
	}

	if offerItems.AllPartiesAccepted() {
		err := fmt.Errorf("offer is already accepted")
		l.Error("", zap.Error(err))
		return err
	}

	approvers, err := t.tradingStore.FindApproversForOffer(offerKey)
	if err != nil {
		l.Error("could not get approvers for offer", zap.Error(err))
		return err
	}

	approvableOfferItemsOnGivingSide, canApproveOfferGivingAnything :=
		approvers.OfferItemsUsersCanGive[loggedInUserKey]
	approvableOfferItemsOnReceivingSide, canApproveReceivingAnything :=
		approvers.OfferItemsUsersCanReceive[loggedInUserKey]

	if !canApproveOfferGivingAnything && !canApproveReceivingAnything {
		return errs.ErrUnauthorized
	}

	var offerItemsPendingGiverApproval []model.OfferItemKey
	if approvableOfferItemsOnGivingSide != nil {
		for _, offerItemKey := range approvableOfferItemsOnGivingSide.Items {
			offerItem := offerItems.GetOfferItem(offerItemKey)
			if offerItem.IsAcceptedByGiver() {
				continue
			}
			offerItemsPendingGiverApproval = append(offerItemsPendingGiverApproval, offerItemKey)
		}
	}
	var offerItemsPendingReceiverApproval []model.OfferItemKey
	if approvableOfferItemsOnReceivingSide != nil {
		for _, offerItemKey := range approvableOfferItemsOnReceivingSide.Items {
			offerItem := offerItems.GetOfferItem(offerItemKey)
			if offerItem.IsAcceptedByReceiver() {
				continue
			}
			offerItemsPendingReceiverApproval = append(offerItemsPendingReceiverApproval, offerItemKey)
		}
	}

	if len(offerItemsPendingReceiverApproval) == 0 && len(offerItemsPendingGiverApproval) == 0 {
		return fmt.Errorf("Nothing left to approve by you")
	}

	err = t.tradingStore.MarkOfferItemsAsAccepted(
		ctx,
		loggedInUserKey,
		model.NewOfferItemKeys(offerItemsPendingGiverApproval),
		model.NewOfferItemKeys(offerItemsPendingReceiverApproval))

	if err != nil {
		return err
	}

	//
	// var blocks []chat.Block
	//
	// blocks = append(blocks, *chat.NewHeaderBlock(
	// 	chat.NewMarkdownObject(fmt.Sprintf(":+1: Good news! [_%s_](/users/%s) has accepted the offer :)", user.GetUsername(), loggedInUserKey.String())),
	// 	nil))
	//
	// for _, user := range users.Items {
	// 	var userDecision *trading.OfferDecision
	// 	for _, decision := range decisions.Items {
	// 		if decision.GetUserKey() == user.GetUserKey() {
	// 			userDecision = decision
	// 			break
	// 		}
	// 	}
	// 	if userDecision == nil {
	// 		err := fmt.Errorf("could not find user decision")
	// 		l.Error("", zap.Error(err))
	// 		return nil, err
	// 	}
	//
	// 	if userDecision.Decision == trading.AcceptedDecision {
	// 		blocks = append(blocks, *chat.NewSectionBlock(
	// 			chat.NewMarkdownObject(fmt.Sprintf("[_%s_](/users/%s) has accepted the offer :relaxed:", user.Username, user.ID)),
	// 			nil,
	// 			nil,
	// 			nil))
	// 	} else if userDecision.Decision == trading.PendingDecision {
	// 		blocks = append(blocks, *chat.NewSectionBlock(
	// 			chat.NewMarkdownObject(fmt.Sprintf("[_%s_](/users/%s) still no answer :expressionless:", user.Username, user.ID)),
	// 			nil,
	// 			nil,
	// 			nil))
	// 	} else if userDecision.Decision == trading.DeclinedDecision {
	// 		blocks = append(blocks, *chat.NewSectionBlock(
	// 			chat.NewMarkdownObject(fmt.Sprintf("[_%s_](/users/%s) declined the offer :slightly_frowning_face:", user.Username, user.ID)),
	// 			nil,
	// 			nil,
	// 			nil))
	// 	}
	// }
	//
	// sendMessage := chat.NewSendConversationMessage(
	// 	loggedInUserKey,
	// 	user.GetUsername(),
	// 	userKeys,
	// 	"",
	// 	blocks,
	// 	[]chat.Attachment{},
	// 	nil,
	// )
	// _, err = t.chatService.SendConversationMessage(ctx, sendMessage)
	// if err != nil {
	// 	l.Error("could not send message", zap.Error(err))
	// 	return nil, err
	// }

	offerItems, err = t.tradingStore.GetOfferItemsForOffer(offerKey)
	if err != nil {
		l.Error("could not get offer items", zap.Error(err))
		return err
	}

	if offerItems.AllPartiesAccepted() {
		err := t.tradingStore.UpdateOfferStatus(offerKey, trading2.AcceptedOffer)
		if err != nil {
			l.Error("could not update offer status", zap.Error(err))
			return err
		}
	}

	usersInOffer, err := t.us.GetByKeys(ctx, approvers.AllUserKeys().Items)
	if err != nil {
		return err
	}

	err = t.checkOfferCompleted(ctx, offer.GroupKey, offerKey, offerItems, loggedInUser, usersInOffer)
	if err != nil {
		return err
	}

	// if offerItems.AllPartiesAccepted() {
	//
	// 	var blocks []chat.Block
	// 	blocks = append(blocks, *chat.NewHeaderBlock(
	// 		chat.NewMarkdownObject(fmt.Sprintf(":champagne: Alright! Everyone accepted the offer!")),
	// 		nil))
	//
	// 	linkBlock := chat.NewSectionBlock(
	// 		chat.NewMarkdownObject(
	// 			fmt.Sprintf("It's now time to do your thing! Once you've kept up with your side of the bargain, "+
	// 				"just mark it as 'completed' by going into your [transactions](/transactions). Also, when other parties "+
	// 				"give you what was agreed upon, you also have to confirm it, so we can exchange hours from your timebanks.")),
	// 		nil,
	// 		nil,
	// 		nil)
	// 	blocks = append(blocks, *linkBlock)
	//
	// 	sendMessage := chat.NewSendConversationMessage(
	// 		loggedInUserKey,
	// 		user.GetUsername(),
	// 		userKeys,
	// 		"",
	// 		blocks,
	// 		[]chat.Attachment{},
	// 		nil,
	// 	)
	// 	_, err = t.chatService.SendConversationMessage(ctx, sendMessage)
	// 	if err != nil {
	// 		l.Error("could not send conversation message", zap.Error(err))
	// 		return nil, err
	// 	}
	//
	// }

	//
	// var resources = resource.NewEmptyResources()
	// if len(offerItems.Items) > 0 {
	// 	getResourcesByKeysResponse, err := t.rs.GetByKeys(ctx, offerItems.GetResourceKeys())
	// 	if err != nil {
	// 		l.Error("could not get resources by keys", zap.Error(err))
	// 		return nil, err
	// 	}
	// 	resources = getResourcesByKeysResponse
	// }
	//
	// if currentUserLastOneToDecide {
	//
	// 	for _, user := range users.Items {
	// 		userKey := user.GetUserKey()
	//
	// 		var blocks []chat.Block
	//
	// 		userItems := offerItems.GetOfferItemsReceivedByUser(user.GetUserKey())
	// 		for _, userItem := range userItems.Items {
	// 			userItemKey := userItem.GetKey()
	// 			userItemIdValue := userItemKey.ID.String()
	// 			if userItem.IsReceivedBy(user.GetUserKey()) {
	//
	// 				if !userItem.IsResourceExchangeItem() {
	// 					continue
	// 				}
	//
	// 				actionId := "confirm_item_received"
	// 				res, _ := resources.GetResource(userItem.GetResourceKey())
	// 				fromUser, _ := users.GetUser(userItem.GetFromUserKey())
	// 				block := chat.NewSectionBlock(
	// 					chat.NewMarkdownObject(fmt.Sprintf("**%s**", res.Summary)),
	// 					[]chat.BlockElement{},
	// 					chat.NewButtonElement(
	// 						chat.NewMarkdownObject(fmt.Sprintf("I received it from **%s**", fromUser.Username)),
	// 						nil,
	// 						&actionId,
	// 						nil,
	// 						&userItemIdValue,
	// 						nil),
	// 					nil)
	// 				blocks = append(blocks, *block)
	//
	// 			} else if userItem.IsGivenBy(user.GetUserKey()) {
	// 				if !userItem.IsResourceExchangeItem() {
	// 					continue
	// 				}
	//
	// 				actionId := "confirm_item_given"
	// 				res, _ := resources.GetResource(userItem.GetResourceKey())
	// 				toUser, _ := users.GetUser(userItem.GetToUserKey())
	// 				block := chat.NewSectionBlock(
	// 					chat.NewMarkdownObject(fmt.Sprintf("**%s**", res.Summary)),
	// 					[]chat.BlockElement{},
	// 					chat.NewButtonElement(
	// 						chat.NewMarkdownObject(fmt.Sprintf("I've given it to **%s**", toUser.Username)),
	// 						nil,
	// 						&actionId,
	// 						nil,
	// 						&userItemIdValue,
	// 						nil),
	// 					nil)
	// 				blocks = append(blocks, *block)
	//
	// 			}
	//
	// 			blocks = append(blocks, *chat.NewDividerBlock())
	// 		}
	//
	// 		if blocks == nil {
	// 			blocks = []chat.Block{}
	// 		}
	//
	// 		if blocks[len(blocks)-1].Type == chat.Divider {
	// 			blocks = blocks[:len(blocks)-1]
	// 		}
	//
	// 		sendMessage := chat.NewSendConversationMessage(
	// 			user.GetUserKey(),
	// 			user.GetUsername(),
	// 			userKeys,
	// 			"",
	// 			blocks,
	// 			[]chat.Attachment{},
	// 			&userKey,
	// 		)
	// 		_, err = t.chatService.SendConversationMessage(nil, sendMessage)
	// 		if err != nil {
	// 			l.Error("could not send conversation message", zap.Error(err))
	// 			return nil, err
	// 		}
	// 	}
	// }

	return nil
}
