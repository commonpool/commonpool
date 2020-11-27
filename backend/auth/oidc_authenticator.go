package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// IAuth OIDC Implementation
type OidcAuthenticator struct {
	appConfig    *config.AppConfig
	oauth2Config oauth2.Config
	oidcConfig   *oidc.Config
	oidcProvider *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	authStore    Store
}

// GetAuthenticatedUser gets the current authenticated user
func (a *OidcAuthenticator) GetAuthUserSession(c echo.Context) UserSession {
	var isAuthenticated = c.Get(IsAuthenticatedKey).(bool)
	if !isAuthenticated {
		return UserSession{
			Username:        "",
			Subject:         "",
			Email:           "",
			IsAuthenticated: false,
		}
	}
	return UserSession{
		Username:        c.Get(SubjectUsernameKey).(string),
		Subject:         c.Get(SubjectKey).(string),
		Email:           c.Get(SubjectEmailKey).(string),
		IsAuthenticated: isAuthenticated,
	}
}

// GetAuthenticatedUser gets the current authenticated user
func (a *OidcAuthenticator) GetAuthUserSession2(c context.Context) UserSession {

	var isAuthenticatedIntf = c.Value(IsAuthenticatedKey)

	if isAuthenticated, ok := isAuthenticatedIntf.(bool); !ok || !isAuthenticated {
		return UserSession{
			Username:        "",
			Subject:         "",
			Email:           "",
			IsAuthenticated: false,
		}
	}

	return UserSession{
		Username:        c.Value(SubjectUsernameKey).(string),
		Subject:         c.Value(SubjectKey).(string),
		Email:           c.Value(SubjectEmailKey).(string),
		IsAuthenticated: true,
	}
}

func (a *OidcAuthenticator) GetAuthUserKey(c echo.Context) model.UserKey {
	return model.NewUserKey(c.Get(SubjectKey).(string))
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
		url := a.appConfig.OidcDiscoveryUrl + "/protocol/openid-connect/logout?" + v.Encode()
		response := &RedirectResponse{
			RedirectResponseMeta{
				RedirectTo: url,
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

			SetIsAuthenticated(c, false)

			rawAccessToken, err := a.getAccessToken(c)
			if err != nil {
				l.Error("could not get raw access token", zap.Error(err))
				return err
			}

			if rawAccessToken == "" {
				l.Warn("access token not present")
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
				err = saveAuthenticatedUser(c, a.authStore, idTokenClaims.Subject, idTokenClaims.PreferredUsername, idTokenClaims.Email)
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
			tokenResponse := &TokenResponse{}
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

			err = saveAuthenticatedUser(c, a.authStore, idTokenClaims.Subject, idTokenClaims.PreferredUsername, idTokenClaims.Email)
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
	var accessTokenFromQuery = c.Request().URL.Query().Get("token")

	accessTokenCookie, err := c.Cookie(accessTokenCookieName)
	if err == nil {
		accessTokenFromCookie = accessTokenCookie.Value
	}

	if accessTokenFromCookie != "" {
		c.Logger().Info("Using access token from cookie")
		rawAccessToken = accessTokenFromCookie
	}

	if accessTokenFromHeader != "" {
		c.Logger().Info("Using access token from header")
		parts := strings.Split(accessTokenFromHeader, " ")
		if len(parts) != 2 {
			err = fmt.Errorf("invalid raw access token: %s", accessTokenFromHeader)
			c.Logger().Error(err)
			c.Response().WriteHeader(400)
			return "", err
		}
		rawAccessToken = parts[1]
	}

	if accessTokenFromQuery != "" {
		c.Logger().Info("Using access token from query")
		rawAccessToken = accessTokenFromQuery
	}

	return rawAccessToken, nil
}

// RedirectToAuth sends redirect request to authenticate
func (a *OidcAuthenticator) RedirectToAuth(c echo.Context) error {
	st, err := newState(c, state)
	if err != nil {
		return err
	}
	codeURL := a.oauth2Config.AuthCodeURL(st)
	response := &RedirectResponse{
		RedirectResponseMeta{
			RedirectTo: codeURL,
		},
	}
	return c.JSON(http.StatusUnauthorized, response)
}

// redirectToHome sends redirect request to homepage
func (a *OidcAuthenticator) redirectToHome(c echo.Context) error {
	response := &RedirectResponse{
		RedirectResponseMeta{
			RedirectTo: "/",
		},
	}
	return c.JSON(http.StatusUnauthorized, response)
}
