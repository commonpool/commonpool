package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	res "github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
	"go.uber.org/zap"
	"time"
)

type TradingService struct {
	tradingStore trading.Store
	groupService group.Service
	rs           res.Store
	us           auth.Store
	chatService  chat.Service
}

var _ trading.Service = &TradingService{}

func NewTradingService(
	tradingStore trading.Store,
	resourceStore res.Store,
	authStore auth.Store,
	chatService chat.Service,
	groupService group.Service) *TradingService {
	return &TradingService{
		tradingStore: tradingStore,
		rs:           resourceStore,
		us:           authStore,
		chatService:  chatService,
		groupService: groupService,
	}
}

func (t TradingService) checkOfferCompleted(ctx context.Context, offerKey model.OfferKey, offerItems *trading.OfferItems, userConfirmingItem *auth.User, usersInOffer *auth.Users) error {

	ctx, l := GetCtx(ctx, "TradingService", "checkOfferCompleted")

	if offerItems.AllUserActionsCompleted() {

		l.Debug("all items have been given and received. Marking offer as completed")

		err := t.tradingStore.UpdateOfferStatus(offerKey, trading.CompletedOffer)
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

		_, err = t.chatService.SendConversationMessage(ctx, chat.NewSendConversationMessage(
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

func (t TradingService) buildOfferCompletedMessage(ctx context.Context, items *trading.OfferItems, users *auth.Users) ([]chat.Block, string, error) {

	ctx, _ = GetCtx(ctx, "TradingService", "buildOfferCompletedMessage")

	var blocks []chat.Block

	mainText := ":champagne: Alright! everybody confirmed having received and given their stuff."
	blocks = append(blocks, *chat.NewHeaderBlock(
		chat.NewMarkdownObject(mainText),
		nil,
	))

	for _, offerItem := range items.Items {

		if offerItem.IsCreditTransfer() {

			creditTransfer := offerItem.(trading.CreditTransferItem)

			var toLink = ""
			var fromLink = ""

			if creditTransfer.To.IsForGroup() {
				toLink = t.chatService.GetGroupLink(creditTransfer.To.GetGroupKey())
			} else if creditTransfer.To.IsForUser() {
				toLink = t.chatService.GetUserLink(creditTransfer.To.GetUserKey())
			}

			if creditTransfer.From.IsForGroup() {
				fromLink = t.chatService.GetGroupLink(creditTransfer.From.GetGroupKey())
			} else if creditTransfer.From.IsForUser() {
				fromLink = t.chatService.GetUserLink(creditTransfer.From.GetUserKey())
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
