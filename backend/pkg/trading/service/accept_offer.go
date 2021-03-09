package service

//
// import (
// 	"context"
// 	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
// 	"github.com/commonpool/backend/pkg/keys"
// 	"github.com/commonpool/backend/pkg/trading/domain"
// )
//
// func (t TradingService) AcceptOffer(ctx context.Context, offerKey keys.OfferKey) error {
//
// 	loggedInUser, err := oidc.GetLoggedInUser(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	loggedInUserKey := loggedInUser.GetUserKey()
//
// 	domainOffer, err := t.offerRepo.Load(ctx, offerKey)
// 	if err != nil {
// 		return err
// 	}
//
// 	approvers, err := t.tradingStore.FindApproversForOffer(offerKey)
// 	if err != nil {
// 		return err
// 	}
//
// 	if err = domainOffer.ApproveAll(loggedInUserKey, func(userKey keys.UserKey, offerItem domain.OfferItem, direction domain.ApprovalDirection) bool {
// 		if direction == domain.Inbound {
// 			return approvers.GetInboundApprovers(offerItem.GetKey()).Contains(userKey)
// 		} else if direction == domain.Outbound {
// 			return approvers.GetOutboundApprovers(offerItem.GetKey()).Contains(userKey)
// 		} else {
// 			return false
// 		}
// 	}); err != nil {
// 		return err
// 	}
//
// 	if err := t.offerRepo.Save(ctx, domainOffer); err != nil {
// 		return err
// 	}
//
// 	return nil
//
// 	// }
// 	//
// 	// offer, err := t.tradingStore.GetOffer(offerKey)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//
// 	// offerItems, err := t.tradingStore.GetOfferItemsForOffer(offerKey)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//
// 	// if offerItems.AllApproved() {
// 	// 	err := fmt.Errorf("offer is already accepted")
// 	// 	return err
// 	// }
// 	//
// 	// if !approvers.HasAnyOfferItemsToApprove(loggedInUserKey) {
// 	// 	return exceptions.ErrUnauthorized
// 	// }
// 	//
// 	// outboundApprovableItems :=
// 	// 	approvers.GetOutboundOfferItems(loggedInUserKey)
// 	//
// 	// inboundApprovableItems :=
// 	// 	approvers.GetInboundOfferItems(loggedInUserKey)
// 	//
// 	// var offerItemsPendingOutboundApproval []keys.OfferItemKey
// 	// if outboundApprovableItems != nil {
// 	// 	for _, offerItemKey := range outboundApprovableItems.Items {
// 	// 		offerItem := offerItems.GetOfferItem(offerItemKey)
// 	// 		if offerItem.IsOutboundApproved() {
// 	// 			continue
// 	// 		}
// 	// 		offerItemsPendingOutboundApproval = append(offerItemsPendingOutboundApproval, offerItemKey)
// 	// 	}
// 	// }
// 	// var offerItemsPendingInboundApproval []keys.OfferItemKey
// 	// if inboundApprovableItems != nil {
// 	// 	for _, offerItemKey := range inboundApprovableItems.Items {
// 	// 		offerItem := offerItems.GetOfferItem(offerItemKey)
// 	// 		if offerItem.IsInboundApproved() {
// 	// 			continue
// 	// 		}
// 	// 		offerItemsPendingInboundApproval = append(offerItemsPendingInboundApproval, offerItemKey)
// 	// 	}
// 	// }
// 	//
// 	// if len(offerItemsPendingInboundApproval) == 0 && len(offerItemsPendingOutboundApproval) == 0 {
// 	// 	return fmt.Errorf("nothing left to approve by you")
// 	// }
// 	//
// 	// err = t.tradingStore.MarkOfferItemsAsAccepted(
// 	// 	ctx,
// 	// 	loggedInUserKey,
// 	// 	keys.NewOfferItemKeys(offerItemsPendingOutboundApproval),
// 	// 	keys.NewOfferItemKeys(offerItemsPendingInboundApproval))
// 	//
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//
// 	// //
// 	// // var blocks []chat.Block
// 	// //
// 	// // blocks = append(blocks, *chat.NewHeaderBlock(
// 	// // 	chat.NewMarkdownObject(fmt.Sprintf(":+1: Good news! [_%s_](/users/%s) has accepted the offer :)", user.GetUsername(), loggedInUserKey.String())),
// 	// // 	nil))
// 	// //
// 	// // for _, user := range users.Items {
// 	// // 	var userDecision *trading.OfferDecision
// 	// // 	for _, decision := range decisions.Items {
// 	// // 		if decision.GetUserKey() == user.GetUserKey() {
// 	// // 			userDecision = decision
// 	// // 			break
// 	// // 		}
// 	// // 	}
// 	// // 	if userDecision == nil {
// 	// // 		err := fmt.Errorf("could not find user decision")
// 	// // 		l.Error("", zap.Error(err))
// 	// // 		return nil, err
// 	// // 	}
// 	// //
// 	// // 	if userDecision.Decision == trading.AcceptedDecision {
// 	// // 		blocks = append(blocks, *chat.NewSectionBlock(
// 	// // 			chat.NewMarkdownObject(fmt.Sprintf("[_%s_](/users/%s) has accepted the offer :relaxed:", user.Username, user.ID)),
// 	// // 			nil,
// 	// // 			nil,
// 	// // 			nil))
// 	// // 	} else if userDecision.Decision == trading.PendingDecision {
// 	// // 		blocks = append(blocks, *chat.NewSectionBlock(
// 	// // 			chat.NewMarkdownObject(fmt.Sprintf("[_%s_](/users/%s) still no answer :expressionless:", user.Username, user.ID)),
// 	// // 			nil,
// 	// // 			nil,
// 	// // 			nil))
// 	// // 	} else if userDecision.Decision == trading.DeclinedDecision {
// 	// // 		blocks = append(blocks, *chat.NewSectionBlock(
// 	// // 			chat.NewMarkdownObject(fmt.Sprintf("[_%s_](/users/%s) declined the offer :slightly_frowning_face:", user.Username, user.ID)),
// 	// // 			nil,
// 	// // 			nil,
// 	// // 			nil))
// 	// // 	}
// 	// // }
// 	// //
// 	// // sendMessage := chat.NewSendConversationMessage(
// 	// // 	loggedInUserKey,
// 	// // 	user.GetUsername(),
// 	// // 	userKeys,
// 	// // 	"",
// 	// // 	blocks,
// 	// // 	[]chat.Attachment{},
// 	// // 	nil,
// 	// // )
// 	// // _, err = t.chatService.SendConversationMessage(ctx, sendMessage)
// 	// // if err != nil {
// 	// // 	l.Error("could not send message", zap.Error(err))
// 	// // 	return nil, err
// 	// // }
// 	//
// 	// offerItems, err = t.tradingStore.GetOfferItemsForOffer(offerKey)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//
// 	// if offerItems.AllApproved() {
// 	// 	err := t.tradingStore.UpdateOfferStatus(offerKey, trading.AcceptedOffer)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// }
// 	//
// 	// usersInOffer, err := t.userStore.GetByKeys(ctx, approvers.AllUserKeys())
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//
// 	// err = t.checkOfferCompleted(ctx, offer.GroupKey, offerKey, offerItems, loggedInUser, usersInOffer)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//
// 	// // if offerItems.AllApproved() {
// 	// //
// 	// // 	var blocks []chat.Block
// 	// // 	blocks = append(blocks, *chat.NewHeaderBlock(
// 	// // 		chat.NewMarkdownObject(fmt.Sprintf(":champagne: Alright! Everyone accepted the offer!")),
// 	// // 		nil))
// 	// //
// 	// // 	linkBlock := chat.NewSectionBlock(
// 	// // 		chat.NewMarkdownObject(
// 	// // 			fmt.Sprintf("It's now time to do your thing! Once you've kept up with your side of the bargain, "+
// 	// // 				"just mark it as 'completed' by going into your [transactions](/transactions). Also, when other parties "+
// 	// // 				"give you what was agreed upon, you also have to confirm it, so we can exchange hours from your timebanks.")),
// 	// // 		nil,
// 	// // 		nil,
// 	// // 		nil)
// 	// // 	blocks = append(blocks, *linkBlock)
// 	// //
// 	// // 	sendMessage := chat.NewSendConversationMessage(
// 	// // 		loggedInUserKey,
// 	// // 		user.GetUsername(),
// 	// // 		userKeys,
// 	// // 		"",
// 	// // 		blocks,
// 	// // 		[]chat.Attachment{},
// 	// // 		nil,
// 	// // 	)
// 	// // 	_, err = t.chatService.SendConversationMessage(ctx, sendMessage)
// 	// // 	if err != nil {
// 	// // 		l.Error("could not send conversation message", zap.Error(err))
// 	// // 		return nil, err
// 	// // 	}
// 	// //
// 	// // }
// 	//
// 	// //
// 	// // var resources = resource.NewEmptyResources()
// 	// // if len(offerItems.Items) > 0 {
// 	// // 	getResourcesByKeysResponse, err := t.rs.GetByKeys(ctx, offerItems.GetResourceKeys())
// 	// // 	if err != nil {
// 	// // 		l.Error("could not get resources by keys", zap.Error(err))
// 	// // 		return nil, err
// 	// // 	}
// 	// // 	resources = getResourcesByKeysResponse
// 	// // }
// 	// //
// 	// // if currentUserLastOneToDecide {
// 	// //
// 	// // 	for _, user := range users.Items {
// 	// // 		userKey := user.GetUserKey()
// 	// //
// 	// // 		var blocks []chat.Block
// 	// //
// 	// // 		userItems := offerItems.GetOfferItemsReceivedByUser(user.GetUserKey())
// 	// // 		for _, userItem := range userItems.Items {
// 	// // 			userItemKey := userItem.GetKey()
// 	// // 			userItemIdValue := userItemKey.ID.String()
// 	// // 			if userItem.IsReceivedBy(user.GetUserKey()) {
// 	// //
// 	// // 				if !userItem.IsResourceExchangeItem() {
// 	// // 					continue
// 	// // 				}
// 	// //
// 	// // 				actionId := "confirm_item_received"
// 	// // 				res, _ := resources.GetResource(userItem.GetResourceKey())
// 	// // 				fromUser, _ := users.GetUser(userItem.GetFromUserKey())
// 	// // 				block := chat.NewSectionBlock(
// 	// // 					chat.NewMarkdownObject(fmt.Sprintf("**%s**", res.Summary)),
// 	// // 					[]chat.BlockElement{},
// 	// // 					chat.NewButtonElement(
// 	// // 						chat.NewMarkdownObject(fmt.Sprintf("I received it from **%s**", fromUser.Username)),
// 	// // 						nil,
// 	// // 						&actionId,
// 	// // 						nil,
// 	// // 						&userItemIdValue,
// 	// // 						nil),
// 	// // 					nil)
// 	// // 				blocks = append(blocks, *block)
// 	// //
// 	// // 			} else if userItem.IsGivenBy(user.GetUserKey()) {
// 	// // 				if !userItem.IsResourceExchangeItem() {
// 	// // 					continue
// 	// // 				}
// 	// //
// 	// // 				actionId := "confirm_item_given"
// 	// // 				res, _ := resources.GetResource(userItem.GetResourceKey())
// 	// // 				toUser, _ := users.GetUser(userItem.GetToUserKey())
// 	// // 				block := chat.NewSectionBlock(
// 	// // 					chat.NewMarkdownObject(fmt.Sprintf("**%s**", res.Summary)),
// 	// // 					[]chat.BlockElement{},
// 	// // 					chat.NewButtonElement(
// 	// // 						chat.NewMarkdownObject(fmt.Sprintf("I've given it to **%s**", toUser.Username)),
// 	// // 						nil,
// 	// // 						&actionId,
// 	// // 						nil,
// 	// // 						&userItemIdValue,
// 	// // 						nil),
// 	// // 					nil)
// 	// // 				blocks = append(blocks, *block)
// 	// //
// 	// // 			}
// 	// //
// 	// // 			blocks = append(blocks, *chat.NewDividerBlock())
// 	// // 		}
// 	// //
// 	// // 		if blocks == nil {
// 	// // 			blocks = []chat.Block{}
// 	// // 		}
// 	// //
// 	// // 		if blocks[len(blocks)-1].Type == chat.Divider {
// 	// // 			blocks = blocks[:len(blocks)-1]
// 	// // 		}
// 	// //
// 	// // 		sendMessage := chat.NewSendConversationMessage(
// 	// // 			user.GetUserKey(),
// 	// // 			user.GetUsername(),
// 	// // 			userKeys,
// 	// // 			"",
// 	// // 			blocks,
// 	// // 			[]chat.Attachment{},
// 	// // 			&userKey,
// 	// // 		)
// 	// // 		_, err = t.chatService.SendConversationMessage(nil, sendMessage)
// 	// // 		if err != nil {
// 	// // 			l.Error("could not send conversation message", zap.Error(err))
// 	// // 			return nil, err
// 	// // 		}
// 	// // 	}
// 	// // }
// 	//
// 	// return nil
// }

