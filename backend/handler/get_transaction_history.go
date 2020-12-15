package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) GetTradingHistory(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetTradingHistory")

	req := web.GetTradingHistoryRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	var userKeys []usermodel.UserKey
	for _, userId := range req.UserIDs {
		userKey := usermodel.NewUserKey(userId)
		userKeys = append(userKeys, userKey)
	}

	tradingHistory, err := h.tradingService.GetTradingHistory(ctx, usermodel.NewUserKeys(userKeys))
	if err != nil {
		return err
	}

	tradingUserKeys := usermodel.NewUserKeys([]usermodel.UserKey{})
	for _, entry := range tradingHistory {
		tradingUserKeys = tradingUserKeys.Append(entry.ToUserID)
		tradingUserKeys = tradingUserKeys.Append(entry.FromUserID)
	}

	users, err := h.authStore.GetByKeys(ctx, tradingUserKeys.Items)
	if err != nil {
		return err
	}

	var responseEntries []web.TradingHistoryEntry
	for _, entry := range tradingHistory {
		var resourceId *string
		if entry.ResourceID != nil {
			resourceIdStr := entry.ResourceID.String()
			resourceId = &resourceIdStr
		}
		fromUser, err := users.GetUser(entry.FromUserID)
		if err != nil {
			return err
		}
		toUser, err := users.GetUser(entry.ToUserID)
		if err != nil {
			return err
		}
		webEntry := web.TradingHistoryEntry{
			Timestamp:         entry.Timestamp.String(),
			FromUserID:        entry.FromUserID.String(),
			FromUsername:      fromUser.Username,
			ToUserID:          entry.ToUserID.String(),
			ToUsername:        toUser.Username,
			ResourceID:        resourceId,
			TimeAmountSeconds: entry.TimeAmountSeconds,
		}
		responseEntries = append(responseEntries, webEntry)
	}

	return c.JSON(http.StatusOK, web.GetTradingHistoryResponse{
		Entries: responseEntries,
	})
}
