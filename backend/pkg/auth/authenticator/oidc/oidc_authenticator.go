package oidc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/auth/store"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/utils"
	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Authenticator OIDC Implementation
type OidcAuthenticator struct {
	appConfig    *config.AppConfig
	oauth2Config oauth2.Config
	oidcConfig   *oidc.Config
	oidcProvider *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	authStore    store.Store
	userRepo     domain.UserRepository
}

func (a *OidcAuthenticator) GetLoggedInUser(ctx context.Context) (models.UserReference, error) {
	return GetLoggedInUser(ctx)
}

func (a *OidcAuthenticator) GetRedirectResponse(request *http.Request) (*authenticator.RedirectResponse, error) {
	st, err := newState(request, state)
	if err != nil {
		return nil, err
	}
	codeURL := a.oauth2Config.AuthCodeURL(st)
	response := &authenticator.RedirectResponse{
		authenticator.RedirectResponseMeta{
			RedirectTo: codeURL,
		},
	}
	return response, nil
}

func (a *OidcAuthenticator) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		SetIsAuthenticated(c, false)
		return a.RedirectToAuth(c)
	}
}

func (a *OidcAuthenticator) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		SetIsAuthenticated(c, false)
		clearCookies(c)
		v := url.Values{
			"redirect_uri": {a.appConfig.BaseUri},
		}
		discoveryUrl := a.appConfig.OidcDiscoveryUrl + "/protocol/openid-connect/logout?" + v.Encode()
		response := &authenticator.RedirectResponse{
			authenticator.RedirectResponseMeta{
				RedirectTo: discoveryUrl,
			},
		}
		return c.JSON(http.StatusOK, response)
	}
}

// Authenticate will check and user accessTokens or refreshTokens and save user info to context
func (a *OidcAuthenticator) Authenticate(redirectOnError bool) echo.MiddlewareFunc {
	ctx := context.Background()
	l := logging.WithContext(ctx)

	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			if os.Getenv("DEBUG") == "true" {
				username := c.Request().Header.Get("X-Debug-Username")
				email := c.Request().Header.Get("X-Debug-Email")
				id := c.Request().Header.Get("X-Debug-User-Id")
				isAuthenticated := c.Request().Header.Get("X-Debug-Is-Authenticated")
				if username != "" && email != "" && id != "" && isAuthenticated == "true" {
					if isAuthenticated == "true" {
						if err := SaveAuthenticatedUser(c, ctx, a.userRepo, a.authStore, id, username, email); err != nil {
							l.Error("could not save authenticated used", zap.Error(err))
							return err
						}
						return handlerFunc(c)
					}
				}
			}

			SetIsAuthenticated(c, false)

			rawAccessToken, err := a.getAccessToken(c)
			if err != nil {
				l.Error("could not get raw access token", zap.Error(err))
				return err
			}

			if rawAccessToken == "" {
				// access token not present
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.RedirectToAuth)
			}

			// verify id token
			idToken, err := a.verifier.Verify(ctx, rawAccessToken)

			if err == nil {

				// successfully got id token
				idTokenClaims := new(JwtClaims)
				err = idToken.Claims(&idTokenClaims)
				if err != nil {

					l.Error("could not get id token claims", zap.Error(err))

					clearCookies(c)
					err = fmt.Errorf("failed to retrieve id token claims %s", err)
					c.Logger().Error(err)
					return a.redirectOrNext(c, redirectOnError, handlerFunc, a.redirectToHome)
				}

				// saving into context
				err = SaveAuthenticatedUser(c, ctx, a.userRepo, a.authStore, idTokenClaims.Subject, idTokenClaims.PreferredUsername, idTokenClaims.Email)
				if err != nil {
					l.Error("could not save authenticated used", zap.Error(err))
					return err
				}

				setAccessTokenCookie(c, rawAccessToken, a.appConfig)

				return handlerFunc(c)

			}

			// from this point, access token was not valid / expired
			refreshTokenCookie := utils.FindCookie(c, refreshTokenCookieName)
			if refreshTokenCookie == nil {
				clearCookies(c)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.RedirectToAuth)
			}

			// prepare and send refresh token request
			refreshTokenFromCookie := refreshTokenCookie.Value
			formValues := url.Values{
				"client_id":     []string{a.appConfig.OidcClientId},
				"client_secret": []string{a.appConfig.OidcClientSecret},
				"grant_type":    []string{"refresh_token"},
				"refresh_token": []string{refreshTokenFromCookie},
				"scope":         []string{"openid email profile"},
			}

			// post-ing to identity server
			res, err := http.PostForm(a.oidcProvider.Endpoint().TokenURL, formValues)
			if err != nil {
				clearCookies(c)
				err = fmt.Errorf("impossible to use refresh token: %s", err)
				c.Logger().Error(err)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.RedirectToAuth)
			}

			// reading response
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				clearCookies(c)
				err = fmt.Errorf("impossible to read refresh token response: %d, %s", res.StatusCode, string(body))
				c.Logger().Error(err)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.RedirectToAuth)
			}

			// checking status code
			if res.StatusCode != http.StatusOK {
				clearCookies(c)
				err = fmt.Errorf("unexpected refresh token response code: %d, %s", res.StatusCode, string(body))
				c.Logger().Error(err)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.RedirectToAuth)
			}

			// unmarshal response
			tokenResponse := &authenticator.TokenResponse{}
			err = json.Unmarshal(body, tokenResponse)
			if err != nil {
				clearCookies(c)
				err = fmt.Errorf("impossible to unmarshal refresh token response: %s, %s", err.Error(), string(body))
				c.Logger().Error(err)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.RedirectToAuth)
			}

			// verify id token
			idToken, err = a.verifier.Verify(ctx, tokenResponse.IdToken)
			if err != nil {
				clearCookies(c)
				err = fmt.Errorf("impossible to verify refreshed token: %s", err.Error())
				c.Logger().Error(err)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.RedirectToAuth)
			}

			// retrieve claims
			idTokenClaims := new(JwtClaims)
			err = idToken.Claims(&idTokenClaims)
			if err != nil {
				clearCookies(c)
				err = fmt.Errorf("failed to verify id token: %s", err.Error())
				c.Logger().Error(err)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.redirectToHome)
			}

			// update cookies
			if tokenResponse.AccessToken != "" {
				setAccessTokenCookie(c, tokenResponse.AccessToken, a.appConfig)
			}
			if tokenResponse.RefreshToken != "" {
				setRefreshTokenCookie(c, tokenResponse.RefreshToken, a.appConfig)
			}

			err = SaveAuthenticatedUser(c, ctx, a.userRepo, a.authStore, idTokenClaims.Subject, idTokenClaims.PreferredUsername, idTokenClaims.Email)
			if err != nil {
				return err
			}

			return handlerFunc(c)

		}
	}

}