// func (t TradingService) notifyItemGivenOrReceived(ctx context.Context, offerItemBeingConfirmed *trading.OfferItem, confirmingUser *auth.User, concernedOfferUsers *auth.Users) error {

// l := logging.WithContext(ctx)
//
// l.Debug("getting offer item resource")
//
// // confirming items is only for "Resource" offer items, so it's safe to assume that
// // the item.resourceKey is not going to be nil
// getResource := t.rs.GetByKey(ctx, resource.NewGetResourceByKeyQuery(offerItemBeingConfirmed.GetResourceKey()))
// if getResource.Error != nil {
// 	l.Error("could not get offer item resource", zap.Error(getResource.Error))
// 	return getResource.Error
// }
// resourceSummary := getResource.Resource.Summary
//
// offerItemFromUserKey := offerItemBeingConfirmed.GetFromUserKey()
// offerItemToUserKey := offerItemBeingConfirmed.GetToUserKey()
//
// // building sentence component for sending message
// var verb string
// var article string
// var otherUserName string
// if offerItemBeingConfirmed.IsGivenBy(confirmingUser.GetUserKey()) {
// 	verb = "given"
// 	article = "to"
//
// 	toUser, err := concernedOfferUsers.GetUser(offerItemToUserKey)
// 	if err != nil {
// 		l.Error("could not get 'to' user", zap.Error(err))
// 		return err
// 	}
//
// 	otherUserName = toUser.Username
//
// } else if offerItemBeingConfirmed.IsReceivedBy(confirmingUser.GetUserKey()) {
// 	verb = "received"
// 	article = "from"
//
// 	fromUser, err := concernedOfferUsers.GetUser(offerItemFromUserKey)
// 	if err != nil {
// 		l.Error("could not get 'from' user", zap.Error(err))
// 		return err
// 	}
//
// 	otherUserName = fromUser.Username
//
// }
//
// _, err := t.chatService.SendConversationMessage(ctx, chat.NewSendConversationMessage(
// 	confirmingUser.GetUserKey(),
// 	confirmingUser.Username,
// 	concernedOfferUsers.GetUserKeys(),
// 	"",
// 	[]chat.Block{
// 		*chat.NewHeaderBlock(chat.NewMarkdownObject(
// 			fmt.Sprintf(":heavy_check_mark: **%s** has confirmed having %s **%s** %s **%s**",
// 				confirmingUser.Username,
// 				verb,
// 				resourceSummary,
// 				article,
// 				otherUserName,
// 			),
// 		),
// 			nil),
// 	},
// 	[]chat.Attachment{},
// 	nil,
// ))
//
// if err != nil {
// 	l.Error("could not send message to users")
// 	return err
// }
//
// 	return nil
// }

