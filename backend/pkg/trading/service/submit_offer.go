package service

import (
	"fmt"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	chatmodel "github.com/commonpool/backend/pkg/chat/model"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/group/model"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/resource"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	model2 "github.com/commonpool/backend/pkg/trading/model"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"strings"
	"time"
)

func (t TradingService) SendOffer(ctx context.Context, groupKey model.GroupKey, offerItems *model2.OfferItems, message string) (*model2.Offer, *model2.OfferItems, error) {

	userSession, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return nil, nil, exceptions.ErrUnauthorized
	}

	// The ownership of a resource can only be moved once in an offer
	if err := assertResourcesAreTransferredOnlyOnce(offerItems); err != nil {
		return nil, nil, err
	}

	// All time based offers must have positive time values. No negative time
	if err := assertTimeOfferItemsHavePositiveTimeValue(offerItems); err != nil {
		return nil, nil, err
	}

	resources, err := t.rs.GetByKeys(ctx, offerItems.GetResourceKeys())
	if err != nil {
		return nil, nil, err
	}

	// Checking that the offer's group can actually 'see' the resource
	if err := t.assertResourcesAreViewableByGroup(resources, groupKey); err != nil {
		return nil, nil, err
	}

	// Checking that an offer does not include transferring a resource to its current owner
	if err := t.assertResourcesAreNotTransferredToTheirCurrentOwner(resources, offerItems); err != nil {
		return nil, nil, err
	}

	// Checking that resource_transfer offerItems refer to object-typed resources
	if err := t.assertResourceTransferOfferItemsReferToObjectResources(resources, offerItems); err != nil {
		return nil, nil, err
	}

	// Checking that service_provision offer-items actually point to a service-typed resource
	if err := t.assertProvideServiceItemsAreForServiceResources(resources, offerItems); err != nil {
		return nil, nil, err
	}

	// Checking that borrowal offer-items actually point to a object-typed resource
	if err := t.assertBorrowOfferItemPointToObjectTypedResource(resources, offerItems); err != nil {
		return nil, nil, err
	}

	offerKey := model2.NewOfferKey(uuid.NewV4())
	offer := model2.NewOffer(offerKey, groupKey, userSession.GetUserKey(), message, nil)

	err = t.tradingStore.SaveOffer(offer, offerItems)
	if err != nil {
		return nil, nil, err
	}

	offer, err = t.tradingStore.GetOffer(offerKey)
	if err != nil {
		return nil, nil, err
	}

	offerItems, err = t.tradingStore.GetOfferItemsForOffer(offerKey)
	if err != nil {
		return nil, nil, err
	}

	if err := t.sendAcceptOrDeclineMessages(ctx, offerItems, offer, userSession); err != nil {
		return nil, nil, err
	}

	if err := t.sendCustomOfferMessage(ctx, userSession, offerItems.GetUserKeys(), message); err != nil {
		return nil, nil, err
	}

	return offer, offerItems, nil
}

func (t TradingService) findAppropriateChannelForOffer(offer *tradingmodel.Offer, offerItems *tradingmodel.OfferItems, approvers *usermodel.UserKeys) (chatmodel.ChannelKey, chatmodel.ChannelType, error) {
	if !offerItems.GetGroupKeys().IsEmpty() {
		channelKey := chatmodel.GetChannelKeyForGroup(offer.GroupKey)
		return channelKey, chatmodel.GroupChannel, nil
	}
	channelKey, err := chatmodel.GetChannelKey(approvers)
	if err != nil {
		return chatmodel.ChannelKey{}, 0, err
	}
	return channelKey, chatmodel.ConversationChannel, nil
}

func (t TradingService) assertResourcesAreViewableByGroup(resources *resource.GetResourceByKeysResponse, groupKey groupmodel.GroupKey) error {
	for _, item := range resources.Resources.Items {
		if !resources.Claims.GroupHasClaim(groupKey, item.Key, resourcemodel.ViewerClaim) {
			return exceptions.ErrResourceNotSharedWithGroup
		}
	}
	return nil
}

func (t TradingService) assertResourcesAreNotTransferredToTheirCurrentOwner(resources *resource.GetResourceByKeysResponse, items *tradingmodel.OfferItems) error {
	for _, offerItem := range items.Items {
		if !offerItem.IsResourceTransfer() {
			continue
		}
		resourceTransfer := offerItem.(*tradingmodel.ResourceTransferItem)
		if resources.Claims.HasClaim(resourceTransfer.To, resourceTransfer.ResourceKey, resourcemodel.OwnershipClaim) {
			return exceptions.ErrCannotTransferResourceToItsOwner
		}
	}
	return nil
}

