package utils

import (
	"fmt"
	"github.com/commonpool/backend/errors"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
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

func ParseBefore(c echo.Context) (*time.Time, error) {
	before, err := ParseQueryParamTimestamp(c, "before")
	if err != nil {
		response := errors.ErrParseBefore(err.Error())
		return nil, &response
	}
	return before, nil
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
		intValue, err := strconv.Atoi(paramAsStr)
		if err != nil {
			response := errors.ErrCannotConvertToInt(paramAsStr, err.Error())
			return 0, &response
		}
		return intValue, nil
	} else {
		return defaultValue, nil
	}
}

func ParseQueryParamTimestamp(c echo.Context, paramName string) (*time.Time, error) {
	paramAsStr := c.QueryParam(paramName)
	if paramAsStr == "" {
		return nil, fmt.Errorf("query param " + paramName + " is required")
	}
	i, err := strconv.ParseInt(paramAsStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("query param " + paramName + " could not be converted to int64")
	}
	tm := time.Unix(i, 0)
	return &tm, nil
}