// //
// func (t TradingService) findAppropriateChannelForOffer(offer *trading.Offer, offerItems *domain.OfferItems, approvers *keys.UserKeys) (keys.ChannelKey, chat.ChannelType, error) {
// 	if !offerItems.GetGroupKeys().IsEmpty() {
// 		channelKey := keys.GetChannelKeyForGroup(offer.GroupKey)
// 		return channelKey, chat.GroupChannel, nil
// 	}
// 	channelKey, err := keys.GetChannelKey(approvers)
// 	if err != nil {
// 		return keys.ChannelKey{}, 0, err
// 	}
// 	return channelKey, chat.ConversationChannel, nil
// }
//
// func (t TradingService) sendCustomOfferMessage(ctx context.Context, fromUser *models.UserSession, userKeys *keys.UserKeys, message string) error {
//
// 	if strings.TrimSpace(message) == "" {
// 		return nil
// 	}
//
// 	sendMsgRequest := service.NewSendConversationMessage(
// 		fromUser.GetUserKey(),
// 		fromUser.GetUsername(),
// 		userKeys,
// 		message,
// 		[]chat.Block{},
// 		[]chat.Attachment{},
// 		nil,
// 	)
//
// 	if _, err := t.chatService.SendConversationMessage(ctx, sendMsgRequest); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (t TradingService) sendAcceptOrDeclineMessages(ctx context.Context, offerItems *domain.OfferItems, offer *trading.Offer, userSession *models.UserSession) error {
//
// 	approvers, err := t.tradingStore.FindApproversForCandidateOffer(offer, offerItems)
// 	if err != nil {
// 		return err
// 	}
//
// 	for _, userKey := range approvers.Items {
// 		chatMessage := t.buildAcceptOrDeclineChatMessage(userKey, offer, offerItems)
// 		sendMsgRequest := service.NewSendConversationMessage(
// 			userSession.GetUserKey(),
// 			userSession.GetUsername(),
// 			approvers,
// 			"New offer",
// 			chatMessage,
// 			[]chat.Attachment{},
// 			&userKey,
// 		)
// 		_, err := t.chatService.SendConversationMessage(ctx, sendMsgRequest)
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// func (t TradingService) buildAcceptOrDeclineChatMessage(recipientUserKey keys.UserKey, offer *trading.Offer, offerItems *domain.OfferItems) []chat.Block {
//
// 	messageBlocks := []chat.Block{
// 		*chat.NewHeaderBlock(
// 			chat.NewMarkdownObject(
// 				fmt.Sprintf("%s is proposing an exchange", offer.GetAuthorKey().GetFrontendLink()),
// 			), nil),
// 	}
//
// 	for _, offerItem := range offerItems.Items {
//
// 		var message string
//
// 		if offerItem.IsResourceTransfer() {
//
// 			resourceTransfer := offerItem.(*domain.ResourceTransferItem)
//
// 			if resourceTransfer.To.IsForUser() {
//
// 				message = fmt.Sprintf("%s would take %s",
// 					resourceTransfer.To.GetUserKey().GetFrontendLink(),
// 					resourceTransfer.ResourceKey.GetFrontendLink(),
// 				)
//
// 			} else if resourceTransfer.To.IsForGroup() {
//
// 				message = fmt.Sprintf("The group %s would take %s",
// 					resourceTransfer.To.GetGroupKey().GetFrontendLink(),
// 					resourceTransfer.ResourceKey.GetFrontendLink(),
// 				)
//
// 			}
//
// 		} else if offerItem.IsServiceProviding() {
//
// 			serviceProvision := offerItem.(*domain.ProvideServiceItem)
//
// 			if serviceProvision.To.IsForGroup() {
//
// 				message = fmt.Sprintf("group %s would get %s worth of %s",
// 					serviceProvision.To.GetGroupKey().GetFrontendLink(),
// 					serviceProvision.Duration.String(),
// 					serviceProvision.ResourceKey.GetFrontendLink(),
// 				)
//
// 			} else if serviceProvision.To.IsForUser() {
//
// 				message = fmt.Sprintf("user %s would get %s worth of %s",
// 					serviceProvision.To.GetUserKey().GetFrontendLink(),
// 					serviceProvision.Duration.String(),
// 					serviceProvision.ResourceKey.GetFrontendLink(),
// 				)
//
// 			}
//
// 		} else if offerItem.IsBorrowingResource() {
//
// 			resourceBorrow := offerItem.(*domain.BorrowResourceItem)
//
// 			if resourceBorrow.To.IsForUser() {
//
// 				message = fmt.Sprintf("user %s would borrow %s for %s",
// 					resourceBorrow.To.GetUserKey().GetFrontendLink(),
// 					resourceBorrow.ResourceKey.GetFrontendLink(),
// 					resourceBorrow.Duration.String(),
// 				)
//
// 			} else if resourceBorrow.To.IsForGroup() {
//
// 				message = fmt.Sprintf("group %s would borrow %s for %s",
// 					resourceBorrow.To.GetGroupKey().GetFrontendLink(),
// 					resourceBorrow.ResourceKey.GetFrontendLink(),
// 					resourceBorrow.Duration.String(),
// 				)
//
// 			}
//
// 		} else if offerItem.IsCreditTransfer() {
//
// 			creditTransfer := offerItem.(*domain.CreditTransferItem)
//
// 			fromLink := ""
// 			if creditTransfer.From.IsForGroup() {
// 				fromLink = creditTransfer.From.GetGroupKey().GetFrontendLink()
// 			} else if creditTransfer.From.IsForUser() {
// 				fromLink = creditTransfer.From.GetUserKey().GetFrontendLink()
// 			}
//
// 			toLink := ""
// 			if creditTransfer.To.IsForGroup() {
// 				toLink = "group " + creditTransfer.To.GetGroupKey().GetFrontendLink()
// 			} else if creditTransfer.To.IsForUser() {
// 				toLink = "user " + creditTransfer.To.GetUserKey().GetFrontendLink()
// 			}
//
// 			message = fmt.Sprintf("user %s would get `%s` of time credits from %s",
// 				toLink,
// 				creditTransfer.Amount.String(),
// 				fromLink,
// 			)
//
// 		}
//
// 		messageBlocks = append(messageBlocks, *chat.NewSectionBlock(chat.NewMarkdownObject(message), nil, nil, nil))
//
// 	}
//
// 	primaryButtonStyle := chat.Primary
// 	dangerButtonStyle := chat.Danger
// 	acceptOfferActionId := "accept_offer"
// 	declineOfferActionId := "decline_offer"
// 	offerId := offer.GetKey().String()
//
// 	messageBlocks = append(messageBlocks, *chat.NewActionBlock([]chat.BlockElement{
// 		*chat.NewButtonElement(chat.NewPlainTextObject("Accept"), &primaryButtonStyle, &acceptOfferActionId, nil, &offerId, nil),
// 		*chat.NewButtonElement(chat.NewPlainTextObject("Decline"), &dangerButtonStyle, &declineOfferActionId, nil, &offerId, nil),
// 	}, nil))
//
// 	linkBlock := chat.NewSectionBlock(
// 		chat.NewMarkdownObject(
// 			fmt.Sprintf("[View offer details](/offers/%s)", offerId)),
// 		nil,
// 		nil,
// 		nil)
//
// 	messageBlocks = append(messageBlocks, *linkBlock)
//
// 	return messageBlocks
//
// }
//
// func assertTimeOfferItemsHavePositiveTimeValue(offerItems *domain.OfferItems) error {
// 	for _, offerItem := range offerItems.Items {
//
// 		var duration time.Duration
// 		if offerItem.IsCreditTransfer() {
// 			duration = offerItem.(*domain.CreditTransferItem).Amount
// 		} else if offerItem.IsBorrowingResource() {
// 			duration = offerItem.(*domain.BorrowResourceItem).Duration
// 		} else if offerItem.IsServiceProviding() {
// 			duration = offerItem.(*domain.ProvideServiceItem).Duration
// 		} else {
// 			continue
// 		}
//
// 		if duration < 0 {
// 			return exceptions.ErrNegativeDuration
// 		}
// 	}
// 	return nil
// }
//
// func assertResourcesAreTransferredOnlyOnce(offerItems *domain.OfferItems) error {
// 	var seenResourceKeys []keys.ResourceKey
// 	for _, item := range offerItems.Items {
// 		if item.IsResourceTransfer() {
// 			resourceTransfer := item.(*domain.ResourceTransferItem)
// 			resourceKey := resourceTransfer.ResourceKey
// 			for _, seenResourceKey := range seenResourceKeys {
// 				if seenResourceKey == resourceKey {
// 					return exceptions.ErrDuplicateResourceInOffer
// 				}
// 			}
// 			seenResourceKeys = append(seenResourceKeys, resourceKey)
// 		}
// 	}
// 	return nil
// }

