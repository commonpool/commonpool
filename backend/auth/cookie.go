package auth

import (
	"github.com/commonpool/backend/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const accessTokenCookieName = "access_token"
const refreshTokenCookieName = "refresh_token"

// setRefreshTokenCookie sets the refresh_token cookie into the response
func setRefreshTokenCookie(c echo.Context, refreshToken string, appConfig *config.AppConfig) {
	setTokenCookie(c, refreshTokenCookieName, refreshToken, appConfig)
}

// setAccessTokenCookie sets the access_token cookie into the response
func setAccessTokenCookie(c echo.Context, accessToken string, appConfig *config.AppConfig) {
	setTokenCookie(c, accessTokenCookieName, accessToken, appConfig)
	c.Set("token", accessToken)
}

// setTokenCookie Sets the cookie into the http response
func setTokenCookie(c echo.Context, cookieName string, value string, appConfig *config.AppConfig) {
	jwtCookie := new(http.Cookie)
	jwtCookie.Name = cookieName
	jwtCookie.Value = value
	if appConfig.SecureCookies {
		jwtCookie.Secure = true
		jwtCookie.HttpOnly = true
	}
	jwtCookie.Path = "/"
	c.SetCookie(jwtCookie)
}

// clearCookies clears authentication cookies and add them to the http response
func clearCookies(c echo.Context) {
	_ = clearCookie(c, refreshTokenCookieName)
	_ = clearCookie(c, accessTokenCookieName)
}

// clearCookie utility method to clear a cookie by name
func clearCookie(c echo.Context, cookieName string) error {
	ck := new(http.Cookie)
	ck.Name = cookieName
	ck.Value = ""
	ck.MaxAge = 0
	ck.Expires = time.Unix(0, 0)
	ck.Path = "/"
	ck.HttpOnly = true
	ck.Secure = true
	c.SetCookie(ck)
	return nil
}
