package handler

import (
	"cp/pkg/acknowledgements"
	"cp/pkg/api"
	"cp/pkg/credits"
	"cp/pkg/groups"
	"cp/pkg/images"
	"cp/pkg/memberships"
	"cp/pkg/messages"
	"cp/pkg/notifications"
	"cp/pkg/posts"
	"cp/pkg/users"
	"cp/pkg/utils"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"os"
)

const (
	GroupKey                       = "Group"
	GroupIDKey                     = "GroupID"
	UserKey                        = "User"
	UserIDKey                      = "UserID"
	MembershipKey                  = "Membership"
	PostIDKey                      = "PostID"
	PostKey                        = "Post"
	AuthenticatedUserKey           = "AuthenticatedUser"
	AuthenticatedUserMembershipKey = "AuthenticatedUserMembership"
	ProfileKey                     = "Profile"
)

type Handler struct {
	cookieStore          *sessions.CookieStore
	groupStore           groups.Store
	membershipStore      memberships.Store
	userStore            users.Store
	postStore            posts.Store
	creditsStore         credits.Store
	acknowledgementStore acknowledgements.Store
	messageStore         messages.Store
	notificationStore    notifications.Store
	imageStore           images.Store
	alertManager         *utils.AlertManager
	db                   *gorm.DB
}

func NewHandler(
	cookieStore *sessions.CookieStore,
	groupStore groups.Store,
	membershipStore memberships.Store,
	userStore users.Store,
	postStore posts.Store,
	creditsStore credits.Store,
	acknowledgementStore acknowledgements.Store,
	messageStore messages.Store,
	notificationStore notifications.Store,
	imageStore images.Store,
	alertManager *utils.AlertManager,
	db *gorm.DB) *Handler {
	return &Handler{
		cookieStore:          cookieStore,
		groupStore:           groupStore,
		membershipStore:      membershipStore,
		userStore:            userStore,
		postStore:            postStore,
		creditsStore:         creditsStore,
		acknowledgementStore: acknowledgementStore,
		messageStore:         messageStore,
		imageStore:           imageStore,
		notificationStore:    notificationStore,
		alertManager:         alertManager,
		db:                   db,
	}
}

func (h *Handler) getGroup(c echo.Context) (*api.Group, error) {
	if group, ok := c.Get(GroupKey).(*api.Group); ok {
		return group, nil
	}
	return nil, echo.ErrInternalServerError
}

func (h *Handler) getGroupID(c echo.Context) (string, error) {
	if groupID, ok := c.Get(GroupIDKey).(string); ok {
		return groupID, nil
	}
	return "", echo.ErrInternalServerError
}

func (h *Handler) groupM() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			groupID := c.Param(GroupIDKey)
			if groupID == "" {
				return echo.ErrNotFound
			}
			c.Set(GroupIDKey, groupID)

			group, err := h.groupStore.Get(groupID)
			if err != nil {
				return echo.ErrNotFound
			}
			c.Set(GroupKey, group)

			return handlerFunc(c)
		}
	}
}

func (h *Handler) getUser(c echo.Context) (*api.User, error) {
	if user, ok := c.Get(UserKey).(*api.User); ok {
		return user, nil
	}
	return nil, echo.ErrInternalServerError
}

func (h *Handler) userM() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Param(UserIDKey)
			if userID == "" {
				return echo.ErrNotFound
			}
			c.Set(UserIDKey, userID)
			user, err := h.userStore.Get(userID)
			if err != nil {
				return echo.ErrNotFound
			}
			c.Set(UserKey, user)
			return handlerFunc(c)
		}
	}
}

func (h *Handler) getMembership(c echo.Context) (*api.Membership, error) {
	if membership, ok := c.Get(MembershipKey).(*api.Membership); ok {
		return membership, nil
	}
	return nil, echo.ErrInternalServerError
}

