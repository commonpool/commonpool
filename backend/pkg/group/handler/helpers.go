package handler

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/handler"
	"go.uber.org/zap"
)

func (h *Handler) getUserNamesForMemberships(ctx context.Context, memberships *domain.Memberships) (models.UserNames, error) {

	ctx, l := handler.GetCtx(ctx, "getUserNamesForMemberships")

	var userNames = models.UserNames{}
	for _, membership := range memberships.Items {
		userKey := membership.GetUserKey()
		_, ok := userNames[userKey]
		if !ok {
			username, err := h.userService.GetUsername(userKey)
			if err != nil {
				l.Error("could not get username", zap.String("user_id", userKey.String()))
				return userNames, err
			}
			userNames[userKey] = username
		}
	}
	return userNames, nil
}

func (h *Handler) getGroupNamesForMemberships(ctx context.Context, memberships *domain.Memberships) (group.Names, error) {

	ctx, l := handler.GetCtx(ctx, "getGroupNamesForMemberships")

	var groupNames = group.Names{}
	for _, membership := range memberships.Items {
		groupKey := membership.GetGroupKey()
		_, ok := groupNames[groupKey]
		if !ok {
			getGroup, err := h.groupService.GetGroup(ctx, groupKey)
			if err != nil {
				l.Error("could not get group", zap.Error(err))
				return groupNames, err
			}
			groupNames[groupKey] = getGroup.Name
		}
	}
	return groupNames, nil
}
