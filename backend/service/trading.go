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
