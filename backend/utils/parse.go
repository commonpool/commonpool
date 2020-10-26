package utils

import (
	"github.com/commonpool/backend/errors"
	"github.com/labstack/echo/v4"
	"strconv"
)

func ParseSkip(c echo.Context) (int, error) {
	skip, err := ParseQueryParamInt(c, "skip", 0)
	if err != nil {
		response := errors.ErrParseSkip(err.Error())
		return 0, &response
	}
	if skip < 0 {
		skip = 0
	}
	return skip, nil
}

func ParseTake(c echo.Context, defaultTake int, maxTake int) (int, error) {
	take, err := ParseQueryParamInt(c, "take", defaultTake)
	if err != nil {
		response := errors.ErrParseTake(err.Error())
		return 0, &response
	}
	if take < 0 {
		take = 0
	}
	if take > maxTake {
		take = maxTake
	}
	return take, nil
}

func ParseQueryParamInt(c echo.Context, paramName string, defaultValue int) (int, error) {
	paramAsStr := c.QueryParam(paramName)
	if paramAsStr != "" {
		int, err := strconv.Atoi(paramAsStr)
		if err != nil {
			response := errors.ErrCannotConvertToInt(paramAsStr, err.Error())
			return 0, &response
		}
		return int, nil
	} else {
		return defaultValue, nil
	}
}
