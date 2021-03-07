package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	group "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	uuid "github.com/satori/go.uuid"
)

func (g GroupService) CreateGroup(ctx context.Context, request *group.CreateGroupRequest) (keys.GroupKey, error) {

	userSession, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return keys.GroupKey{}, err
	}

	groupKey := keys.NewGroupKey(uuid.NewV4())
	grp := domain.NewGroup(groupKey)
	if err := grp.CreateGroup(userSession.GetUserKey(), domain.GroupInfo{
		Name:        request.Name,
		Description: request.Description,
	}); err != nil {
		return keys.GroupKey{}, err
	}

	if err := g.groupRepo.Save(ctx, grp); err != nil {
		return keys.GroupKey{}, err
	}

	return groupKey, nil

	//
	//
	// grp, membership, err := g.groupStore.CreateGroupAndMembership(ctx, request.GroupKey, userSession.GetUserKey(), request.ResourceName, request.Description)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// channelKey := keys.GetChannelKeyForGroup(request.GroupKey)
	// channel, err := g.chatService.CreateChannel(ctx, channelKey, chat.GroupChannel)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// channelSubscriptionKey := keys.NewChannelSubscriptionKey(channel.GetKey(), userSession.GetUserKey())
	// _, err = g.chatService.SubscribeToChannel(ctx, channelSubscriptionKey, grp.ResourceName)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// _, err = g.chatService.SendGroupMessage(ctx, service.NewSendGroupMessage(grp.GetKey(), userSession.GetUserKey(), "Commonpool", "Bienvenue!", []chat.Block{}, []chat.Attachment{}, nil))
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return &g.CreateGroupResponse{
	// 	Group:      grp,
	// 	Membership: membership,
	// }, nil

}
