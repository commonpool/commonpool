package handler

import (
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
)

func parseTargetFromQueryParams(c echo.Context, typeQueryParam string, valueQueryParam string) (*keys.Target, error) {
	typeParam := c.QueryParams().Get(typeQueryParam)
	if typeParam != "" {
		typeValue, err := keys.ParseOfferItemTargetType(typeParam)
		if err != nil {
			return nil, err
		}
		targetType := &typeValue
		targetIdStr := c.QueryParams().Get(valueQueryParam)
		if targetIdStr == "" {
			return nil, exceptions.ErrQueryParamRequired(valueQueryParam)
		}
		if targetType.IsGroup() {
			groupKey, err := keys.ParseGroupKey(targetIdStr)
			if err != nil {
				return nil, err
			}
			return keys.NewGroupTarget(groupKey), nil
		} else if targetType.IsUser() {
			userKey := keys.NewUserKey(targetIdStr)
			return keys.NewUserTarget(userKey), nil
		}
	}
	return nil, nil
}
