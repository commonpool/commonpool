package service

import (
	"context"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	group2 "github.com/commonpool/backend/pkg/group"
)

func (g GroupService) CreateGroup(ctx context.Context, request *group2.CreateGroupRequest) (*group2.CreateGroupResponse, error) {

	userSession, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	grp, membership, err := g.groupStore.CreateGroupAndMembership(ctx, request.GroupKey, userSession.GetUserKey(), request.Name, request.Description)
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

	return &group2.CreateGroupResponse{
		ChannelKey:      channel.GetKey(),
		SubscriptionKey: channelSubscription.GetKey(),
		Group:           grp,
		Membership:      membership,
	}, nil

}
