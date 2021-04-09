package handler

import (
	"context"
	"cp/pkg/api"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func (h *Handler) handleOauthCallback(c echo.Context) error {
	session, err := h.getSession(c)
	if err != nil {
		return err
	}

	if c.Request().URL.Query().Get("state") != session.Values["state"] {
		return fmt.Errorf("invalid state parameter")
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		err := fmt.Errorf("could not get oidc authenticator: %w", err)
		return err
	}

	token, err := authenticator.Config.Exchange(context.TODO(), c.Request().URL.Query().Get("code"))
	if err != nil {
		c.Logger().Error(fmt.Errorf("could not exchange token: %w", err))
		return echo.ErrUnauthorized
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.Logger().Error(fmt.Errorf("no id_token found in oauth2 token."))
		return echo.ErrInternalServerError
	}

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("OIDC_CLIENT_ID"),
	}

	idToken, err := authenticator.provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to verify id_token: %w", err))
		return echo.ErrInternalServerError
	}

	var profile api.Profile
	if err := idToken.Claims(&profile); err != nil {
		c.Logger().Error(fmt.Errorf("failed to retrieve profile from id_token: %w", err))
		return echo.ErrInternalServerError
	}

	session.Values["id_token"] = rawIDToken
	session.Values["email"] = profile.Email
	session.Values["username"] = profile.Username
	session.Values["id"] = profile.ID
	session.Values["groups"] = profile.Groups

	if err := session.Save(c.Request(), c.Response().Writer); err != nil {
		c.Logger().Error(fmt.Errorf("failed to save session: %w", err))
		return echo.ErrInternalServerError
	}

	return c.Redirect(http.StatusSeeOther, "/")

}
