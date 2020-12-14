package handler

import (
	"github.com/commonpool/backend/pkg/user"
	"github.com/commonpool/backend/utils"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

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
func (h *Handler) SearchUsers(c echo.Context) error {
	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return NewErrResponse(c, err)
	}

	qry := c.QueryParam("query")

	userQuery := user.UserQuery{
		Query: qry,
		Skip:  skip,
		Take:  take,
	}

	users, err := h.authStore.Find(userQuery)
	if err != nil {
		return err
	}

	responseItems := make([]web.UserInfoResponse, len(users))
	for i, user := range users {
		responseItems[i] = web.UserInfoResponse{
			Id:       user.ID,
			Username: user.Username,
		}
	}

	response := web.UsersInfoResponse{
		Users: responseItems,
		Take:  take,
		Skip:  skip,
	}

	return c.JSON(http.StatusOK, response)

}
