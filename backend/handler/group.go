package handler

import (
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
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

	authUserKey := h.authorization.GetAuthUserKey(c)

	req := web.CreateGroupRequest{}
	if err := c.Bind(&req); err != nil {
		return NewErrResponse(c, err)
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		response := errors.ErrValidation("name is required")
		return NewErrResponse(c, &response)
	}

	var groupKey = model.NewGroupKey(uuid.NewV4())

	var createGroupRequest = group.NewCreateGroupRequest(groupKey, authUserKey, req.Name, req.Description)
	var createGroup = h.groupStore.CreateGroup(createGroupRequest)
	if createGroup.Error != nil {
		return NewErrResponse(c, createGroup.Error)
	}

	var getGroupRequest = group.NewGetGroupRequest(groupKey)
	var getGroup = h.groupStore.GetGroup(getGroupRequest)
	if getGroup.Error != nil {
		return NewErrResponse(c, getGroup.Error)
	}

	var response = web.NewCreateGroupResponse(getGroup.Group)
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

	groupKey, err := model.ParseGroupKey(c.Param("id"))
	if err != nil {
		return NewErrResponse(c, err)
	}

	var getGroupRequest = group.NewGetGroupRequest(groupKey)
	var getGroup = h.groupStore.GetGroup(getGroupRequest)
	if getGroup.Error != nil {
		return NewErrResponse(c, getGroup.Error)
	}

	var response = web.NewGetGroupResponse(getGroup.Group)
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

	authUserKey := h.authorization.GetAuthUserKey(c)

	getMembershipsRequest := group.NewGetMembershipsForUserRequest(authUserKey, model.AnyMembershipStatus())
	getMemberships := h.groupStore.GetMembershipsForUser(getMembershipsRequest)
	if getMemberships.Error != nil {
		return NewErrResponse(c, getMemberships.Error)
	}

	memberships := getMemberships.Memberships

	groupNames, err := h.getGroupNamesForMemberships(memberships)
	if err != nil {
		return NewErrResponse(c, err)
	}

	userNames, err := h.getUserNamesForMemberships(memberships)
	if err != nil {
		return NewErrResponse(c, err)
	}

	response := web.NewGetUserMembershipsResponse(getMemberships.Memberships, groupNames, userNames)
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

	var membershipStatus = model.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := model.ParseMembershipStatus(statusStr)
		if err != nil {
			return NewErrResponse(c, err)
		}
		membershipStatus = &ms
	}

	userKey := model.NewUserKey(c.Param("id"))

	getMembershipsRequest := group.NewGetMembershipsForUserRequest(userKey, membershipStatus)
	getMemberships := h.groupStore.GetMembershipsForUser(getMembershipsRequest)
	if getMemberships.Error != nil {
		return NewErrResponse(c, getMemberships.Error)
	}

	memberships := getMemberships.Memberships

	groupNames, err := h.getGroupNamesForMemberships(memberships)
	if err != nil {
		return NewErrResponse(c, err)
	}

	userNames, err := h.getUserNamesForMemberships(memberships)
	if err != nil {
		return NewErrResponse(c, err)
	}

	response := web.NewGetUserMembershipsResponse(getMemberships.Memberships, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}

// GetGroup godoc
// @Summary Gets a group memberships
// @Description Gets the members of a group
// @ID getGroupMemberships
// @Tags groups
// @Param id path string true "ID of the group" (format:uuid)
// @Accept json
// @Produce json
// @Success 200 {object} web.GetGroupMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:id/memberships [get]
func (h *Handler) GetGroupMemberships(c echo.Context) error {

	groupKey, err := model.ParseGroupKey(c.Param("id"))
	if err != nil {
		return NewErrResponse(c, err)
	}

	var getGroupRequest = group.NewGetGroupRequest(groupKey)
	var getGroup = h.groupStore.GetGroup(getGroupRequest)
	if getGroup.Error != nil {
		return NewErrResponse(c, getGroup.Error)
	}

	getGroupMembershipsRequest := group.NewGetMembershipsForGroupRequest(groupKey)
	getGroupMemberships := h.groupStore.GetMembershipsForGroup(getGroupMembershipsRequest)
	if getGroupMemberships.Error != nil {
		return NewErrResponse(c, err)
	}

	memberships := getGroupMemberships.Memberships

	userNames, err := h.getUserNamesForMemberships(memberships)
	if err != nil {
		return NewErrResponse(c, err)
	}

	groupNames, err := h.getGroupNamesForMemberships(memberships)
	if err != nil {
		return NewErrResponse(c, err)
	}

	response := web.NewGetUserMembershipsResponse(memberships, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}

func (h *Handler) getUserNamesForMemberships(memberships []model.Membership) (model.UserNames, error) {
	var userNames = model.UserNames{}
	for _, membership := range memberships {
		userKey := membership.GetUserKey()
		_, ok := userNames[userKey]
		if !ok {
			username, err := h.authStore.GetUsername(userKey)
			if err != nil {
				return userNames, err
			}
			userNames[userKey] = username
		}
	}
	return userNames, nil
}

func (h *Handler) getGroupNamesForMemberships(memberships []model.Membership) (model.GroupNames, error) {
	var groupNames = model.GroupNames{}
	for _, membership := range memberships {
		groupKey := membership.GetGroupKey()
		_, ok := groupNames[groupKey]
		if !ok {
			getGroup := h.groupStore.GetGroup(group.NewGetGroupRequest(groupKey))
			if getGroup.Error != nil {
				return groupNames, getGroup.Error
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

	authUserKey := h.authorization.GetAuthUserKey(c)

	req := web.InviteUserRequest{}
	if err := c.Bind(&req); err != nil {
		return NewErrResponse(c, err)
	}

	// Retrieve the group
	groupKey, err := model.ParseGroupKey(c.Param("id"))
	if err != nil {
		return NewErrResponse(c, err)
	}
	getGroup := h.groupStore.GetGroup(group.NewGetGroupRequest(groupKey))
	if getGroup.Error != nil {
		return NewErrResponse(c, getGroup.Error)
	}
	groupNames := model.GroupNames{
		groupKey: getGroup.Group.Name,
	}

	// Retrieve the invited user key
	userKey := model.NewUserKey(req.UserID)
	username, err := h.authStore.GetUsername(userKey)
	if err != nil {
		return NewErrResponse(c, err)
	}
	userNames := model.UserNames{
		userKey: username,
	}

	// Check that inviter can actually invite on that group
	authMembershipKey := model.NewMembershipKey(groupKey, authUserKey)
	authPermissions := h.groupStore.GetGroupPermissionsForUser(group.NewGetMembershipPermissionsRequest(authMembershipKey))
	if authPermissions.Error != nil {
		return NewErrResponse(c, authPermissions.Error)
	}
	if !authPermissions.MembershipPermissions.IsAdmin {
		return NewErrResponse(c, fmt.Errorf("forbidden"))
	}

	// Create the membership
	membershipKey := model.NewMembershipKey(groupKey, userKey)
	invite := h.groupStore.Invite(group.InviteRequest{
		MembershipKey: membershipKey,
		InvitedBy:     group.GroupParty,
	})

	if invite.Error != nil {
		return NewErrResponse(c, invite.Error)
	}

	// Retrieve the created membership
	getMembership := h.groupStore.GetMembership(group.NewGetMembershipRequest(membershipKey))
	if getMembership.Error != nil {
		return NewErrResponse(c, getMembership.Error)
	}

	// Respond to query
	response := web.NewInviteUserResponse(getMembership.Membership, groupNames, userNames)
	return c.JSON(http.StatusAccepted, response)

}
