package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/chat"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/resource"
	trading2 "github.com/commonpool/backend/pkg/trading"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	transaction2 "github.com/commonpool/backend/pkg/transaction"
	"github.com/commonpool/backend/pkg/user"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"time"
)

type TradingService struct {
	tradingStore       trading2.Store
	transactionService transaction2.Service
	groupService       group2.Service
	resourceStore      resource.Store
	userStore          user.Store
	chatService        chat.Service
}

var _ trading2.Service = &TradingService{}

func NewTradingService(
	tradingStore trading2.Store,
	resourceStore resource.Store,
	authStore user.Store,
	chatService chat.Service,
	groupService group2.Service,
	transactionService transaction2.Service) *TradingService {
	return &TradingService{
		tradingStore:       tradingStore,
		resourceStore:      resourceStore,
		userStore:          authStore,
		chatService:        chatService,
		groupService:       groupService,
		transactionService: transactionService,
	}
}

func (t TradingService) checkOfferCompleted(
	ctx context.Context,
	groupKey group2.GroupKey,
	offerKey tradingmodel.OfferKey,
	offerItems *tradingmodel.OfferItems,
	userConfirmingItem usermodel.UserReference,
	usersInOffer *user.Users) error {

	if offerItems.AllPartiesAccepted() && offerItems.AllUserActionsCompleted() {
		for _, offerItem := range offerItems.Items {
			if offerItem.IsCreditTransfer() {
				creditTransfer := offerItem.(*tradingmodel.CreditTransferItem)
				_, err := t.transactionService.TimeCreditsExchanged(groupKey, creditTransfer.From, creditTransfer.To, creditTransfer.Amount)
				if err != nil {
					return err
				}
			}
			if offerItem.IsServiceProviding() {
				serviceProvision := offerItem.(*tradingmodel.ProvideServiceItem)
				_, err := t.transactionService.ServiceWasProvided(groupKey, serviceProvision.ResourceKey, serviceProvision.Duration)
				if err != nil {
					return err
				}
			}
			if offerItem.IsBorrowingResource() {
				borrowResource := offerItem.(*tradingmodel.BorrowResourceItem)
				_, err := t.transactionService.ResourceWasBorrowed(groupKey, borrowResource.ResourceKey, borrowResource.To, borrowResource.Duration)
				if err != nil {
					return err
				}
			}
			if offerItem.IsResourceTransfer() {
				transfer := offerItem.(*tradingmodel.ResourceTransferItem)
				_, err := t.transactionService.ResourceWasTaken(groupKey, transfer.ResourceKey, transfer.To)
				if err != nil {
					return err
				}
			}
		}

		err := t.tradingStore.UpdateOfferStatus(offerKey, tradingmodel.CompletedOffer)
		if err != nil {
			return err
		}

		blocks, mainText, err := t.buildOfferCompletedMessage(ctx, offerItems, usersInOffer)
		if err != nil {
			return err
		}

		_, err = t.chatService.SendConversationMessage(ctx, chat.NewSendConversationMessage(
			userConfirmingItem.GetUserKey(),
			userConfirmingItem.GetUsername(),
			usersInOffer.GetUserKeys(),
			fmt.Sprintf(mainText),
			blocks,
			[]chat.Attachment{},
			nil,
		))
	}
	return nil
}

func (t TradingService) buildOfferCompletedMessage(ctx context.Context, items *tradingmodel.OfferItems, users *user.Users) ([]chat.Block, string, error) {

	var blocks []chat.Block

	mainText := ":champagne: Alright! everybody confirmed having received and given their stuff."
	blocks = append(blocks, *chat.NewHeaderBlock(
		chat.NewMarkdownObject(mainText),
		nil,
	))

	for _, offerItem := range items.Items {

		if offerItem.IsCreditTransfer() {

			creditTransfer := offerItem.(*tradingmodel.CreditTransferItem)

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

func (t TradingService) checkIfAllItemsCompleted(ctx context.Context, loggerInUser usermodel.UserReference, offerItem tradingmodel.OfferItem) error {

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

	allUsersInOffer, err := t.userStore.GetByKeys(ctx, approvers.AllUserKeys())
	if err != nil {
		return err
	}

	return t.checkOfferCompleted(ctx, offer.GroupKey, offer.Key, offerItems, loggerInUser, allUsersInOffer)

}
