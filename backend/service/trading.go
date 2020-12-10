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
	"github.com/commonpool/backend/transaction"
	"go.uber.org/zap"
	ctx "golang.org/x/net/context"
	"time"
)

type TradingService struct {
	tradingStore       trading.Store
	transactionService transaction.Service
	groupService       group.Service
	rs                 res.Store
	us                 auth.Store
	chatService        chat.Service
}

var _ trading.Service = &TradingService{}

func NewTradingService(
	tradingStore trading.Store,
	resourceStore res.Store,
	authStore auth.Store,
	chatService chat.Service,
	groupService group.Service,
	transactionService transaction.Service) *TradingService {
	return &TradingService{
		tradingStore:       tradingStore,
		rs:                 resourceStore,
		us:                 authStore,
		chatService:        chatService,
		groupService:       groupService,
		transactionService: transactionService,
	}
}

func (t TradingService) checkOfferCompleted(ctx context.Context, groupKey model.GroupKey, offerKey model.OfferKey, offerItems *trading.OfferItems, userConfirmingItem model.UserReference, usersInOffer *auth.Users) error {

	ctx, l := GetCtx(ctx, "TradingService", "checkOfferCompleted")

	if offerItems.AllPartiesAccepted() && offerItems.AllUserActionsCompleted() {
		for _, offerItem := range offerItems.Items {
			if offerItem.IsCreditTransfer() {
				creditTransfer := offerItem.(*trading.CreditTransferItem)
				_, err := t.transactionService.TimeCreditsExchanged(groupKey, creditTransfer.From, creditTransfer.To, creditTransfer.Amount)
				if err != nil {
					return err
				}
			}
			if offerItem.IsServiceProviding() {
				serviceProvision := offerItem.(*trading.ProvideServiceItem)
				_, err := t.transactionService.ServiceWasProvided(groupKey, serviceProvision.ResourceKey, serviceProvision.Duration)
				if err != nil {
					return err
				}
			}
			if offerItem.IsBorrowingResource() {
				borrowResource := offerItem.(*trading.BorrowResourceItem)
				_, err := t.transactionService.ResourceWasBorrowed(groupKey, borrowResource.ResourceKey, borrowResource.To, borrowResource.Duration)
				if err != nil {
					return err
				}
			}
			if offerItem.IsResourceTransfer() {
				transfer := offerItem.(*trading.ResourceTransferItem)
				_, err := t.transactionService.ResourceWasTaken(groupKey, transfer.ResourceKey, transfer.To)
				if err != nil {
					return err
				}
			}
		}

		err := t.tradingStore.UpdateOfferStatus(offerKey, trading.CompletedOffer)
		if err != nil {
			l.Error("could not mark offer as completed", zap.Error(err))
			return err
		}

		blocks, mainText, err := t.buildOfferCompletedMessage(ctx, offerItems, usersInOffer)
		if err != nil {
			l.Debug("could not build offer completion message", zap.Error(err))
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

			creditTransfer := offerItem.(*trading.CreditTransferItem)

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

func (t TradingService) FindTargetsForOfferItem(
	ctx ctx.Context,
	groupKey model.GroupKey,
	itemType trading.OfferItemType,
	from *model.Target,
	to *model.Target) (*model.Targets, error) {

	membershipStatus := group.ApprovedMembershipStatus
	membershipsForGroup, err := t.groupService.GetGroupMemberships(ctx, &group.GetMembershipsForGroupRequest{
		GroupKey:         groupKey,
		MembershipStatus: &membershipStatus,
	})
	if err != nil {
		return nil, err
	}

	group, err := t.groupService.GetGroup(ctx, &group.GetGroupRequest{
		Key: groupKey,
	})
	if err != nil {
		return nil, err
	}

	var targets []*model.Target

	groupTarget := model.NewGroupTarget(group.Group.Key)

	if to == nil || !to.Equals(groupTarget) {
		targets = append(targets, groupTarget)
	}

	for _, membership := range membershipsForGroup.Memberships.Items {
		userTarget := model.NewUserTarget(membership.GetUserKey())
		if to == nil || !to.Equals(userTarget) {
			targets = append(targets, userTarget)
		}
	}

	return model.NewTargets(targets), nil
}
