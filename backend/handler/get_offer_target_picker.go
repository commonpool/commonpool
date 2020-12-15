package handler

import (
	"github.com/commonpool/backend/model"
	group2 "github.com/commonpool/backend/pkg/group"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/handler"
	model2 "github.com/commonpool/backend/pkg/trading/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) HandleOfferItemTargetPicker(c echo.Context) error {

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

	var userKeys []usermodel.UserKey
	var groupKeys []groupmodel.GroupKey
	for _, target := range targets.Items {
		if target.IsForUser() {
			userKeys = append(userKeys, target.GetUserKey())
		} else if target.IsForGroup() {
			groupKeys = append(groupKeys, target.GetGroupKey())
		}
	}

	for _, grpKey := range groupKeys {
		grp, err := h.groupService.GetGroup(ctx, &group2.GetGroupRequest{
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
