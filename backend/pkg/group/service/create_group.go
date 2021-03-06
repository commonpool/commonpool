package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/service"
	group "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/keys"
)

func (g GroupService) CreateGroup(ctx context.Context, request *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {

	userSession, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	grp, membership, err := g.groupStore.CreateGroupAndMembership(ctx, request.GroupKey, userSession.GetUserKey(), request.Name, request.Description)
	if err != nil {
		return nil, err
	}

	channelKey := keys.GetChannelKeyForGroup(request.GroupKey)
	channel, err := g.chatService.CreateChannel(ctx, channelKey, chat.GroupChannel)
	if err != nil {
		return nil, err
	}

	channelSubscriptionKey := keys.NewChannelSubscriptionKey(channel.GetKey(), userSession.GetUserKey())
	_, err = g.chatService.SubscribeToChannel(ctx, channelSubscriptionKey, grp.Name)
	if err != nil {
		return nil, err
	}

	_, err = g.chatService.SendGroupMessage(ctx, service.NewSendGroupMessage(grp.GetKey(), userSession.GetUserKey(), "Commonpool", "Bienvenue!", []chat.Block{}, []chat.Attachment{}, nil))
	if err != nil {
		return nil, err
	}

	return &group.CreateGroupResponse{
		Group:      grp,
		Membership: membership,
	}, nil

}
