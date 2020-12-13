package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
)

func (g GroupService) CancelOrDeclineInvitation(ctx context.Context, request *group.CancelOrDeclineInvitationRequest) error {

	ctx, _ = GetCtx(ctx, "GroupService", "CancelOrDeclineInvitation")

	membershipKey := request.MembershipKey

	userSession, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	membership, err := g.groupStore.GetMembership(ctx, membershipKey)
	if err != nil {
		return err
	}

	wasMember := membership.UserConfirmed && membership.GroupConfirmed

	if membership.GetUserKey() == userSession.GetUserKey() {
		// user is declining invitaiton from group

		// todo: check if last membership

		err = g.groupStore.DeleteMembership(ctx, membershipKey)
		if err != nil {
			return err
		}

	} else {
		// group is declining invitation from user

		adminMembershipKey := model.NewMembershipKey(membershipKey.GroupKey, userSession.GetUserKey())
		adminMembership, err := g.groupStore.GetMembership(ctx, adminMembershipKey)
		if err != nil {
			return err
		}

		if !adminMembership.IsAdmin {
			err := fmt.Errorf("cannot decline invitation if not admin")
			return err
		}

		err = g.groupStore.DeleteMembership(ctx, membershipKey)
		if err != nil {
			return err
		}
	}

	if wasMember {
		usernameLeavingGroup, err := g.authStore.GetUsername(membershipKey.UserKey)
		if err != nil {
			return err
		}

		amqpChannel, err := g.amqpClient.GetChannel()
		if err != nil {
			return err
		}

		channelSubscriptionKey := model.NewChannelSubscriptionKey(membershipKey.GroupKey.GetChannelKey(), membershipKey.UserKey)
		err = g.chatService.UnsubscribeFromChannel(ctx, channelSubscriptionKey)
		if err != nil {
			return err
		}

		channelKey := request.MembershipKey.GroupKey.GetChannelKey()
		err = amqpChannel.ExchangeUnbind(ctx, membershipKey.UserKey.GetExchangeName(), "", amqp.WebsocketMessagesExchange, false, map[string]interface{}{
			"event_type": "chat.message",
			"channel_id": channelKey.String(),
			"x-match":    "all",
		})
		if err != nil {
			return err
		}

		grp, err := g.groupStore.GetGroup(ctx, request.MembershipKey.GroupKey)
		if err != nil {
			return err
		}

		text := fmt.Sprintf("%s has left #%s", usernameLeavingGroup, grp.Name)
		message := chat.NewContextBlock([]chat.BlockElement{
			chat.NewMarkdownObject(text)},
			nil,
		)

		_, err = g.chatService.SendGroupMessage(ctx, chat.NewSendGroupMessage(request.MembershipKey.GroupKey, membershipKey.UserKey, usernameLeavingGroup, text, []chat.Block{*message}, []chat.Attachment{}, nil))
		if err != nil {
			return err
		}
	}

	return nil

}
