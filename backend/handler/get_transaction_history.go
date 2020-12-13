package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) GetTradingHistory(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "GetTradingHistory")

	req := web.GetTradingHistoryRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	var userKeys []model.UserKey
	for _, userId := range req.UserIDs {
		userKey := model.NewUserKey(userId)
		userKeys = append(userKeys, userKey)
	}

	tradingHistory, err := h.tradingService.GetTradingHistory(ctx, model.NewUserKeys(userKeys))
	if err != nil {
		return err
	}

	tradingUserKeys := model.NewUserKeys([]model.UserKey{})
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
