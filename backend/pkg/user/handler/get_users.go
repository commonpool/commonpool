package handler

import (
	"github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/pkg/user"
	"github.com/commonpool/backend/pkg/utils"
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
func (h *UserHandler) SearchUsers(c echo.Context) error {
	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return handler.NewErrResponse(c, err)
	}

	qry := c.QueryParam("query")

	userQuery := user.Query{
		Query: qry,
		Skip:  skip,
		Take:  take,
	}

	users, err := h.userService.Find(userQuery)
	if err != nil {
		return err
	}

	responseItems := make([]web.UserInfoResponse, len(users.Items))
	for i, u := range users.Items {
		responseItems[i] = web.UserInfoResponse{
			Id:       u.ID,
			Username: u.Username,
		}
	}

	response := web.UsersInfoResponse{
		Users: responseItems,
		Take:  take,
		Skip:  skip,
	}

	return c.JSON(http.StatusOK, response)

}