func (t TradingService) assertResourceTransferOfferItemsReferToObjectResources(resources *resource.GetResourceByKeysResponse, items *tradingmodel.OfferItems) error {
	for _, offerItem := range items.Items {
		if !offerItem.IsResourceTransfer() {
			continue
		}
		resourceTransfer := offerItem.(*tradingmodel.ResourceTransferItem)
		r, err := resources.Resources.GetResource(resourceTransfer.ResourceKey)
		if err != nil {
			return err
		}
		if !r.IsObject() {
			return exceptions.ErrResourceTransferOfferItemsMustReferToObjectResources
		}
	}
	return nil
}

func (t TradingService) assertProvideServiceItemsAreForServiceResources(resources *resource.GetResourceByKeysResponse, items *tradingmodel.OfferItems) error {
	for _, offerItem := range items.Items {
		if !offerItem.IsServiceProviding() {
			continue
		}
		serviceProvision := offerItem.(*tradingmodel.ProvideServiceItem)
		r, err := resources.Resources.GetResource(serviceProvision.ResourceKey)
		if err != nil {
			return err
		}
		if !r.IsService() {
			return exceptions.ErrServiceProvisionOfferItemsMustPointToServiceResources
		}
	}
	return nil
}

func (t TradingService) assertBorrowOfferItemPointToObjectTypedResource(resources *resource.GetResourceByKeysResponse, items *tradingmodel.OfferItems) error {
	for _, offerItem := range items.Items {
		if !offerItem.IsBorrowingResource() {
			continue
		}
		itemBorrow := offerItem.(*tradingmodel.BorrowResourceItem)
		r, err := resources.Resources.GetResource(itemBorrow.ResourceKey)
		if err != nil {
			return err
		}
		if !r.IsObject() {
			return exceptions.ErrBorrowOfferItemMustReferToObjectTypedResource
		}
	}
	return nil
}

func (t TradingService) sendCustomOfferMessage(ctx context.Context, fromUser *auth.UserSession, userKeys *usermodel.UserKeys, message string) error {

	if strings.TrimSpace(message) == "" {
		return nil
	}

	sendMsgRequest := chat.NewSendConversationMessage(
		fromUser.GetUserKey(),
		fromUser.GetUsername(),
		userKeys,
		message,
		[]chatmodel.Block{},
		[]chatmodel.Attachment{},
		nil,
	)

	if _, err := t.chatService.SendConversationMessage(ctx, sendMsgRequest); err != nil {
		return err
	}

	return nil
}

