package auth

import (
	"context"
	"github.com/commonpool/backend/config"
	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/oauth2"
	"net/http"
)

const oauthCallbackPath = "/oauth2/callback"

// NewAuth Setup auth middleware
func NewAuth(e *echo.Group, appConfig *config.AppConfig, groupPrefix string, as Store) IAuth {

	ctx := context.Background()

	// Create the oidc provider
	provider, err := oidc.NewProvider(ctx, appConfig.OidcDiscoveryUrl)
	if err != nil {
		panic(err)
	}

	// Get the OAUTH config
	oauth2Config := getOauth2Config(appConfig, provider, groupPrefix)

	// Get the OIDC config
	oidcConfig := getOidcConfig(appConfig, false)
	verifier := provider.Verifier(oidcConfig)

	// Create the OidcAuthenticator object
	authz := OidcAuthenticator{
		appConfig:    appConfig,
		oauth2Config: oauth2Config,
		oidcConfig:   oidcConfig,
		oidcProvider: provider,
		verifier:     verifier,
		authStore:    as,
	}

	// Setup the router for authentication
	e.GET(oauthCallbackPath, func(c echo.Context) error {

		// decode the state
		st, err := decodeState(ctx, c.Request().URL.Query().Get("state"))
		if err != nil {
			log.Error(err)
			return err
		}

		if st.State != state {
			return c.String(http.StatusBadRequest, "State mismatch")
		}

		code := c.Request().URL.Query().Get("code")
		c.Logger().Info("code: ", code)
		oauth2token, err := authz.oauth2Config.Exchange(ctx, code)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
		}

		rawIdToken, ok := oauth2token.Extra("id_token").(string)
		if !ok {
			return c.String(http.StatusInternalServerError, "No id field in oauth2 token")
		}

		idToken, err := authz.verifier.Verify(ctx, rawIdToken)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to verify ID Token:"+err.Error())
		}

		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *JwtClaims
		}{oauth2token, new(JwtClaims)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// update cookies
		if oauth2token.AccessToken != "" {
			setAccessTokenCookie(c, oauth2token.AccessToken, appConfig)
		}
		if oauth2token.RefreshToken != "" {
			setRefreshTokenCookie(c, oauth2token.RefreshToken, appConfig)
		}

		err = saveAuthenticatedUser(c, as, resp.IDTokenClaims.Subject, resp.IDTokenClaims.Email, resp.IDTokenClaims.Email)
		if err != nil {
			return err
		}

		c.Logger().Info("desired url: ", st.DesiredUrl)
		return c.Redirect(http.StatusFound, st.DesiredUrl)

	})

	return &authz
}
