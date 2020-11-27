package handler

import (
	"context"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// CreateGroup godoc
// @Summary Creates a group
// @Description Creates a group and sets the authenticated user as the owner
// @ID createGroup
// @Tags groups
// @Param group body web.CreateGroupRequest true "Group to create"
// @Accept json
// @Produce json
// @Success 200 {object} web.CreateGroupResponse
// @Failure 400 {object} utils.Error
// @Router /groups [post]
func (h *Handler) CreateGroup(c echo.Context) error {

	ctx, l := GetEchoContext(c, "CreateGroup")

	l.Debug("creating group")

	req := web.CreateGroupRequest{}
	if err := c.Bind(&req); err != nil {
		l.Error("could not bind CreateGroupRequest", zap.Error(err))
		return NewErrResponse(c, err)
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		response := errors.ErrValidation("name is required")
		return NewErrResponse(c, &response)
	}

	var groupKey = model.NewGroupKey(uuid.NewV4())

	var createGroupRequest = group.NewCreateGroupRequest(groupKey, req.Name, req.Description)
	createGroupResponse, err := h.groupService.CreateGroup(ctx, createGroupRequest)
	if err != nil {
		l.Error("could not create group", zap.Error(err))
		return err
	}

	var response = web.NewCreateGroupResponse(createGroupResponse.Group)
	return c.JSON(http.StatusCreated, response)

}

// GetGroup godoc
// @Summary Gets a group
// @Description Gets a group
// @ID getGroup
// @Tags groups
// @Param id path string true "ID of the group" (format:uuid)
// @Accept json
// @Produce json
// @Success 200 {object} web.GetGroupResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:id [get]
func (h *Handler) GetGroup(c echo.Context) error {

	ctx, l := GetEchoContext(c, "GetGroup")

	l.Debug("getting group")

	groupKey, err := model.ParseGroupKey(c.Param("id"))
	if err != nil {
		l.Error("could not parse group key", zap.Error(err))
		return NewErrResponse(c, err)
	}

	getGroupResponse, err := h.groupService.GetGroup(ctx, group.NewGetGroupRequest(groupKey))
	if err != nil {
		l.Error("could not get group", zap.Error(err))
		return err
	}

	var response = web.NewGetGroupResponse(getGroupResponse.Group)
	return c.JSON(http.StatusOK, response)

}

// GetLoggedInUserMemberships godoc
// @Summary Gets currently logged in user memberships
// @Description Gets the memberships for the currently logged in user
// @ID getLoggedInUserMemberships
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.GetUserMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /my/memberships [get]
func (h *Handler) GetLoggedInUserMemberships(c echo.Context) error {

	ctx, l := GetEchoContext(c, "GetLoggedInUserMemberships")

	l.Debug("getting logged in user memberships")

	authUserKey := h.authorization.GetAuthUserKey(c)

	userMembershipsResponse, err := h.groupService.GetUserMemberships(ctx, group.NewGetMembershipsForUserRequest(authUserKey, group.AnyMembershipStatus()))
	if err != nil {
		l.Error("could not get logged in user memberships", zap.Error(err))
		return NewErrResponse(c, err)
	}

	memberships := userMembershipsResponse.Memberships

	groupNames, err := h.getGroupNamesForMemberships(ctx, memberships)
	if err != nil {
		return NewErrResponse(c, err)
	}

	userNames, err := h.getUserNamesForMemberships(ctx, memberships)
	if err != nil {
		return NewErrResponse(c, err)
	}

	response := web.NewGetUserMembershipsResponse(userMembershipsResponse.Memberships, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}

// GetUserMemberships godoc
// @Summary Gets memberships for a given user
// @Description Gets the memberships for a given user
// @ID getUserMemberships
// @Param id path string true "ID of the user" (format:uuid)
// @Param status status MembershipStatus true "status of the membership"
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.GetUserMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /users/:id/memberships [get]
func (h *Handler) GetUserMemberships(c echo.Context) error {

	ctx, l := GetEchoContext(c, "GetUserMemberships")

	l.Debug("getting user memberships")

	var membershipStatus = group.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := group.ParseMembershipStatus(statusStr)
		if err != nil {

			l.Error("could not parse desired membership 'status' query param",
				zap.Error(err),
				zap.String("status", statusStr))

			return NewErrResponse(c, err)

		}
		membershipStatus = &ms
	}

	userKey := model.NewUserKey(c.Param("id"))
	getMembershipsForUserRequest := group.NewGetMembershipsForUserRequest(userKey, membershipStatus)
	getMembershipsResponse, err := h.groupService.GetUserMemberships(ctx, getMembershipsForUserRequest)
	if err != nil {
		l.Error("could not get user memberships", zap.Error(err))
		return err
	}

	groupNames, err := h.getGroupNamesForMemberships(ctx, getMembershipsResponse.Memberships)
	if err != nil {
		l.Error("could not get group names for memberships", zap.Error(err))
		return NewErrResponse(c, err)
	}

	userNames, err := h.getUserNamesForMemberships(ctx, getMembershipsResponse.Memberships)
	if err != nil {
		l.Error("could not get user names for memberships", zap.Error(err))
		return NewErrResponse(c, err)
	}

	response := web.NewGetUserMembershipsResponse(getMembershipsResponse.Memberships, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}

// GetUserMemberships godoc
// @Summary Gets the membership for a given user and group
// @Description Gets the membership for a given user and group
// @ID getMembership
// @Param groupId path string true "ID of the group" (format:uuid)
// @Param userId path string true "ID of the user"
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.GetMembershipResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:groupId/memberships/:userId [get]
func (h *Handler) GetMembership(c echo.Context) error {

	ctx, l := GetEchoContext(c, "GetMembership")

	l.Debug("getting memberships")

	userKey := model.NewUserKey(c.Param("userId"))

	groupKey, err := model.ParseGroupKey(c.Param("groupId"))
	if err != nil {
		l.Error("could not parse group key", zap.Error(err))
		return NewErrResponse(c, err)
	}

	membershipKey := model.NewMembershipKey(groupKey, userKey)
	getMembershipsRequest := group.NewGetMembershipRequest(membershipKey)
	getMemberships, err := h.groupService.GetMembership(ctx, getMembershipsRequest)
	if err != nil {
		l.Error("could not get membership", zap.Error(err))
		return NewErrResponse(c, err)
	}

	var memberships = group.NewMemberships([]group.Membership{*getMemberships.Membership})

	groupNames, err := h.getGroupNamesForMemberships(ctx, memberships)
	if err != nil {
		l.Error("could not get group names for memberships", zap.Error(err))
		return NewErrResponse(c, err)
	}

	userNames, err := h.getUserNamesForMemberships(ctx, memberships)
	if err != nil {
		l.Error("could not get user names for memberships", zap.Error(err))
		return NewErrResponse(c, err)
	}

	response := web.NewGetMembershipResponse(getMemberships.Membership, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}

// GetGroup godoc
// @Summary Gets a group memberships
// @Description Gets the members of a group
// @ID getGroupMemberships
// @Tags groups
// @Param id path string true "ID of the group" (format:uuid)
// @Param status status MembershipStatus true "status of the membership"
// @Accept json
// @Produce json
// @Success 200 {object} web.GetGroupMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:id/memberships [get]
func (h *Handler) GetGroupMemberships(c echo.Context) error {

	ctx, l := GetEchoContext(c, "GetGroupMemberships")

	var membershipStatus = group.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := group.ParseMembershipStatus(statusStr)
		if err != nil {
			l.Error("could not parse desired membership 'status' query param", zap.String("status", statusStr))
			return NewErrResponse(c, err)
		}
		membershipStatus = &ms
	}

	groupKey, err := model.ParseGroupKey(c.Param("id"))
	if err != nil {
		l.Error("could not parse group key", zap.Error(err))
		return NewErrResponse(c, err)
	}

	_, err = h.groupService.GetGroup(ctx, group.NewGetGroupRequest(groupKey))
	if err != nil {
		l.Error("could not get group", zap.Error(err))
		return NewErrResponse(c, err)
	}

	getGroupMemberships, err := h.groupService.GetGroupsMemberships(ctx, group.NewGetMembershipsForGroupRequest(groupKey, membershipStatus))
	if err != nil {
		l.Error("could not get group memberships", zap.Error(err))
		return NewErrResponse(c, err)
	}

	userNames, err := h.getUserNamesForMemberships(ctx, getGroupMemberships.Memberships)
	if err != nil {
		l.Error("could not get user names for memberships", zap.Error(err))
		return NewErrResponse(c, err)
	}

	groupNames, err := h.getGroupNamesForMemberships(ctx, getGroupMemberships.Memberships)
	if err != nil {
		l.Error("could not get group names for memberships", zap.Error(err))
		return NewErrResponse(c, err)
	}

	response := web.NewGetUserMembershipsResponse(getGroupMemberships.Memberships, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}

func (h *Handler) getUserNamesForMemberships(ctx context.Context, memberships *group.Memberships) (auth.UserNames, error) {

	ctx, l := GetCtx(ctx, "getUserNamesForMemberships")

	l.Debug("getting user names for memberships")

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

func (h *Handler) getGroupNamesForMemberships(ctx context.Context, memberships *group.Memberships) (group.GroupNames, error) {

	ctx, l := GetCtx(ctx, "getGroupNamesForMemberships")

	l.Debug("getting group names for memberships")

	var groupNames = group.GroupNames{}
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

// GetGroup godoc
// @Summary User picker for group invite
// @Description Finds users to invite on a group
// @ID inviteMemberPicker
// @Tags groups
// @Param id path string true "ID of the group" (format:uuid)
// @Accept json
// @Produce json
// @Success 200 {object} web.GetGroupMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:id/invite-member-picker [get]
func (h *Handler) GetUsersForGroupInvitePicker(c echo.Context) error {
	skip, err := utils.ParseSkip(c)
	if err != nil {
		return NewErrResponse(c, err)
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return NewErrResponse(c, err)
	}

	qry := c.QueryParam("query")

	groupKey, err := model.ParseGroupKey(c.Param("id"))
	if err != nil {
		return NewErrResponse(c, err)
	}

	userQuery := auth.UserQuery{
		Query:      qry,
		Skip:       skip,
		Take:       take,
		NotInGroup: &groupKey,
	}

	users, err := h.authStore.Find(userQuery)
	if err != nil {
		return NewErrResponse(c, err)
	}

	responseItems := make([]web.UserInfoResponse, len(users))
	for i, user := range users {
		responseItems[i] = web.UserInfoResponse{
			Id:       user.ID,
			Username: user.Username,
		}
	}

	response := web.GetUsersForGroupInvitePickerResponse{
		Users: responseItems,
		Take:  take,
		Skip:  skip,
	}

	return c.JSON(http.StatusOK, response)

}

// GetGroup godoc
// @Summary Invite a user to a group
// @Description Invite a user to a group
// @ID inviteUser
// @Tags groups
// @Param id path string true "ID of the group" (format:uuid)
// @Param invite body web.InviteUserRequest true "User to invite"
// @Accept json
// @Produce json
// @Success 200 {object} web.InviteUserResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:id/invite [get]
func (h *Handler) InviteUser(c echo.Context) error {

	ctx, l := GetEchoContext(c, "InviteUser")

	l.Debug("inviting user to group")

	req := web.InviteUserRequest{}
	if err := c.Bind(&req); err != nil {
		l.Error("could not bind web.InviteUserRequest", zap.Error(err))
		return NewErrResponse(c, err)
	}

	groupKey, err := model.ParseGroupKey(c.Param("id"))
	if err != nil {
		l.Error("could not parse group key", zap.String("group_id", c.Param("id")))
		return NewErrResponse(c, err)
	}

	invitedUserKey := model.NewUserKey(req.UserID)
	newMembershipKey := model.NewMembershipKey(groupKey, invitedUserKey)

	membership, err := h.groupService.SendGroupInvitation(ctx, group.NewInviteRequest(newMembershipKey))
	if err != nil {
		l.Error("could not send invitation", zap.Error(err))
		return err
	}

	getGroup, err := h.groupService.GetGroup(ctx, group.NewGetGroupRequest(groupKey))
	if err != nil {
		l.Error("could not get group", zap.Error(err))
		return NewErrResponse(c, err)
	}
	groupNames := group.GroupNames{
		groupKey: getGroup.Group.Name,
	}

	username, err := h.authStore.GetUsername(membership.Membership.GetUserKey())
	if err != nil {
		l.Error("could not get username", zap.Error(err))
		return NewErrResponse(c, err)
	}
	userNames := auth.UserNames{
		newMembershipKey.UserKey: username,
	}

	// Respond to query
	response := web.NewInviteUserResponse(membership.Membership, groupNames, userNames)
	return c.JSON(http.StatusAccepted, response)

}

// AcceptInvitation godoc
// @Summary Accept a group invitation
// @Description Accept a group invitation
// @ID acceptInvitation
// @Tags groups
// @Param groupId path string true "ID of the group" (format:uuid)
// @Param userId path string true "ID of the user"
// @Accept json
// @Produce json
// @Success 200 {object} web.AcceptInvitationResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:groupId/memberships/:userId/accept [post]
func (h *Handler) AcceptInvitation(c echo.Context) error {

	ctx, l := GetEchoContext(c, "AcceptInvitation")

	groupKey, err := model.ParseGroupKey(c.Param("groupId"))
	if err != nil {
		l.Error("could not parse group key", zap.Error(err))
		return NewErrResponse(c, err)
	}
	userKey := model.NewUserKey(c.Param("userId"))

	l = l.With(zap.Object("user", userKey), zap.Object("group", groupKey))

	membershipKey := model.NewMembershipKey(groupKey, userKey)
	acceptInvitationResponse, err := h.groupService.AcceptInvitation(ctx, group.NewAcceptInvitationRequest(membershipKey))
	if err != nil {
		l.Error("could not accept invitation", zap.Error(err))
		return err
	}

	memberships := group.NewMemberships([]group.Membership{*acceptInvitationResponse.Membership})

	userNames, err := h.getUserNamesForMemberships(ctx, memberships)
	if err != nil {
		l.Error("could not get user names for memberships", zap.Error(err))
		return NewErrResponse(c, err)
	}

	groupNames, err := h.getGroupNamesForMemberships(ctx, memberships)
	if err != nil {
		l.Error("could not get group names for memberships", zap.Error(err))
		return NewErrResponse(c, err)
	}

	response := web.NewAcceptInvitationResponse(acceptInvitationResponse.Membership, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}

// DeclineInvitation godoc
// @Summary declines a group invitation
// @Description declines a group invitation
// @ID declineInvitation
// @Tags groups
// @Param groupId path string true "ID of the group" (format:uuid)
// @Param userId path string true "ID of the user"
// @Accept json
// @Produce json
// @Success 200 {object} web.DeclineInvitationResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:groupId/memberships/:userId/decline [post]
func (h *Handler) DeclineInvitation(c echo.Context) error {

	ctx, l := GetEchoContext(c, "DeclineInvitation")

	l.Debug("declining invitation")

	groupKey, err := model.ParseGroupKey(c.Param("groupId"))
	if err != nil {
		return NewErrResponse(c, err)
	}

	userKey := model.NewUserKey(c.Param("userId"))
	membershipKey := model.NewMembershipKey(groupKey, userKey)

	err = h.groupService.DeclineInvitation(ctx, group.NewDelineInvitationRequest(membershipKey))
	if err != nil {
		l.Error("could not decline invitation", zap.Error(err))
		return err
	}

	return c.NoContent(http.StatusAccepted)

}

// LeaveGroup godoc
// @Summary leave group
// @Description declines a group invitation
// @ID leaveGroup
// @Tags groups
// @Param groupId path string true "ID of the group" (format:uuid)
// @Param userId path string true "ID of the user"
// @Accept json
// @Produce json
// @Success 200 {object} web.LeaveGroupResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:groupId/memberships/:userId [delete]
func (h *Handler) LeaveGroup(c echo.Context) error {

	ctx, l := GetEchoContext(c, "LeaveGroup")

	l.Debug("leaving group")

	groupKey, err := model.ParseGroupKey(c.Param("groupId"))
	if err != nil {
		return NewErrResponse(c, err)
	}

	userKey := model.NewUserKey(c.Param("userId"))
	membershipKey := model.NewMembershipKey(groupKey, userKey)

	err = h.groupService.LeaveGroup(ctx, group.NewLeaveGroupRequest(membershipKey))
	if err != nil {
		l.Error("could not decline invitation", zap.Error(err))
		return err
	}

	return c.NoContent(http.StatusAccepted)

}
