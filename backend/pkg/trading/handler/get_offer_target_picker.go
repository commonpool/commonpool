package handler

import (
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/resource/model"
	model2 "github.com/commonpool/backend/pkg/trading/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *TradingHandler) HandleOfferItemTargetPicker(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleOfferItemTargetPicker")

	groupKey, err := groupmodel.ParseGroupKey(c.QueryParams().Get("group_id"))
	if err != nil {
		return err
	}

	offerItemType, err := model2.ParseOfferItemType(c.QueryParams().Get("type"))
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

	var items []web.OfferGroupOrUserPickerItem

	groups, err := h.groupService.GetGroupsByKeys(ctx, targets.GetGroupKeys())
	if err != nil {
		return err
	}

	for _, group := range groups.Items {
		groupId := group.GetKey().String()
		items = append(items, web.OfferGroupOrUserPickerItem{
			Type:    model.GroupTarget,
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
		items = append(items, web.OfferGroupOrUserPickerItem{
			Type:   model.UserTarget,
			UserID: &userKey,
			Name:   item.Username,
		})
	}

	result := &web.OfferGroupOrUserPickerResult{
		Items: items,
	}

	return c.JSON(http.StatusOK, result)
}
