package handler

import (
	"github.com/commonpool/backend/model"
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
	userKey := model.NewUserKey(userId)

	user := &model.User{}
	err := h.authStore.GetByKey(userKey, user)

	if err != nil {
		return NewErrResponse(c, err)
	}

	response := web.UserInfoResponse{
		Username: user.Username,
		Id:       user.ID,
	}
	return c.JSON(http.StatusOK, response)

}
