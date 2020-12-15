package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	chatmodel "github.com/commonpool/backend/pkg/chat/model"
	group "github.com/commonpool/backend/pkg/group"
)

func (g GroupService) CreateGroup(ctx context.Context, request *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {

	userSession, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	grp, membership, err := g.groupStore.CreateGroupAndMembership(ctx, request.GroupKey, userSession.GetUserKey(), request.Name, request.Description)
	if err != nil {
		return nil, err
	}

	channelKey := chatmodel.GetChannelKeyForGroup(request.GroupKey)
	channel, err := g.chatService.CreateChannel(ctx, channelKey, chatmodel.GroupChannel)
	if err != nil {
		return nil, err
	}

	channelSubscriptionKey := chatmodel.NewChannelSubscriptionKey(channel.GetKey(), userSession.GetUserKey())
	channelSubscription, err := g.chatService.SubscribeToChannel(ctx, channelSubscriptionKey, grp.Name)
	if err != nil {
		return nil, err
	}

	_, err = g.chatService.SendGroupMessage(ctx, chat.NewSendGroupMessage(grp.GetKey(), userSession.GetUserKey(), "Commonpool", "Bienvenue!", []chatmodel.Block{}, []chatmodel.Attachment{}, nil))
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
