package service

import (
	"context"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"go.uber.org/zap"
)

func (g GroupService) CreateGroup(ctx context.Context, request *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {

	ctx, l := GetCtx(ctx, "GroupService", "CreateGroup")

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("could not get user session", zap.Error(err))
		return nil, err
	}

	grp, err := g.groupStore.CreateGroup(ctx, request.GroupKey, userSession.GetUserKey(), request.Name, request.Description)
	if err != nil {
		l.Error("could not save group to store", zap.Error(err))
		return nil, err
	}

	membershipKey := model.NewMembershipKey(grp.GetKey(), userSession.GetUserKey())
	membership, err := g.groupStore.CreateMembership(ctx, membershipKey, true, true, true, false, true, true)
	if err != nil {
		l.Error("could not create membership", zap.Error(err))
		return nil, err
	}

	channel, err := g.chatService.CreateChannel(ctx, request.GroupKey.GetChannelKey(), chat.GroupChannel)
	if err != nil {
		l.Error("could not create channel", zap.Error(err))
		return nil, err
	}

	channelSubscriptionKey := model.NewChannelSubscriptionKey(channel.GetKey(), userSession.GetUserKey())
	channelSubscription, err := g.chatService.SubscribeToChannel(ctx, channelSubscriptionKey, grp.Name)
	if err != nil {
		l.Error("could not create channel subscription", zap.Error(err))
		return nil, err
	}

	_, err = g.chatService.SendGroupMessage(ctx, chat.NewSendGroupMessage(grp.GetKey(), userSession.GetUserKey(), "Commonpool", "Bienvenue!", []chat.Block{}, []chat.Attachment{}, nil))
	if err != nil {
		l.Error("could not send group message", zap.Error(err))
		return nil, err
	}

	return &group.CreateGroupResponse{
		ChannelKey:      channel.GetKey(),
		SubscriptionKey: channelSubscription.GetKey(),
		Group:           grp,
		Membership:      membership,
	}, nil

}
