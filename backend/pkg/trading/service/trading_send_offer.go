package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/chat"
	resource2 "github.com/commonpool/backend/pkg/resource"
	trading2 "github.com/commonpool/backend/pkg/trading"
	uuid "github.com/satori/go.uuid"
	ctx "golang.org/x/net/context"
	"strings"
	"time"
)

func (t TradingService) SendOffer(ctx ctx.Context, groupKey model.GroupKey, offerItems *trading2.OfferItems, message string) (*trading2.Offer, *trading2.OfferItems, error) {

	userSession, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return nil, nil, errors.ErrUnauthorized
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

	offerKey := model.NewOfferKey(uuid.NewV4())
	offer := trading2.NewOffer(offerKey, groupKey, userSession.GetUserKey(), message, nil)

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

func (t TradingService) findAppropriateChannelForOffer(offer *trading2.Offer, offerItems *trading2.OfferItems, approvers *model.UserKeys) (model.ChannelKey, chat.ChannelType, error) {
	if !offerItems.GetGroupKeys().IsEmpty() {
		return offer.GroupKey.GetChannelKey(), chat.GroupChannel, nil
	}
	channelKey, err := approvers.GetChannelKey()
	if err != nil {
		return model.ChannelKey{}, 0, err
	}
	return channelKey, chat.ConversationChannel, nil
}

func (t TradingService) assertResourcesAreViewableByGroup(resources *resource2.GetResourceByKeysResponse, groupKey model.GroupKey) error {
	for _, item := range resources.Resources.Items {
		if !resources.Claims.GroupHasClaim(groupKey, item.Key, resource2.ViewerClaim) {
			return errors.ErrResourceNotSharedWithGroup
		}
	}
	return nil
}

func (t TradingService) assertResourcesAreNotTransferredToTheirCurrentOwner(resources *resource2.GetResourceByKeysResponse, items *trading2.OfferItems) error {
	for _, offerItem := range items.Items {
		if !offerItem.IsResourceTransfer() {
			continue
		}
		resourceTransfer := offerItem.(*trading2.ResourceTransferItem)
		if resources.Claims.HasClaim(resourceTransfer.To, resourceTransfer.ResourceKey, resource2.OwnershipClaim) {
			return errors.ErrCannotTransferResourceToItsOwner
		}
	}
	return nil
}

func (t TradingService) assertResourceTransferOfferItemsReferToObjectResources(resources *resource2.GetResourceByKeysResponse, items *trading2.OfferItems) error {
	for _, offerItem := range items.Items {
		if !offerItem.IsResourceTransfer() {
			continue
		}
		resourceTransfer := offerItem.(*trading2.ResourceTransferItem)
		r, err := resources.Resources.GetResource(resourceTransfer.ResourceKey)
		if err != nil {
			return err
		}
		if !r.IsObject() {
			return errors.ErrResourceTransferOfferItemsMustReferToObjectResources
		}
	}
	return nil
}

func (t TradingService) assertProvideServiceItemsAreForServiceResources(resources *resource2.GetResourceByKeysResponse, items *trading2.OfferItems) error {
	for _, offerItem := range items.Items {
		if !offerItem.IsServiceProviding() {
			continue
		}
		serviceProvision := offerItem.(*trading2.ProvideServiceItem)
		r, err := resources.Resources.GetResource(serviceProvision.ResourceKey)
		if err != nil {
			return err
		}
		if !r.IsService() {
			return errors.ErrServiceProvisionOfferItemsMustPointToServiceResources
		}
	}
	return nil
}

func (t TradingService) assertBorrowOfferItemPointToObjectTypedResource(resources *resource2.GetResourceByKeysResponse, items *trading2.OfferItems) error {
	for _, offerItem := range items.Items {
		if !offerItem.IsBorrowingResource() {
			continue
		}
		itemBorrow := offerItem.(*trading2.BorrowResourceItem)
		r, err := resources.Resources.GetResource(itemBorrow.ResourceKey)
		if err != nil {
			return err
		}
		if !r.IsObject() {
			return errors.ErrBorrowOfferItemMustReferToObjectTypedResource
		}
	}
	return nil
}

func (t TradingService) sendCustomOfferMessage(ctx context.Context, fromUser *auth.UserSession, userKeys *model.UserKeys, message string) error {

	if strings.TrimSpace(message) == "" {
		return nil
	}

	sendMsgRequest := chat.NewSendConversationMessage(
		fromUser.GetUserKey(),
		fromUser.GetUsername(),
		userKeys,
		message,
		[]chat.Block{},
		[]chat.Attachment{},
		nil,
	)

	if _, err := t.chatService.SendConversationMessage(ctx, sendMsgRequest); err != nil {
		return err
	}

	return nil
}

func (t TradingService) sendAcceptOrDeclineMessages(ctx context.Context, offerItems *trading2.OfferItems, offer *trading2.Offer, userSession *auth.UserSession) error {

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
			[]chat.Attachment{},
			&userKey,
		)
		_, err := t.chatService.SendConversationMessage(ctx, sendMsgRequest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t TradingService) buildAcceptOrDeclineChatMessage(recipientUserKey model.UserKey, offer *trading2.Offer, offerItems *trading2.OfferItems) []chat.Block {

	messageBlocks := []chat.Block{
		*chat.NewHeaderBlock(
			chat.NewMarkdownObject(
				fmt.Sprintf("%s is proposing an exchange", t.chatService.GetUserLink(offer.GetAuthorKey())),
			), nil),
	}

	for _, offerItem := range offerItems.Items {

		var message string

		if offerItem.IsResourceTransfer() {

			resourceTransfer := offerItem.(*trading2.ResourceTransferItem)

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

			serviceProvision := offerItem.(*trading2.ProvideServiceItem)

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

			resourceBorrow := offerItem.(*trading2.BorrowResourceItem)

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

			creditTransfer := offerItem.(*trading2.CreditTransferItem)

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

		messageBlocks = append(messageBlocks, *chat.NewSectionBlock(chat.NewMarkdownObject(message), nil, nil, nil))

	}

	primaryButtonStyle := chat.Primary
	dangerButtonStyle := chat.Danger
	acceptOfferActionId := "accept_offer"
	declineOfferActionId := "decline_offer"
	offerId := offer.GetKey().String()

	messageBlocks = append(messageBlocks, *chat.NewActionBlock([]chat.BlockElement{
		*chat.NewButtonElement(chat.NewPlainTextObject("Accept"), &primaryButtonStyle, &acceptOfferActionId, nil, &offerId, nil),
		*chat.NewButtonElement(chat.NewPlainTextObject("Decline"), &dangerButtonStyle, &declineOfferActionId, nil, &offerId, nil),
	}, nil))

	linkBlock := chat.NewSectionBlock(
		chat.NewMarkdownObject(
			fmt.Sprintf("[View offer details](/offers/%s)", offerId)),
		nil,
		nil,
		nil)

	messageBlocks = append(messageBlocks, *linkBlock)

	return messageBlocks

}

func assertTimeOfferItemsHavePositiveTimeValue(offerItems *trading2.OfferItems) error {
	for _, offerItem := range offerItems.Items {

		var duration time.Duration
		if offerItem.IsCreditTransfer() {
			duration = offerItem.(*trading2.CreditTransferItem).Amount
		} else if offerItem.IsBorrowingResource() {
			duration = offerItem.(*trading2.BorrowResourceItem).Duration
		} else if offerItem.IsServiceProviding() {
			duration = offerItem.(*trading2.ProvideServiceItem).Duration
		} else {
			continue
		}

		if duration < 0 {
			return errors.ErrNegativeDuration
		}
	}
	return nil
}

func assertResourcesAreTransferredOnlyOnce(offerItems *trading2.OfferItems) error {
	var seenResourceKeys []model.ResourceKey
	for _, item := range offerItems.Items {
		if item.IsResourceTransfer() {
			resourceTransfer := item.(*trading2.ResourceTransferItem)
			resourceKey := resourceTransfer.ResourceKey
			for _, seenResourceKey := range seenResourceKeys {
				if seenResourceKey == resourceKey {
					return errors.ErrDuplicateResourceInOffer
				}
			}
			seenResourceKeys = append(seenResourceKeys, resourceKey)
		}
	}
	return nil
}
