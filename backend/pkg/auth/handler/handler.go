package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/auth/service"
	"github.com/commonpool/backend/pkg/auth/store"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	userStore     store.Store
	userService   service.Service
	authenticator authenticator.Authenticator
	appConfig     *config.AppConfig
}

func NewUserHandler(appConfig *config.AppConfig, userService service.Service, authenticator authenticator.Authenticator) *UserHandler {
	return &UserHandler{
		userService:   userService,
		authenticator: authenticator,
		appConfig:     appConfig,
	}
}

func (h *UserHandler) Register(g *echo.Group) {
	usersGroup := g.Group("/users", h.authenticator.Authenticate(true))
	usersGroup.GET("", h.SearchUsers)
	usersGroup.GET("/:id", h.GetUserInfo)
	authGroup := g.Group("/auth")
	authGroup.Any("/login", h.authenticator.Login())
	authGroup.Any("/logout", h.authenticator.Logout())
	h.authenticator.Register(g)
	session := g.Group("/session", h.authenticator.Authenticate(false))
	session.GET("/info", h.SessionInfo)
}

type UsersInfoResponse struct {
	Users []UserInfoResponse `json:"users"`
	Take  int                `json:"take"`
	Skip  int                `json:"skip"`
}

type UserInfoResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type UserAuthResponse struct {
	IsAuthenticated bool   `json:"isAuthenticated"`
	Username        string `json:"username"`
	Id              string `json:"id"`
}

// GetUserInfo godoc
// @Summary Returns information about a user
// @Description Returns information about the given user
// @ID getUserInfo
// @Param id path string true "User id" format(uuid)
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} web.UserInfoResponse
// @Failure 400 {object} utils.Error
// @Router /users/:id [get]
func (h *UserHandler) GetUserInfo(c echo.Context) error {

	userId := c.Param("id")
	userKey := keys.NewUserKey(userId)

	user, err := h.userService.GetUser(userKey)

	if err != nil {
		return err
	}

	response := UserInfoResponse{
		Username: user.Username,
		Id:       user.ID,
	}
	return c.JSON(http.StatusOK, response)

}

// GetUserInfo godoc
// @Summary Finds users
// @Description Paginated query on users
// @ID searchUsers
// @Param id query string false "Query"
// @Param take query int false "Take"
// @Param skip query int false "Skip"
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} web.UsersInfoResponse
// @Failure 400 {object} utils.Error
// @Router /users [get]
func (h *UserHandler) SearchUsers(c echo.Context) error {
	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	qry := c.QueryParam("query")

	userQuery := store.Query{
		Query: qry,
		Skip:  skip,
		Take:  take,
	}

	users, err := h.userService.Find(userQuery)
	if err != nil {
		return err
	}

	responseItems := make([]UserInfoResponse, len(users.Items))
	for i, u := range users.Items {
		responseItems[i] = UserInfoResponse{
			Id:       u.ID,
			Username: u.Username,
		}
	}

	response := UsersInfoResponse{
		Users: responseItems,
		Take:  take,
		Skip:  skip,
	}

	return c.JSON(http.StatusOK, response)

}

// SessionInfo godoc
// @Summary Returns information about myself
// @Description Returns information about the currently authenticated user
// @ID whoAmI
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} web.UserAuthResponse
// @Failure 400 {object} utils.Error
// @Router /meta/who-am-i [get]
func (h *UserHandler) SessionInfo(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "SessionInfo")

	loggedInUser, err := oidc.GetLoggedInUser(ctx)

	if err != nil {
		return c.JSON(http.StatusOK, &UserAuthResponse{
			IsAuthenticated: false,
		})
	}

	return c.JSON(http.StatusOK, UserAuthResponse{
		IsAuthenticated: true,
		Username:        loggedInUser.Username,
		Id:              loggedInUser.Subject,
	})
}