/**
func (t TradingService) buildOfferCompletedMessage(ctx context.Context, items *tradingdomain.OfferItems, users *models.Users) ([]chat.Block, string, error) {

	var blocks []chat.Block

	mainText := ":champagne: Alright! everybody confirmed having received and given their stuff."
	blocks = append(blocks, *chat.NewHeaderBlock(
		chat.NewMarkdownObject(mainText),
		nil,
	))

	for _, offerItem := range items.Items {

		if offerItem.IsCreditTransfer() {

			creditTransfer := offerItem.(*tradingdomain.CreditTransferItem)

			var toLink = ""
			var fromLink = ""

			if creditTransfer.To.IsForGroup() {
				toLink = creditTransfer.To.GetGroupKey().GetFrontendLink()
			} else if creditTransfer.To.IsForUser() {
				toLink = creditTransfer.To.GetUserKey().GetFrontendLink()
			}

			if creditTransfer.From.IsForGroup() {
				fromLink = creditTransfer.From.GetGroupKey().GetFrontendLink()
			} else if creditTransfer.From.IsForUser() {
				fromLink = creditTransfer.From.GetUserKey().GetFrontendLink()
			}

			blocks = append(blocks, *chat.NewSectionBlock(
				chat.NewMarkdownObject(fmt.Sprintf("%s received `%s` timebank credits from %s",
					toLink,
					creditTransfer.Amount.Truncate(time.Minute*1).String(),
					fromLink,
				)),
				nil,
				nil,
				nil,
			))
		}
	}

	return blocks, mainText, nil

}

*/