func (h *Handler) memberM(optional bool) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			groupID := c.Param(GroupIDKey)
			userID := c.Param(UserIDKey)

			if groupID == "" {
				return echo.ErrBadRequest
			}
			if userID == "" {
				return echo.ErrBadRequest
			}

			membership, err := h.membershipStore.Get(groupID, userID)
			if optional && errors.Is(err, echo.ErrNotFound) {
				var m *api.Membership
				c.Set(MembershipKey, m)
			} else if err != nil {
				return err
			} else {
				c.Set(MembershipKey, membership)
			}

			return handlerFunc(c)
		}
	}
}

func (h *Handler) getPost(c echo.Context) (*api.Post, error) {
	if post, ok := c.Get(PostKey).(*api.Post); ok {
		return post, nil
	}
	return nil, echo.ErrInternalServerError
}

func (h *Handler) postM(optional bool) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			groupID, err := h.getGroupID(c)
			if err != nil {
				return err
			}

			postID := c.Param(PostIDKey)
			c.Set(PostIDKey, "")

			if postID == "" {
				if optional {
					var post *api.Post
					c.Set(PostKey, post)
					return handlerFunc(c)
				} else {
					return echo.ErrBadRequest
				}
			}

			c.Set(PostIDKey, postID)
			post, err := h.postStore.Get(postID)
			if optional && errors.Is(err, echo.ErrNotFound) {
				c.Set(PostIDKey, "")
				c.Set(PostKey, nil)
				return handlerFunc(c)
			}
			if err != nil {
				return err
			}
			if post.GroupID != groupID {
				return echo.ErrBadRequest
			}
			c.Set(PostIDKey, post.ID)
			c.Set(PostKey, post)
			return handlerFunc(c)
		}
	}
}

func (h *Handler) getAuthenticatedUser(c echo.Context) (*api.User, error) {
	if authenticatedUser, ok := c.Get(AuthenticatedUserKey).(*api.User); ok {
		if authenticatedUser == nil {
			return nil, echo.ErrUnauthorized
		}
		return authenticatedUser, nil
	}
	return nil, echo.ErrInternalServerError
}

func (h *Handler) authM(optional bool) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			profile, err := GetProfile(h.cookieStore, c)
			if err != nil {
				return err
			}

			if !optional && profile == nil {
				return echo.ErrUnauthorized
			}

			if profile == nil {
				var authenticatedUser *api.User
				c.Set(AuthenticatedUserKey, authenticatedUser)
				c.Set(ProfileKey, profile)
				return handlerFunc(c)
			}

			c.Set(ProfileKey, profile)

			user := &api.User{
				ID:       profile.ID,
				Username: profile.Username,
				Email:    profile.Email,
			}
			if err := h.userStore.Upsert(user); err != nil {
				return err
			}

			user, err = h.userStore.Get(user.ID)
			if err != nil {
				return err
			}

			c.Set(AuthenticatedUserKey, user)

			return handlerFunc(c)

		}
	}
}

func (h *Handler) getAuthenticatedUserMembership(c echo.Context) (*api.Membership, error) {
	if membership, ok := c.Get(AuthenticatedUserMembershipKey).(*api.Membership); ok {
		return membership, nil
	}
	return nil, echo.ErrInternalServerError
}

func (h *Handler) authMemberM(optional bool) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			groupID, err := h.getGroupID(c)
			if err != nil {
				return err
			}

			authenticatedUser, err := h.getAuthenticatedUser(c)
			if err != nil {
				return err
			}

			membership, err := h.membershipStore.Get(groupID, authenticatedUser.ID)
			if optional && errors.Is(err, echo.ErrNotFound) {
				var m *api.Membership
				c.Set(AuthenticatedUserMembershipKey, m)
			} else if err != nil {
				return err
			} else {
				c.Set(AuthenticatedUserMembershipKey, membership)
			}

			return handlerFunc(c)
		}
	}
}

func (h *Handler) isInGroupM(group string) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			profile, err := GetProfile(h.cookieStore, c)
			if err != nil {
				return err
			}
			if profile == nil {
				return echo.ErrUnauthorized
			}
			if !profile.IsInGroup(group) {
				return echo.ErrForbidden
			}

			return handlerFunc(c)
		}
	}
}

