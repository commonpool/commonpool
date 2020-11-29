package service

import (
	"context"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
)

func (g GroupService) CreateGroup(ctx context.Context, request *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {

	ctx, _ = GetCtx(ctx, "GroupService", "CreateGroup")

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		return nil, err
	}

	grp, err := g.groupStore.CreateGroup(ctx, request.GroupKey, userSession.GetUserKey(), request.Name, request.Description)
	if err != nil {
		return nil, err
	}

	membershipKey := model.NewMembershipKey(grp.GetKey(), userSession.GetUserKey())
	membership, err := g.groupStore.CreateMembership(ctx, membershipKey, true, true, true, false, true, true)
	if err != nil {
		return nil, err
	}

	channel, err := g.chatService.CreateChannel(ctx, request.GroupKey.GetChannelKey(), chat.GroupChannel)
	if err != nil {
		return nil, err
	}

	channelSubscriptionKey := model.NewChannelSubscriptionKey(channel.GetKey(), userSession.GetUserKey())
	channelSubscription, err := g.chatService.SubscribeToChannel(ctx, channelSubscriptionKey, grp.Name)
	if err != nil {
		return nil, err
	}

	_, err = g.chatService.SendGroupMessage(ctx, chat.NewSendGroupMessage(grp.GetKey(), userSession.GetUserKey(), "Commonpool", "Bienvenue!", []chat.Block{}, []chat.Attachment{}, nil))
	if err != nil {
		return nil, err
	}

	return &group.CreateGroupResponse{
		ChannelKey:      channel.GetKey(),
		SubscriptionKey: channelSubscription.GetKey(),
		Group:           grp,
		Membership:      membership,
	}, nil

}
