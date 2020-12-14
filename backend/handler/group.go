package handler

import (
	"context"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/pkg/handler"
	"go.uber.org/zap"
)

func (h *Handler) getUserNamesForMemberships(ctx context.Context, memberships *group.Memberships) (auth.UserNames, error) {

	ctx, l := handler.GetCtx(ctx, "getUserNamesForMemberships")

	var userNames = auth.UserNames{}
	for _, membership := range memberships.Items {
		userKey := membership.GetUserKey()
		_, ok := userNames[userKey]
		if !ok {
			username, err := h.authStore.GetUsername(userKey)
			if err != nil {
				l.Error("could not get username", zap.String("user_id", userKey.String()))
				return userNames, err
			}
			userNames[userKey] = username
		}
	}
	return userNames, nil
}

func (h *Handler) getGroupNamesForMemberships(ctx context.Context, memberships *group.Memberships) (group.Names, error) {

	ctx, l := handler.GetCtx(ctx, "getGroupNamesForMemberships")

	var groupNames = group.Names{}
	for _, membership := range memberships.Items {
		groupKey := membership.GetGroupKey()
		_, ok := groupNames[groupKey]
		if !ok {
			getGroup, err := h.groupService.GetGroup(ctx, group.NewGetGroupRequest(groupKey))
			if err != nil {
				l.Error("could not get group", zap.Error(err))
				return groupNames, err
			}
			groupNames[groupKey] = getGroup.Group.Name
		}
	}
	return groupNames, nil
}
