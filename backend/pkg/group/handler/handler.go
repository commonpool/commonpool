package handler

import (
	"github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	handler2 "github.com/commonpool/backend/pkg/auth/handler"
	"github.com/commonpool/backend/pkg/exceptions"
	group "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/queries"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Handler struct {
	groupService                 group.Service
	auth                         authenticator.Authenticator
	getGroup                     *queries.GetGroup
	getMembership                *queries.GetMembershipReadModel
	getGroupMemberships          *queries.GetGroupMemberships
	getUserMemberships           *queries.GetUserMemberships
	getUsersForGroupInvitePicker *queries.GetUsersForGroupInvite
}

func NewHandler(
	groupService group.Service,
	auth authenticator.Authenticator,
	getGroup *queries.GetGroup,
	getMembership *queries.GetMembershipReadModel,
	getGroupMemberships *queries.GetGroupMemberships,
	getUserMemberships *queries.GetUserMemberships,
	getUsersForGroupInvitePicker *queries.GetUsersForGroupInvite) *Handler {
	return &Handler{
		groupService:                 groupService,
		auth:                         auth,
		getGroup:                     getGroup,
		getMembership:                getMembership,
		getGroupMemberships:          getGroupMemberships,
		getUserMemberships:           getUserMemberships,
		getUsersForGroupInvitePicker: getUsersForGroupInvitePicker,
	}
}

func (h *Handler) Register(g *echo.Group) {

	groups := g.Group("/groups", h.auth.Authenticate(true))
	groups.POST("", h.CreateGroup)
	groups.GET("/:id", h.GetGroup)
	groups.GET("/:id/memberships", h.GetGroupMemberships)
	groups.GET("/:id/memberships/:userId", h.GetMembership)
	groups.GET("/:id/invite-member-picker", h.GetUsersForGroupInvitePicker)

	memberships := g.Group("/memberships", h.auth.Authenticate(true))
	memberships.GET("", h.GetUserMemberships)
	memberships.POST("", h.CreateOrAcceptMembership)
	memberships.DELETE("", h.CancelOrDeclineInvitation)

}

type GetMembershipResponse struct {
	Membership *readmodels.MembershipReadModel `json:"membership"`
}
type GetMembershipsResponse struct {
	Memberships []*readmodels.MembershipReadModel `json:"memberships"`
}
type GetGroupResponse struct {
	Group *readmodels.GroupReadModel `json:"group"`
}
type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type CreateOrAcceptInvitationRequest struct {
	UserKey  keys.UserKey  `json:"userId"`
	GroupKey keys.GroupKey `json:"groupId"`
}
type CancelOrDeclineInvitationRequest struct {
	UserKey  keys.UserKey  `json:"userId"`
	GroupKey keys.GroupKey `json:"groupId"`
}
type GetUsersForGroupInvitePickerResponse struct {
	Users []*handler2.UserInfoResponse `json:"users"`
	Take  int                          `json:"take"`
	Skip  int                          `json:"skip"`
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

	ctx, _ := handler.GetEchoContext(c, "GetGroup")

	groupKey, err := keys.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}
	rm, err := h.getGroup.Get(ctx, groupKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetGroupResponse{
		Group: rm,
	})

}

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

	ctx, _ := handler.GetEchoContext(c, "CreateGroup")

	req := CreateGroupRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return exceptions.ErrValidation("name is required")
	}

	groupKey, err := h.groupService.CreateGroup(ctx, group.NewCreateGroupRequest(req.Name, req.Description))
	if err != nil {
		return err
	}

	errChan := make(chan error)
	grpChan := make(chan *readmodels.GroupReadModel)

	go func() {
		err = retry.Do(func() error {
			grp, err := h.groupService.GetGroup(ctx, groupKey)
			if err != nil {
				return err
			}
			grpChan <- grp
			return nil
		})
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case grp := <-grpChan:
		return c.JSON(http.StatusCreated, GetGroupResponse{
			Group: grp,
		})
	}

}

