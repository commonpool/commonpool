package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	url2 "net/url"
	"os"
)

func (h *Handler) handleLogin(c echo.Context) error {

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to generate random state: %w", err))
		return echo.ErrInternalServerError
	}

	state := base64.StdEncoding.EncodeToString(b)
	session, err := h.getSession(c)
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get session: %w", err))
		return echo.ErrInternalServerError
	}
	session.Values["state"] = state
	if err := session.Save(c.Request(), c.Response().Writer); err != nil {
		c.Logger().Error(fmt.Errorf("failed to save session: %w", err))
		return echo.ErrInternalServerError
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get authenticator: %w", err))
		return echo.ErrInternalServerError
	}

	return c.Redirect(http.StatusTemporaryRedirect, authenticator.Config.AuthCodeURL(state))

}

func (h *Handler) handleLogout(c echo.Context) error {

	session, err := h.getSession(c)
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get session: %w", err))
		return echo.ErrInternalServerError
	}

	session.Values = map[interface{}]interface{}{}
	if err := session.Save(c.Request(), c.Response().Writer); err != nil {
		c.Logger().Error(fmt.Errorf("failed to save session: %w", err))
		return echo.ErrInternalServerError
	}

	uri, err := url2.Parse(fmt.Sprintf("%s/protocol/openid-connect/logout?redirect_uri=", os.Getenv("OIDC_DISCOVERY_URL")))
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to parse oidc discovery url: %w", err))
		return echo.ErrInternalServerError
	}
	redirectUri := fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host)
	q := uri.Query()
	q.Set("redirect_uri", redirectUri)
	uri.RawQuery = q.Encode()

	uriStr := uri.String()
	return c.Redirect(http.StatusTemporaryRedirect, uriStr)

}
