package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/utils"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetRecentlyActiveSubscriptions
// @Summary Returns the latest user message threads
// @Description This endpoint returns the latest messaging threads for the currently logged in user.
// @ID getLatestThreads
// @Param take query int false "Number of threads to take" minimum(0) maximum(100) default(10)
// @Param skip query int false "Number of threads to skip" minimum(0) default(0)
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {object} web.GetLatestSubscriptionsResponse
// @Failure 400 {object} utils.Error
// @Router /chat/subscriptions [get]
func (chatHandler *ChatHandler) GetRecentlyActiveSubscriptions(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetRecentlyActiveSubscriptions")

	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	channelSubscriptions, err := chatHandler.chatService.GetSubscriptionsForUser(ctx, take, skip)
	if err != nil {
		return err
	}

	mappedSubscriptions, err := MapChannelSubscriptions(ctx, chatHandler.chatService, channelSubscriptions)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, web.GetLatestSubscriptionsResponse{
		Subscriptions: mappedSubscriptions,
	})

}
