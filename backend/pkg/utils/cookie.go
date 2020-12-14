package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func FindCookie(c echo.Context, name string) *http.Cookie {
	for _, cookie := range c.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
