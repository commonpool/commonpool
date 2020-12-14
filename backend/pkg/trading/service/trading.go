package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/chat"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/resource"
	trading2 "github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/service"
	"github.com/commonpool/backend/transaction"
	"go.uber.org/zap"
	ctx "golang.org/x/net/context"
	"time"
)

type TradingService struct {
	tradingStore       trading2.Store
	transactionService transaction.Service
	groupService       group2.Service
	rs                 resource.Store
	us                 auth.Store
	chatService        chat.Service
}

var _ trading2.Service = &TradingService{}

func NewTradingService(
	tradingStore trading2.Store,
	resourceStore resource.Store,
	authStore auth.Store,
	chatService chat.Service,
	groupService group2.Service,
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

func (t TradingService) checkOfferCompleted(ctx context.Context, groupKey model.GroupKey, offerKey model.OfferKey, offerItems *trading2.OfferItems, userConfirmingItem model.UserReference, usersInOffer *auth.Users) error {

	ctx, l := service.GetCtx(ctx, "TradingService", "checkOfferCompleted")

	if offerItems.AllPartiesAccepted() && offerItems.AllUserActionsCompleted() {
		for _, offerItem := range offerItems.Items {
			if offerItem.IsCreditTransfer() {
				creditTransfer := offerItem.(*trading2.CreditTransferItem)
				_, err := t.transactionService.TimeCreditsExchanged(groupKey, creditTransfer.From, creditTransfer.To, creditTransfer.Amount)
				if err != nil {
					return err
				}
			}
			if offerItem.IsServiceProviding() {
				serviceProvision := offerItem.(*trading2.ProvideServiceItem)
				_, err := t.transactionService.ServiceWasProvided(groupKey, serviceProvision.ResourceKey, serviceProvision.Duration)
				if err != nil {
					return err
				}
			}
			if offerItem.IsBorrowingResource() {
				borrowResource := offerItem.(*trading2.BorrowResourceItem)
				_, err := t.transactionService.ResourceWasBorrowed(groupKey, borrowResource.ResourceKey, borrowResource.To, borrowResource.Duration)
				if err != nil {
					return err
				}
			}
			if offerItem.IsResourceTransfer() {
				transfer := offerItem.(*trading2.ResourceTransferItem)
				_, err := t.transactionService.ResourceWasTaken(groupKey, transfer.ResourceKey, transfer.To)
				if err != nil {
					return err
				}
			}
		}

		err := t.tradingStore.UpdateOfferStatus(offerKey, trading2.CompletedOffer)
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

func (t TradingService) buildOfferCompletedMessage(ctx context.Context, items *trading2.OfferItems, users *auth.Users) ([]chat.Block, string, error) {

	ctx, _ = service.GetCtx(ctx, "TradingService", "buildOfferCompletedMessage")

	var blocks []chat.Block

	mainText := ":champagne: Alright! everybody confirmed having received and given their stuff."
	blocks = append(blocks, *chat.NewHeaderBlock(
		chat.NewMarkdownObject(mainText),
		nil,
	))

	for _, offerItem := range items.Items {

		if offerItem.IsCreditTransfer() {

			creditTransfer := offerItem.(*trading2.CreditTransferItem)

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
	itemType trading2.OfferItemType,
	from *model.Target,
	to *model.Target) (*model.Targets, error) {

	membershipStatus := group2.ApprovedMembershipStatus
	membershipsForGroup, err := t.groupService.GetGroupMemberships(ctx, &group2.GetMembershipsForGroupRequest{
		GroupKey:         groupKey,
		MembershipStatus: &membershipStatus,
	})
	if err != nil {
		return nil, err
	}

	group, err := t.groupService.GetGroup(ctx, &group2.GetGroupRequest{
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
