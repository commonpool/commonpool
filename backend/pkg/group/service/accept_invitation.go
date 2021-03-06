package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth"
	group2 "github.com/commonpool/backend/pkg/group"
)

func (g GroupService) CreateOrAcceptInvitation(ctx context.Context, request *group2.CreateOrAcceptInvitationRequest) error {

	userSession, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	group, err := g.groupRepo.Load(ctx, request.MembershipKey.GroupKey)
	if err != nil {
		return err
	}

	if err := group.JoinGroup(userSession.GetUserKey(), request.MembershipKey.UserKey); err != nil {
		return err
	}

	if err := g.groupRepo.Save(ctx, group); err != nil {
		return err
	}

	return nil

	// if acceptedMembership.GroupConfirmed && acceptedMembership.UserConfirmed {
	// 	// add user to group channel
	//
	// 	grp, err := g.groupStore.GetGroup(ctx, request.MembershipKey.GroupKey)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	usernameJoiningGroup, err := g.authStore.GetUsername(membershipKey.UserKey)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	channelKey := keys.GetChannelKeyForGroup(request.MembershipKey.GroupKey)
	//
	// 	channelSubscriptionKey := keys.NewChannelSubscriptionKey(channelKey, acceptedMembership.GetUserKey())
	// 	_, err = g.chatService.SubscribeToChannel(ctx, channelSubscriptionKey, grp.Name)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	amqpChannel, err := g.amqpClient.GetChannel()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	err = amqpChannel.ExchangeBind(ctx, membershipKey.UserKey.GetExchangeName(), "", mq.WebsocketMessagesExchange, false, map[string]interface{}{
	// 		"event_type": "chat.message",
	// 		"channel_id": channelKey.String(),
	// 		"x-match":    "all",
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	text := fmt.Sprintf("%s has joined #%s", usernameJoiningGroup, grp.Name)
	// 	message := chat.NewContextBlock([]chat.BlockElement{
	// 		chat.NewMarkdownObject(text)},
	// 		nil,
	// 	)
	//
	// 	_, err = g.chatService.SendGroupMessage(ctx, chat.NewSendGroupMessage(request.MembershipKey.GroupKey, membershipKey.UserKey, usernameJoiningGroup, text, []chat.Block{*message}, []chat.Attachment{}, nil))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// }
	//
	// return &group2.CreateOrAcceptInvitationResponse{
	// 	Membership: acceptedMembership,
	// }, nil

}
