package chat

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
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
func (h *ChatHandler) GetRecentlyActiveSubscriptions(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetRecentlyActiveSubscriptions")

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return errors.ErrUnauthorized
	}
	loggedInUserKey := model.NewUserKey(loggedInUser.Subject)

	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	userSubscriptions, err := h.chatService.GetUserSubscriptions(ctx, loggedInUserKey, take, skip)
	if err != nil {
		return err
	}

	var items []web.Subscription
	for _, subscription := range userSubscriptions.Items {
		channel, err := h.chatService.GetChannel(ctx, subscription.GetChannelKey())
		if err != nil {
			return err
		}
		items = append(items, *web.MapSubscription(channel, &subscription))
	}

	if items == nil {
		items = []web.Subscription{}
	}

	return c.JSON(http.StatusOK, web.GetLatestSubscriptionsResponse{
		Subscriptions: items,
	})

}
