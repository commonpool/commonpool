package handler

import (
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
	userAuth := h.authorization.GetAuthUserSession(c)
	response := web.UserAuthResponse{
		IsAuthenticated: userAuth.IsAuthenticated,
		Username:        userAuth.Username,
		Id:              userAuth.Subject,
	}
	return c.JSON(http.StatusOK, response)
}