// GetMembership godoc
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

	ctx, _ := handler.GetEchoContext(c, "GetMembership")

	userKey := keys.NewUserKey(c.Param("userId"))

	groupKey, err := keys.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	membershipKey := keys.NewMembershipKey(groupKey, userKey)
	var membership *readmodels.MembershipReadModel
	err = retry.Do(func() error {
		var err error
		membership, err = h.getMembership.Get(ctx, membershipKey)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	resp := &GetMembershipResponse{
		Membership: membership,
	}

	return c.JSON(http.StatusOK, resp)

}

// GetUserMemberships godoc
// @Summary Gets memberships for a given user
// @Description Gets the memberships for a given user
// @ID getUserMemberships
// @Param user_id query string false "ID of the user. If not set, defaults to the logged in user id" (format:uuid)
// @Param status status MembershipStatus true "status of the membership"
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.GetMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /memberships [get]
func (h *Handler) GetUserMemberships(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetUserMemberships")

	var membershipStatus = domain.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := domain.ParseMembershipStatus(statusStr)
		if err != nil {
			return err
		}
		membershipStatus = &ms
	}

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	userKey := loggedInUser.GetUserKey()

	userIdStr := c.Param("user_id")
	if userIdStr != "" {
		userKey = keys.NewUserKey(userIdStr)
	}

	m, err := h.getUserMemberships.Get(ctx, userKey, membershipStatus)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GetMembershipsResponse{
		Memberships: m,
	})

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

	ctx, _ := handler.GetEchoContext(c, "GetGroupMemberships")

	var membershipStatus = domain.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := domain.ParseMembershipStatus(statusStr)
		if err != nil {
			return err
		}
		membershipStatus = &ms
	}

	groupKey, err := keys.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	m, err := h.getGroupMemberships.Get(ctx, groupKey, membershipStatus)
	if err != nil {
		return err
	}

	response := GetMembershipsResponse{
		Memberships: m,
	}
	return c.JSON(http.StatusOK, response)

}

// CreateOrAcceptMembership godoc
// @Summary Accept a group invitation
// @Description Accept a group invitation
// @ID acceptInvitation
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.CreateOrAcceptInvitationResponse
// @Failure 400 {object} utils.Error
// @Router /groups/memberships [post]
func (h *Handler) CreateOrAcceptMembership(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "CreateOrAcceptInvitation")
	req := CreateOrAcceptInvitationRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	membershipKey := keys.NewMembershipKey(req.GroupKey, req.UserKey)
	err := h.groupService.CreateOrAcceptInvitation(ctx, group.NewAcceptInvitationRequest(membershipKey))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)

}

// CancelOrDeclineInvitation godoc
// @Summary declines a group invitation
// @Description declines a group invitation
// @ID declineInvitation
// @Tags groups
// @Accept json
// @Produce json
// @Success 202 {object} web.CancelOrDeclineInvitationResponse
// @Failure 400 {object} utils.Error
// @Router /memberships [delete]
func (h *Handler) CancelOrDeclineInvitation(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "CancelOrDeclineInvitation")
	req := CancelOrDeclineInvitationRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	membershipKey := keys.NewMembershipKey(req.GroupKey, req.UserKey)
	err := h.groupService.CancelOrDeclineInvitation(ctx, group.NewDeclineInvitationRequest(membershipKey))
	if err != nil {
		l.Error("could not decline invitation", zap.Error(err))
		return err
	}
	return c.NoContent(http.StatusAccepted)

}

// GetUsersForGroupInvitePicker godoc
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
	ctx := handler.GetContext(c)

	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	qry := c.QueryParam("query")

	groupKey, err := keys.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	users, err := h.getUsersForGroupInvitePicker.Get(ctx, groupKey, qry, skip, take)
	if err != nil {
		return err
	}

	responseItems := make([]*handler2.UserInfoResponse, len(users))
	for i, u := range users {
		responseItems[i] = &handler2.UserInfoResponse{
			Id:       u.UserKey,
			Username: u.Username,
		}
	}

	response := GetUsersForGroupInvitePickerResponse{
		Users: responseItems,
		Take:  take,
		Skip:  skip,
	}

	return c.JSON(http.StatusOK, response)

}
