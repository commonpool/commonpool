package handler

import (
	. "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c echo.Context, loggedInUserKey model.UserKey, sharedWithGroups *model.GroupKeys) (error, bool) {

	ctx, l := handler.GetEchoContext(c, "ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf")

	var membershipStatus = group.ApprovedMembershipStatus

	userMemberships, err := h.groupService.GetUserMemberships(ctx, group.NewGetMembershipsForUserRequest(loggedInUserKey, &membershipStatus))
	if err != nil {
		l.Error("could not get user memberships", zap.Error(err))
		return err, true
	}

	// Checking if resource is shared with groups the user is part of
	for _, sharedWith := range sharedWithGroups.Items {
		hasMembershipInGroup := userMemberships.Memberships.ContainsMembershipForGroup(sharedWith)
		if !hasMembershipInGroup {
			return c.String(http.StatusBadRequest, "cannot share resource with a group you are not part of"), true
		}
	}
	return nil, false
}

func (h *Handler) parseGroupKeys(c echo.Context, sharedWith []web.InputResourceSharing) (*model.GroupKeys, error, bool) {
	sharedWithGroupKeys := make([]model.GroupKey, len(sharedWith))
	for i := range sharedWith {
		groupKeyStr := sharedWith[i].GroupID
		groupKey, err := group.ParseGroupKey(groupKeyStr)
		if err != nil {
			return nil, c.String(http.StatusBadRequest, "invalid group key : "+groupKeyStr), true
		}
		sharedWithGroupKeys[i] = groupKey
	}
	return model.NewGroupKeys(sharedWithGroupKeys), nil, false
}

func NewResourceResponse(res *resource.Resource, creatorUsername string, creatorId string, sharedWithGroups *group.Groups) web.Resource {

	//goland:noinspection GoPreferNilSlice
	var sharings = []web.OutputResourceSharing{}
	for _, withGroup := range sharedWithGroups.Items {
		sharings = append(sharings, web.OutputResourceSharing{
			GroupID:   withGroup.Key.String(),
			GroupName: withGroup.Name,
		})
	}

	return web.Resource{
		Id:               res.Key.String(),
		Type:             res.Type,
		SubType:          res.SubType,
		Description:      res.Description,
		Summary:          res.Summary,
		CreatedBy:        creatorUsername,
		CreatedById:      creatorId,
		CreatedAt:        res.CreatedAt,
		ValueInHoursFrom: res.ValueInHoursFrom,
		ValueInHoursTo:   res.ValueInHoursTo,
		SharedWith:       sharings,
	}
}

func NewErrResponse(c echo.Context, err error) error {
	res, ok := err.(*ErrorResponse)
	if !ok {
		statusCode := http.StatusInternalServerError
		return c.JSON(statusCode, NewError(err.Error(), "", statusCode))
	}
	return c.JSON(res.StatusCode, res)
}
