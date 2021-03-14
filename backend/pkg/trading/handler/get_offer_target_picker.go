package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/labstack/echo/v4"
	"net/http"
)

type OfferGroupOrUserPickerItem struct {
	Type    keys.TargetType `json:"type"`
	UserID  *string         `json:"userId"`
	GroupID *string         `json:"groupId"`
	Name    string          `json:"name"`
}

type OfferGroupOrUserPickerResult struct {
	Items []OfferGroupOrUserPickerItem `json:"items"`
}

func (h *TradingHandler) HandleOfferItemTargetPicker(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleOfferItemTargetPicker")

	groupKey, err := keys.ParseGroupKey(c.QueryParams().Get("group_id"))
	if err != nil {
		return err
	}

	offerItemType, err := domain.ParseOfferItemType(c.QueryParams().Get("type"))
	if err != nil {
		return err
	}

	from, err := parseTargetFromQueryParams(c, "from_type", "from_id")
	if err != nil {
		return err
	}

	to, err := parseTargetFromQueryParams(c, "to_type", "to_id")
	if err != nil {
		return err
	}

	targets, err := h.tradingService.FindTargetsForOfferItem(ctx, groupKey, offerItemType, from, to)
	if err != nil {
		return err
	}

	var items []OfferGroupOrUserPickerItem

	groups, err := h.groupService.GetGroupsByKeys(ctx, targets.GetGroupKeys())
	if err != nil {
		return err
	}

	for _, group := range groups {
		groupId := group.GroupKey.String()
		items = append(items, OfferGroupOrUserPickerItem{
			Type:    keys.GroupTarget,
			GroupID: &groupId,
			Name:    group.Name,
		})
	}

	users, err := h.getUsersByKeys.Get(targets.GetUserKeys())
	if err != nil {
		return err
	}

	for userKey, user := range users {
		userID := userKey.String()
		items = append(items, OfferGroupOrUserPickerItem{
			Type:   keys.UserTarget,
			UserID: &userID,
			Name:   user.Username,
		})
	}

	result := &OfferGroupOrUserPickerResult{
		Items: items,
	}

	return c.JSON(http.StatusOK, result)
}