func (h *Handler) Register(e *echo.Echo) {

	uploadDir := os.Getenv("PUBLIC_DIR")
	if uploadDir == "" {
		uploadDir = "public"
	}
	e.Static("/", uploadDir)

	e.GET("/", h.handleHomeView, h.authM(true)).Name = "get_home"

	a := e.Group("/auth")
	a.GET("/login", h.handleLogin).Name = "get_auth_login"
	a.GET("/logout", h.handleLogout).Name = "get_auth_logout"
	a.GET("/callback", h.handleOauthCallback).Name = "get_auth_callback"

	gs := e.Group("/groups", h.authM(false))
	gs.GET("", h.handleGroupsView).Name = "get_groups"
	gs.GET("/new", h.handleNewGroup).Name = "get_groups_new"
	gs.POST("/new", h.handleNewGroup).Name = "post_groups_new"
	gs.GET("/edit", h.handleEditGroup).Name = "get_groups_edit"

	g := gs.Group(fmt.Sprintf("/:%s", GroupIDKey), h.groupM())
	g.GET("", h.handleGroupPostsView, h.authMemberM(true)).Name = "get_group_posts"
	g.GET("/send", h.handleGroupSend, h.authMemberM(false)).Name = "get_group_send"
	g.POST("/send", h.handleGroupSend, h.authMemberM(false)).Name = "post_group_send"
	g.GET("/members", h.handleGroupMembersView, h.authMemberM(false)).Name = "get_group_members"
	g.GET("/acknowledgements", h.handleGetGroupAcknowledgements, h.authMemberM(false)).Name = "get_group_acknowledgements"
	g.GET("/settings", h.handleGroupSettings, h.authMemberM(false)).Name = "get_group_settings"
	g.POST("/delete", h.handleGroupDelete, h.authMemberM(false)).Name = "post_group_delete"
	g.GET("/history", h.handleGetGroupHistory, h.authMemberM(false)).Name = "get_group_history"
	g.GET("/posts/new", h.handlePostEdit, h.authMemberM(false), h.postM(true)).Name = "get_group_post_new"
	g.POST("/posts/new", h.handlePostEdit, h.authMemberM(false), h.postM(true)).Name = "post_group_post_new"

	p := g.Group(fmt.Sprintf("/posts/:%s", PostIDKey), h.postM(false))
	p.GET("", h.handlePostView, h.authMemberM(true)).Name = "get_group_post"
	p.GET("/edit", h.handlePostEdit, h.authMemberM(false)).Name = "get_group_post_edit"
	p.POST("/edit", h.handlePostEdit, h.authMemberM(false)).Name = "post_group_form_edit"
	p.POST("/delete", h.handlePostDelete, h.authMemberM(false)).Name = "post_group_delete"
	p.POST("/message", h.handlePostMessage, h.authMemberM(false)).Name = "post_group_post_message"

	m := g.Group(fmt.Sprintf("/users/:%s", UserIDKey), h.userM())
	m.POST("/join", h.handleGroupJoin, h.authMemberM(true), h.memberM(true)).Name = "post_group_join"
	m.POST("/leave", h.handleGroupLeave, h.authMemberM(false), h.memberM(false)).Name = "post_group_leave"
	m.POST("/permissions", h.handleGroupSetPermission, h.authMemberM(false), h.memberM(false)).Name = "post_group_permissions"

	u := e.Group(fmt.Sprintf("/users/:%s", UserIDKey), h.authM(false), h.userM())
	u.GET("", h.handleGetUserPosts).Name = "get_user_posts"
	u.GET("/groups", h.handleGetUserGroups).Name = "get_user_groups"
	u.GET("/notifications", h.handleGetUserNotifications).Name = "get_user_notifications"
	u.GET("/acknowledgements", h.handleGetUserAcknowledgements).Name = "get_user_acknowledgements"
	u.GET("/profile", h.handleGetUserProfile).Name = "get_user_profile"
	u.GET("/profile/edit", h.handleEditUserProfile).Name = "get_user_profile_edit"
	u.POST("/profile/edit", h.handleEditUserProfile).Name = "post_user_profile_edit"

	adm := e.Group("/admin", h.authM(false), h.isInGroupM("administrators"))
	adm.GET("", h.handleAdmin)
	adm.POST("/clear", h.handleAdminClearAll)
}
