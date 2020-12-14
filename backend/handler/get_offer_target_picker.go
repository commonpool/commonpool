package handler

import (
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/trading"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) HandleOfferItemTargetPicker(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleOfferItemTargetPicker")

	groupKey, err := group.ParseGroupKey(c.QueryParams().Get("group_id"))
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

	items := []web.OfferGroupOrUserPickerItem{}

	userKeys := []model.UserKey{}
	groupKeys := []model.GroupKey{}
	for _, target := range targets.Items {
		if target.IsForUser() {
			userKeys = append(userKeys, target.GetUserKey())
		} else if target.IsForGroup() {
			groupKeys = append(groupKeys, target.GetGroupKey())
		}
	}

	for _, grpKey := range groupKeys {
		grp, err := h.groupService.GetGroup(ctx, &group.GetGroupRequest{
			Key: grpKey,
		})
		if err != nil {
			return err
		}
		groupId := grp.Group.GetKey().String()
		items = append(items, web.OfferGroupOrUserPickerItem{
			Type:    model.GroupTarget,
			GroupID: &groupId,
			Name:    grp.Group.Name,
		})
	}

	users, err := h.authStore.GetByKeys(ctx, userKeys)
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
