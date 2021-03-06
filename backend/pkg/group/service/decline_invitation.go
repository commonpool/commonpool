package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	group2 "github.com/commonpool/backend/pkg/group"
)

func (g GroupService) CancelOrDeclineInvitation(ctx context.Context, request *group2.CancelOrDeclineInvitationRequest) error {

	userSession, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	group, err := g.groupRepo.Load(ctx, request.MembershipKey.GroupKey)
	if err != nil {
		return err
	}

	if err := group.CancelMembership(userSession.GetUserKey(), request.MembershipKey.UserKey); err != nil {
		return err
	}

	if err := g.groupRepo.Save(ctx, group); err != nil {
		return err
	}

	return nil

	//
	//
	//
	// membershipKey := request.MembershipKey
	//
	// userSession, err := oidc.GetLoggedInUser(ctx)
	// if err != nil {
	// 	return err
	// }
	//
	// membership, err := g.groupStore.GetMembership(ctx, membershipKey)
	// if err != nil {
	// 	return err
	// }
	//
	// wasMember := membership.HasBothPartiesConfirmed()
	//
	// if membership.GetUserKey() == userSession.GetUserKey() {
	// 	// user is declining invitaiton from group
	//
	// 	// todo: check if last membership
	//
	// 	err = g.groupStore.DeleteMembership(ctx, membershipKey)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// } else {
	// 	// group is declining invitation from user
	//
	// 	adminMembershipKey := keys.NewMembershipKey(membershipKey.GroupKey, userSession.GetUserKey())
	// 	adminMembership, err := g.groupStore.GetMembership(ctx, adminMembershipKey)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	if !adminMembership.IsAdmin() {
	// 		err := fmt.Errorf("cannot decline invitation if not admin")
	// 		return err
	// 	}
	//
	// 	err = g.groupStore.DeleteMembership(ctx, membershipKey)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	//
	// if wasMember {
	// 	usernameLeavingGroup, err := g.authStore.GetUsername(membershipKey.UserKey)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	amqpChannel, err := g.amqpClient.GetChannel()
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	channelKey := keys.GetChannelKeyForGroup(membershipKey.GroupKey)
	//
	// 	channelSubscriptionKey := keys.NewChannelSubscriptionKey(channelKey, membershipKey.UserKey)
	// 	err = g.chatService.UnsubscribeFromChannel(ctx, channelSubscriptionKey)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	err = amqpChannel.ExchangeUnbind(ctx, membershipKey.UserKey.GetExchangeName(), "", mq.WebsocketMessagesExchange, false, map[string]interface{}{
	// 		"event_type": "chat.message",
	// 		"channel_id": channelKey.String(),
	// 		"x-match":    "all",
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	grp, err := g.groupStore.GetGroup(ctx, request.MembershipKey.GroupKey)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	text := fmt.Sprintf("%s has left #%s", usernameLeavingGroup, grp.Name)
	// 	message := chat.NewContextBlock([]chat.BlockElement{
	// 		chat.NewMarkdownObject(text)},
	// 		nil,
	// 	)
	//
	// 	_, err = g.chatService.SendGroupMessage(ctx, service.NewSendGroupMessage(request.MembershipKey.GroupKey, membershipKey.UserKey, usernameLeavingGroup, text, []chat.Block{*message}, []chat.Attachment{}, nil))
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	//
	// return nil

}
