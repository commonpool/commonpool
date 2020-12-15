package handler

import (
	"github.com/commonpool/backend/pkg/exceptions"
	group2 "github.com/commonpool/backend/pkg/group"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/handler"
	model3 "github.com/commonpool/backend/pkg/resource/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func (h *ResourceHandler) ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c echo.Context, loggedInUserKey usermodel.UserKey, sharedWithGroups *groupmodel.GroupKeys) (error, bool) {

	ctx, l := handler.GetEchoContext(c, "ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf")

	var membershipStatus = groupmodel.ApprovedMembershipStatus

	userMemberships, err := h.groupService.GetUserMemberships(ctx, group2.NewGetMembershipsForUserRequest(loggedInUserKey, &membershipStatus))
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

func (h *ResourceHandler) parseGroupKeys(c echo.Context, sharedWith []web.InputResourceSharing) (*groupmodel.GroupKeys, error, bool) {
	sharedWithGroupKeys := make([]groupmodel.GroupKey, len(sharedWith))
	for i := range sharedWith {
		groupKeyStr := sharedWith[i].GroupID
		groupKey, err := groupmodel.ParseGroupKey(groupKeyStr)
		if err != nil {
			return nil, c.String(http.StatusBadRequest, "invalid group key : "+groupKeyStr), true
		}
		sharedWithGroupKeys[i] = groupKey
	}
	return groupmodel.NewGroupKeys(sharedWithGroupKeys), nil, false
}

func NewResourceResponse(res *model3.Resource, creatorUsername string, creatorId string, sharedWithGroups *groupmodel.Groups) web.Resource {

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
	res, ok := err.(*exceptions.ErrorResponse)
	if !ok {
		statusCode := http.StatusInternalServerError
		return c.JSON(statusCode, exceptions.NewError(err.Error(), "", statusCode))
	}
	return c.JSON(res.StatusCode, res)
}