package handler

import (
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

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
func (h *Handler) GetUserInfo(c echo.Context) error {

	userId := c.Param("id")
	userKey := usermodel.NewUserKey(userId)

	user, err := h.authStore.GetByKey(userKey)

	if err != nil {
		return NewErrResponse(c, err)
	}

	response := web.UserInfoResponse{
		Username: user.Username,
		Id:       user.ID,
	}
	return c.JSON(http.StatusOK, response)

}
