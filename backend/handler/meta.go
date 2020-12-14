package handler

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// WhoAmI godoc
// @Summary Returns information about myself
// @Description Returns information about the currently authenticated user
// @ID whoAmI
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} web.UserAuthResponse
// @Failure 400 {object} utils.Error
// @Router /meta/who-am-i [get]
func (h *Handler) WhoAmI(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "WhoAmI")

	loggedInUser, err := auth.GetLoggedInUser(ctx)

	if err != nil {
		return c.JSON(http.StatusOK, &web.UserAuthResponse{
			IsAuthenticated: false,
		})
	}

	return c.JSON(http.StatusOK, web.UserAuthResponse{
		IsAuthenticated: true,
		Username:        loggedInUser.Username,
		Id:              loggedInUser.Subject,
	})
}
