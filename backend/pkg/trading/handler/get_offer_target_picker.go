package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/labstack/echo/v4"
	"net/http"
)

type OfferGroupOrUserPickerItem struct {
	Type    resource.TargetType `json:"type"`
	UserID  *string             `json:"userId"`
	GroupID *string             `json:"groupId"`
	Name    string              `json:"name"`
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

	offerItemType, err := trading.ParseOfferItemType(c.QueryParams().Get("type"))
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

	for _, group := range groups.Items {
		groupId := group.GetKey().String()
		items = append(items, OfferGroupOrUserPickerItem{
			Type:    resource.GroupTarget,
			GroupID: &groupId,
			Name:    group.Name,
		})
	}

	users, err := h.userService.GetByKeys(ctx, targets.GetUserKeys())
	if err != nil {
		return err
	}

	for _, item := range users.Items {
		userKey := item.GetUserKey().String()
		items = append(items, OfferGroupOrUserPickerItem{
			Type:   resource.UserTarget,
			UserID: &userKey,
			Name:   item.Username,
		})
	}

	result := &OfferGroupOrUserPickerResult{
		Items: items,
	}

	return c.JSON(http.StatusOK, result)
}
