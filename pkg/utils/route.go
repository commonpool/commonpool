package utils

import "github.com/labstack/echo/v4"

func GetRouteName(c echo.Context) (string, error){
	for _, route := range c.Echo().Routes() {
		if route.Method == c.Request().Method && route.Path == c.Path() {
			return route.Name, nil
		}
	}
	return "", echo.ErrInternalServerError
}