func (t TradingService) sendAcceptOrDeclineMessages(ctx context.Context, offerItems *tradingmodel.OfferItems, offer *tradingmodel.Offer, userSession *auth.UserSession) error {

	approvers, err := t.tradingStore.FindApproversForCandidateOffer(offer, offerItems)
	if err != nil {
		return err
	}

	for _, userKey := range approvers.Items {
		chatMessage := t.buildAcceptOrDeclineChatMessage(userKey, offer, offerItems)
		sendMsgRequest := chat.NewSendConversationMessage(
			userSession.GetUserKey(),
			userSession.GetUsername(),
			approvers,
			"New offer",
			chatMessage,
			[]chatmodel.Attachment{},
			&userKey,
		)
		_, err := t.chatService.SendConversationMessage(ctx, sendMsgRequest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t TradingService) buildAcceptOrDeclineChatMessage(recipientUserKey usermodel.UserKey, offer *tradingmodel.Offer, offerItems *tradingmodel.OfferItems) []chatmodel.Block {

	messageBlocks := []chatmodel.Block{
		*chatmodel.NewHeaderBlock(
			chatmodel.NewMarkdownObject(
				fmt.Sprintf("%s is proposing an exchange", t.chatService.GetUserLink(offer.GetAuthorKey())),
			), nil),
	}

	for _, offerItem := range offerItems.Items {

		var message string

		if offerItem.IsResourceTransfer() {

			resourceTransfer := offerItem.(*tradingmodel.ResourceTransferItem)

			if resourceTransfer.To.IsForUser() {

				message = fmt.Sprintf("%s would take %s",
					t.chatService.GetUserLink(resourceTransfer.To.GetUserKey()),
					t.chatService.GetResourceLink(resourceTransfer.ResourceKey),
				)

			} else if resourceTransfer.To.IsForGroup() {

				message = fmt.Sprintf("The group %s would take %s",
					t.chatService.GetGroupLink(resourceTransfer.To.GetGroupKey()),
					t.chatService.GetResourceLink(resourceTransfer.ResourceKey),
				)

			}

		} else if offerItem.IsServiceProviding() {

			serviceProvision := offerItem.(*tradingmodel.ProvideServiceItem)

			if serviceProvision.To.IsForGroup() {

				message = fmt.Sprintf("group %s would get %s worth of %s",
					t.chatService.GetGroupLink(serviceProvision.To.GetGroupKey()),
					serviceProvision.Duration.String(),
					t.chatService.GetResourceLink(serviceProvision.ResourceKey),
				)

			} else if serviceProvision.To.IsForUser() {

				message = fmt.Sprintf("user %s would get %s worth of %s",
					t.chatService.GetUserLink(serviceProvision.To.GetUserKey()),
					serviceProvision.Duration.String(),
					t.chatService.GetResourceLink(serviceProvision.ResourceKey),
				)

			}

		} else if offerItem.IsBorrowingResource() {

			resourceBorrow := offerItem.(*tradingmodel.BorrowResourceItem)

			if resourceBorrow.To.IsForUser() {

				message = fmt.Sprintf("user %s would borrow %s for %s",
					t.chatService.GetUserLink(resourceBorrow.To.GetUserKey()),
					t.chatService.GetResourceLink(resourceBorrow.ResourceKey),
					resourceBorrow.Duration.String(),
				)

			} else if resourceBorrow.To.IsForGroup() {

				message = fmt.Sprintf("group %s would borrow %s for %s",
					t.chatService.GetGroupLink(resourceBorrow.To.GetGroupKey()),
					t.chatService.GetResourceLink(resourceBorrow.ResourceKey),
					resourceBorrow.Duration.String(),
				)

			}

		} else if offerItem.IsCreditTransfer() {

			creditTransfer := offerItem.(*tradingmodel.CreditTransferItem)

			fromLink := ""
			if creditTransfer.From.IsForGroup() {
				fromLink = t.chatService.GetGroupLink(creditTransfer.From.GetGroupKey())
			} else if creditTransfer.From.IsForUser() {
				fromLink = t.chatService.GetUserLink(creditTransfer.From.GetUserKey())
			}

			toLink := ""
			if creditTransfer.To.IsForGroup() {
				toLink = "group " + t.chatService.GetGroupLink(creditTransfer.To.GetGroupKey())
			} else if creditTransfer.To.IsForUser() {
				toLink = "user " + t.chatService.GetUserLink(creditTransfer.To.GetUserKey())
			}

			message = fmt.Sprintf("user %s would get `%s` of time credits from %s",
				toLink,
				creditTransfer.Amount.String(),
				fromLink,
			)

		}

		messageBlocks = append(messageBlocks, *chatmodel.NewSectionBlock(chatmodel.NewMarkdownObject(message), nil, nil, nil))

	}

	primaryButtonStyle := chatmodel.Primary
	dangerButtonStyle := chatmodel.Danger
	acceptOfferActionId := "accept_offer"
	declineOfferActionId := "decline_offer"
	offerId := offer.GetKey().String()

	messageBlocks = append(messageBlocks, *chatmodel.NewActionBlock([]chatmodel.BlockElement{
		*chatmodel.NewButtonElement(chatmodel.NewPlainTextObject("Accept"), &primaryButtonStyle, &acceptOfferActionId, nil, &offerId, nil),
		*chatmodel.NewButtonElement(chatmodel.NewPlainTextObject("Decline"), &dangerButtonStyle, &declineOfferActionId, nil, &offerId, nil),
	}, nil))

	linkBlock := chatmodel.NewSectionBlock(
		chatmodel.NewMarkdownObject(
			fmt.Sprintf("[View offer details](/offers/%s)", offerId)),
		nil,
		nil,
		nil)

	messageBlocks = append(messageBlocks, *linkBlock)

	return messageBlocks

}

func assertTimeOfferItemsHavePositiveTimeValue(offerItems *tradingmodel.OfferItems) error {
	for _, offerItem := range offerItems.Items {

		var duration time.Duration
		if offerItem.IsCreditTransfer() {
			duration = offerItem.(*tradingmodel.CreditTransferItem).Amount
		} else if offerItem.IsBorrowingResource() {
			duration = offerItem.(*tradingmodel.BorrowResourceItem).Duration
		} else if offerItem.IsServiceProviding() {
			duration = offerItem.(*tradingmodel.ProvideServiceItem).Duration
		} else {
			continue
		}

		if duration < 0 {
			return exceptions.ErrNegativeDuration
		}
	}
	return nil
}

func assertResourcesAreTransferredOnlyOnce(offerItems *tradingmodel.OfferItems) error {
	var seenResourceKeys []resourcemodel.ResourceKey
	for _, item := range offerItems.Items {
		if item.IsResourceTransfer() {
			resourceTransfer := item.(*tradingmodel.ResourceTransferItem)
			resourceKey := resourceTransfer.ResourceKey
			for _, seenResourceKey := range seenResourceKeys {
				if seenResourceKey == resourceKey {
					return exceptions.ErrDuplicateResourceInOffer
				}
			}
			seenResourceKeys = append(seenResourceKeys, resourceKey)
		}
	}
	return nil
}