func (a *OidcAuthenticator) redirectOrNext(c echo.Context, redirectOnError bool, handlerFunc echo.HandlerFunc, redirectTo func(c echo.Context) error) error {
	if !redirectOnError {
		return handlerFunc(c)
	} else {
		return redirectTo(c)
	}
}

// getAccessToken Retrieves the access token from different sources (header, query or cookie)
func (a *OidcAuthenticator) getAccessToken(c echo.Context) (string, error) {
	var rawAccessToken = ""
	var accessTokenFromCookie = ""
	var accessTokenFromHeader = c.Request().Header.Get("Authorization")
	var accessTokenFromQuery = c.Request().URL.Query()["token"]

	accessTokenCookie, err := c.Cookie(accessTokenCookieName)
	if err == nil {
		accessTokenFromCookie = accessTokenCookie.Value
	}

	if accessTokenFromCookie != "" {
		rawAccessToken = accessTokenFromCookie
	}

	if accessTokenFromHeader != "" {
		parts := strings.Split(accessTokenFromHeader, " ")
		if len(parts) != 2 {
			err = fmt.Errorf("invalid raw access token: %s", accessTokenFromHeader)
			c.Logger().Error(err)
			c.Response().WriteHeader(400)
			return "", err
		}
		rawAccessToken = parts[1]
	}

	if len(accessTokenFromQuery) != 0 && accessTokenFromQuery[0] != "" {
		rawAccessToken = accessTokenFromQuery[0]
	}

	return rawAccessToken, nil
}

// RedirectToAuth sends redirect request to authenticate
func (a *OidcAuthenticator) RedirectToAuth(c echo.Context) error {
	st, err := newState(c.Request(), state)
	if err != nil {
		return err
	}
	codeURL := a.oauth2Config.AuthCodeURL(st)
	response := &authenticator.RedirectResponse{
		authenticator.RedirectResponseMeta{
			RedirectTo: codeURL,
		},
	}
	return c.JSON(http.StatusUnauthorized, response)
}

// redirectToHome sends redirect request to homepage
func (a *OidcAuthenticator) redirectToHome(c echo.Context) error {
	response := &authenticator.RedirectResponse{
		authenticator.RedirectResponseMeta{
			RedirectTo: "/",
		},
	}
	return c.JSON(http.StatusUnauthorized, response)
}

func (a *OidcAuthenticator) Register(c *echo.Group) {

	// Setup the router for authentication
	c.GET(oauthCallbackPath, func(c echo.Context) error {

		ctx := handler.GetContext(c)
		log := logging.WithContext(ctx)

		// decode the state
		st, err := decodeState(ctx, c.Request().URL.Query().Get("state"))
		if err != nil {
			log.Error("could not decode state", zap.Error(err))
			return err
		}

		if st.State != state {
			log.Error("state mismatch")
			return c.String(http.StatusBadRequest, "State mismatch")
		}

		rawIdToken := c.Request().URL.Query().Get("id_token")
		refreshToken := c.Request().URL.Query().Get("refresh_token")

		if rawIdToken == "" {

			code := c.Request().URL.Query().Get("code")

			oauth2token, err := a.oauth2Config.Exchange(ctx, code)
			if err != nil {
				log.Error("code exchange error", zap.Error(err))
				return c.String(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
			}

			rawIdTokenFromCode, ok := oauth2token.Extra("id_token").(string)
			if !ok {
				log.Error("no id token", zap.Error(err))
				return c.String(http.StatusInternalServerError, "No id field in oauth2 token")
			}

			rawIdToken = rawIdTokenFromCode
			refreshToken = oauth2token.RefreshToken

		}

		idToken, err := a.verifier.Verify(ctx, rawIdToken)
		if err != nil {
			log.Error("invalid id token", zap.Error(err))
			return c.String(http.StatusInternalServerError, "Failed to verify ID Token:"+err.Error())
		}

		resp := struct {
			OAuth2Token   *oidc.IDToken
			IDTokenClaims *JwtClaims
		}{idToken, new(JwtClaims)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		setAccessTokenCookie(c, rawIdToken, a.appConfig)
		setRefreshTokenCookie(c, refreshToken, a.appConfig)

		err = SaveAuthenticatedUser(c, ctx, a.userRepo, a.authStore, resp.IDTokenClaims.Subject, resp.IDTokenClaims.Email, resp.IDTokenClaims.Email)
		if err != nil {
			return err
		}

		c.Logger().Info("desired url: ", st.DesiredUrl)
		return c.Redirect(http.StatusFound, st.DesiredUrl)

	})
}