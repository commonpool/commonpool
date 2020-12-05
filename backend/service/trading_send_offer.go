package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
	uuid "github.com/satori/go.uuid"
	context2 "golang.org/x/net/context"
	"net/http"
	"strings"
)

var ErrDuplicateResourceInOffer = errors.NewWebServiceException("resource can only appear once in an offer", "ErrDuplicateResourceInOffer", http.StatusBadRequest)
var ErrNegativeTimeOfferItem = errors.NewWebServiceException("time offers must have positive time value", "ErrNegativeTimeOfferItem", http.StatusBadRequest)
var ErrResourceMustBeTradedByOwner = errors.NewWebServiceException("resource an only be traded by their owner", "ErrResourceMustBeTradedByOwner", http.StatusForbidden)

func (t TradingService) SendOffer(ctx context2.Context, offerItems *trading.OfferItems, message string) (*trading.Offer, *trading.OfferItems, *trading.OfferDecisions, error) {

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		return nil, nil, nil, errors.ErrUnauthorized
	}

	if err := assertResourcesAppearOnlyOnceInOffer(offerItems); err != nil {
		return nil, nil, nil, err
	}
	if err := assertTimeOfferItemsHavePositiveTimeValue(offerItems); err != nil {
		return nil, nil, nil, err
	}

	offerResources, err := t.rs.GetByKeys(ctx, offerItems.GetResourceKeys())
	if err != nil {
		return nil, nil, nil, err
	}

	if err := assertResourcesAreTradedByTheirOwner(offerItems, offerResources); err != nil {
		return nil, nil, nil, err
	}

	offer := trading.NewOffer(model.NewOfferKey(uuid.NewV4()), userSession.GetUserKey(), message, nil)

	decisions := t.createOfferDecisions(offer, offerItems)

	offer, items, decisions, err := t.tradingStore.SaveOffer(offer, offerItems, decisions)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := t.sendAcceptOrDeclineMessages(ctx, offerItems, offer, userSession); err != nil {
		return nil, nil, nil, err
	}

	if err := t.sendCustomOfferMessage(ctx, userSession, offerItems.GetUserKeys(), message); err != nil {
		return nil, nil, nil, err
	}

	return offer, items, decisions, nil
}

func (t TradingService) createOfferDecisions(offer *trading.Offer, offerItems *trading.OfferItems) *trading.OfferDecisions {
	var offerDecisions []*trading.OfferDecision
	for _, userKey := range offerItems.GetUserKeys().Items {
		offerDecision := trading.NewOfferDecision(offer.GetKey(), userKey, trading.PendingDecision)
		offerDecisions = append(offerDecisions, offerDecision)
	}
	return trading.NewOfferDecisions(offerDecisions)
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

func (t TradingService) sendAcceptOrDeclineMessages(ctx context.Context, offerItems *trading.OfferItems, offer *trading.Offer, userSession *auth.UserSession) error {
	userKeys := offerItems.GetUserKeys()
	for _, userKey := range userKeys.Items {
		chatMessage := t.buildAcceptOrDeclineChatMessage(userKey, offer, offerItems)
		sendMsgRequest := chat.NewSendConversationMessage(
			userSession.GetUserKey(),
			userSession.GetUsername(),
			userKeys,
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

func (t TradingService) buildAcceptOrDeclineChatMessage(userKey model.UserKey, offer *trading.Offer, offerItems *trading.OfferItems) []chat.Block {

	messageBlocks := []chat.Block{
		*chat.NewHeaderBlock(
			chat.NewMarkdownObject(
				fmt.Sprintf("%s is proposing an exchange", t.chatService.GetUserLink(offer.GetAuthorKey())),
			), nil),
	}

	for _, offerItem := range offerItems.Items {

		var message string

		if userKey == offerItem.GetFromUserKey() && offerItem.IsTimeExchangeItem() {
			message = fmt.Sprintf("%s would like %s of your time", t.chatService.GetUserLink(offerItem.GetToUserKey()), offerItem.FormatOfferedTimeInSeconds())
		} else if userKey == offerItem.GetToUserKey() && offerItem.IsTimeExchangeItem() {
			message = fmt.Sprintf("you would get %s from %s's timebank", offerItem.FormatOfferedTimeInSeconds(), t.chatService.GetUserLink(offerItem.GetToUserKey()))
		} else if userKey == offerItem.GetFromUserKey() && offerItem.IsResourceExchangeItem() {
			message = fmt.Sprintf("%s would get %s from you", t.chatService.GetUserLink(offerItem.GetToUserKey()), t.chatService.GetResourceLink(offerItem.GetResourceKey()))
		} else if userKey == offerItem.GetToUserKey() && offerItem.IsResourceExchangeItem() {
			message = fmt.Sprintf("you would get %s from %s", t.chatService.GetResourceLink(offerItem.GetResourceKey()), t.chatService.GetUserLink(offerItem.GetToUserKey()))
		} else if offerItem.IsTimeExchangeItem() {
			message = fmt.Sprintf("%s would get %s from %s's timebank", t.chatService.GetUserLink(offerItem.GetToUserKey()), offerItem.FormatOfferedTimeInSeconds(), t.chatService.GetUserLink(offerItem.GetFromUserKey()))
		} else if offerItem.IsResourceExchangeItem() {
			message = fmt.Sprintf("%s would get %s from %s", t.chatService.GetUserLink(offerItem.GetToUserKey()), t.chatService.GetResourceLink(offerItem.GetResourceKey()), t.chatService.GetUserLink(offerItem.GetFromUserKey()))
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

func assertResourcesAreTradedByTheirOwner(offerItems *trading.OfferItems, resources *resource.Resources) error {
	for _, resource := range resources.Items {
		offerItemForResource, _ := offerItems.GetOfferItemInvolvingResource(resource.GetKey())
		if offerItemForResource.GetFromUserKey() != resource.GetOwnerKey() {
			return ErrResourceMustBeTradedByOwner
		}
	}
	return nil
}

func assertTimeOfferItemsHavePositiveTimeValue(offerItems *trading.OfferItems) error {
	for _, offerItem := range offerItems.Items {
		if offerItem.IsTimeExchangeItem() {
			if *offerItem.OfferedTimeInSeconds <= 0 {
				return ErrNegativeTimeOfferItem
			}
		}
	}
	return nil
}

func assertResourcesAppearOnlyOnceInOffer(offerItems *trading.OfferItems) error {
	var seenResourceKeys []model.ResourceKey
	for _, item := range offerItems.Items {
		if item.IsResourceExchangeItem() {
			resourceKey := item.GetResourceKey()
			for _, seenResourceKey := range seenResourceKeys {
				if seenResourceKey == resourceKey {
					return ErrDuplicateResourceInOffer
				}
			}
			seenResourceKeys = append(seenResourceKeys, resourceKey)
		}
	}
	return nil
}
