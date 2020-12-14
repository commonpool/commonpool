package auth

import (
	"context"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/user"
	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"net/http"
)

const oauthCallbackPath = "/oauth2/callback"

// NewAuth Setup auth middleware
func NewAuth(e *echo.Group, appConfig *config.AppConfig, groupPrefix string, as user.Store) Authenticator {

	ctx := context.Background()
	l := logging.WithContext(ctx)

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
			l.Error("state mismatch")
			return c.String(http.StatusBadRequest, "State mismatch")
		}

		rawIdToken := c.Request().URL.Query().Get("id_token")
		refreshToken := c.Request().URL.Query().Get("refresh_token")

		if rawIdToken == "" {

			code := c.Request().URL.Query().Get("code")

			oauth2token, err := authz.oauth2Config.Exchange(ctx, code)
			if err != nil {
				l.Error("code exchange error", zap.Error(err))
				return c.String(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
			}

			rawIdTokenFromCode, ok := oauth2token.Extra("id_token").(string)
			if !ok {
				l.Error("no id token", zap.Error(err))
				return c.String(http.StatusInternalServerError, "No id field in oauth2 token")
			}

			rawIdToken = rawIdTokenFromCode
			refreshToken = oauth2token.RefreshToken

		}

		idToken, err := authz.verifier.Verify(ctx, rawIdToken)
		if err != nil {
			l.Error("invalid id token", zap.Error(err))
			return c.String(http.StatusInternalServerError, "Failed to verify ID Token:"+err.Error())
		}

		resp := struct {
			OAuth2Token   *oidc.IDToken
			IDTokenClaims *JwtClaims
		}{idToken, new(JwtClaims)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		setAccessTokenCookie(c, rawIdToken, appConfig)
		setRefreshTokenCookie(c, refreshToken, appConfig)

		err = saveAuthenticatedUser(c, as, resp.IDTokenClaims.Subject, resp.IDTokenClaims.Email, resp.IDTokenClaims.Email)
		if err != nil {
			return err
		}

		c.Logger().Info("desired url: ", st.DesiredUrl)
		return c.Redirect(http.StatusFound, st.DesiredUrl)

	})

	return &authz
}
