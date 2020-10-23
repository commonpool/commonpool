package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
	"github.com/coreos/go-oidc"
	echo "github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
	oauthCallbackPath      = "/oauth2/callback"
	SubjectKey             = "auth-subject"
	SubjectUsernameKey     = "auth-preferred-username"
	SubjectEmailKey        = "auth-email"
	IsAuthenticatedKey     = "is-authenticated"
	state                  = "someState"
)

// UserSession Holds data for the currently authenticated user
type UserSession struct {
	Username        string
	Subject         string
	Email           string
	IsAuthenticated bool
}

type Nonce struct {
	DesiredUrl string `json:"des"`
	State      string `json:"state"`
}

// JwtClaims Claims part of the oidc response
type JwtClaims struct {
	Issuer  string `json:"iss"`
	Subject string `json:"sub"`
	// Audience          []string `json:"aud"`
	Expiration        int64  `json:"exp"`
	IssuedAt          int64  `json:"iat"`
	Id                string `json:"jti"`
	Type              string `json:"typ"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	EmailVerified     bool   `json:"email_verified"`
}

// TokenResponse response from the oidc provider
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	IdToken          string `json:"id_token"`
	NotBeforePolicy  int64  `json:"notBeforePolicy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

// Interface for authorization module
type IAuth interface {
	Authenticate(redirectOnError bool) echo.MiddlewareFunc
	GetAuthUserSession(c echo.Context) UserSession
}

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
	var isAuthenticated bool = c.Get(IsAuthenticatedKey).(bool)
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

// Ascertain that OidcAuthenticator implements IAuth
var _ IAuth = &OidcAuthenticator{}

// RedirectResponseMeta API response when redirecting
type RedirectResponseMeta struct {
	RedirectTo string `json:"redirectTo"`
}

// RedirectResponse API response container
type RedirectResponse struct {
	Meta RedirectResponseMeta `json:"meta"`
}

// Authenticate will check and user accessTokens or refreshTokens and save user info to context
func (a *OidcAuthenticator) Authenticate(redirectOnError bool) echo.MiddlewareFunc {
	ctx := context.Background()

	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			SetIsAuthenticated(c, false)

			rawAccessToken, err := a.getAccessToken(c)
			if err != nil {
				return err
			}

			if rawAccessToken == "" {
				// access token not present
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.redirectToAuth)
			}

			// verify id token
			idToken, err := a.verifier.Verify(ctx, rawAccessToken)

			if err == nil {

				// successfully got id token
				idTokenClaims := new(JwtClaims)
				err = idToken.Claims(&idTokenClaims)
				if err != nil {
					clearCookies(c)
					err = fmt.Errorf("failed to retrieve id token claims %s", err)
					c.Logger().Error(err)
					return a.redirectOrNext(c, redirectOnError, handlerFunc, a.redirectToHome)
				}

				// saving into context
				err = saveAuthenticatedUser(c, a.authStore, idTokenClaims.Subject, idTokenClaims.PreferredUsername, idTokenClaims.Email)
				if err != nil {
					return err
				}

				return handlerFunc(c)

			}

			// from this point, access token was not valid / expired
			refreshTokenCookie := utils.FindCookie(c, refreshTokenCookieName)
			if refreshTokenCookie == nil {
				clearCookies(c)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.redirectToAuth)
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
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.redirectToAuth)
			}

			// reading response
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				clearCookies(c)
				err = fmt.Errorf("impossible to read refresh token response: %d, %s", res.StatusCode, string(body))
				c.Logger().Error(err)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.redirectToAuth)
			}

			// checking status code
			if res.StatusCode != http.StatusOK {
				clearCookies(c)
				err = fmt.Errorf("unexpected refresh token response code: %d, %s", res.StatusCode, string(body))
				c.Logger().Error(err)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.redirectToAuth)
			}

			// unmarshal response
			tokenResponse := &TokenResponse{}
			err = json.Unmarshal(body, tokenResponse)
			if err != nil {
				clearCookies(c)
				err = fmt.Errorf("impossible to unmarshal refresh token response: %s, %s", err.Error(), string(body))
				c.Logger().Error(err)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.redirectToAuth)
			}

			// verify id token
			idToken, err = a.verifier.Verify(ctx, tokenResponse.IdToken)
			if err != nil {
				clearCookies(c)
				err = fmt.Errorf("impossible to verify refreshed token: %s", err.Error())
				c.Logger().Error(err)
				return a.redirectOrNext(c, redirectOnError, handlerFunc, a.redirectToAuth)
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

func saveAuthenticatedUser(c echo.Context, store Store, sub string, username string, email string) error {
	SetIsAuthenticated(c, true)
	setSubject(c, sub)
	setUsername(c, username)
	setEmail(c, email)
	SetIsAuthenticated(c, true)
	return saveUserInfo(store, sub, email, username)
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
		rawAccessToken = accessTokenFromHeader
	}

	if accessTokenFromQuery != "" {
		c.Logger().Info("Using access token from query")
		rawAccessToken = accessTokenFromQuery
	}
	return rawAccessToken, nil
}

func newState(c echo.Context, state string) (string, error) {

	origin := c.Request().Header.Get("Origin")
	referrer := c.Request().Referer()

	c.Logger().Info("Encoding nonce with url, origin: ", origin, " referrer: ", referrer)
	st := Nonce{
		DesiredUrl: referrer,
		State:      state,
	}
	bytes, err := json.Marshal(st)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(bytes)
	return b64, nil
}

func decodeState(state string) (*Nonce, error) {
	bytes, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		return nil, err
	}
	nonce := &Nonce{}
	err = json.Unmarshal(bytes, nonce)
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

// redirectToAuth sends redirect request to authenticate
func (a *OidcAuthenticator) redirectToAuth(c echo.Context) error {
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

// NewAuth Setup auth middleware
func NewAuth(e *echo.Group, appConfig *config.AppConfig, groupPrefix string, as Store) IAuth {

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, appConfig.OidcDiscoveryUrl)
	if err != nil {
		panic(err)
	}
	oauth2Config := getOauth2Config(appConfig, provider, groupPrefix)
	oidcConfig := getOidcConfig(appConfig, false)
	verifier := provider.Verifier(oidcConfig)

	authz := OidcAuthenticator{
		appConfig:    appConfig,
		oauth2Config: oauth2Config,
		oidcConfig:   oidcConfig,
		oidcProvider: provider,
		verifier:     verifier,
		authStore:    as,
	}

	e.GET(oauthCallbackPath, func(c echo.Context) error {

		st, err := decodeState(c.Request().URL.Query().Get("state"))
		if err != nil {
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

// saveUserInfo saves the logged in user info th database
func saveUserInfo(as Store, sub string, email string, username string) error {
	return as.Upsert(model.NewUserKey(sub), email, username)
}

// getOidcConfig gets the OIDC provider config
func getOidcConfig(appConfig *config.AppConfig, skipClientIdCheck bool) *oidc.Config {
	oidcConfig := &oidc.Config{
		ClientID:          appConfig.OidcClientId,
		SkipClientIDCheck: skipClientIdCheck,
	}
	return oidcConfig
}

// getOauth2Config gets the OAUTH2 provider config
func getOauth2Config(appConfig *config.AppConfig, provider *oidc.Provider, pathPrefix string) oauth2.Config {
	oauth2Config := oauth2.Config{
		ClientID:     appConfig.OidcClientId,
		ClientSecret: appConfig.OidcClientSecret,
		RedirectURL:  appConfig.BaseUri + pathPrefix + oauthCallbackPath,
		Endpoint:     provider.Endpoint(),
		Scopes: []string{
			oidc.ScopeOpenID,
			"profile",
			"email",
		},
	}
	return oauth2Config
}

// setRefreshTokenCookie sets the refresh_token cookie into the response
func setRefreshTokenCookie(c echo.Context, refreshToken string, appConfig *config.AppConfig) {
	setTokenCookie(c, refreshTokenCookieName, refreshToken, appConfig)
}

// setAccessTokenCookie sets the access_token cookie into the response
func setAccessTokenCookie(c echo.Context, accessToken string, appConfig *config.AppConfig) {
	setTokenCookie(c, accessTokenCookieName, accessToken, appConfig)
}

func setTokenCookie(c echo.Context, cookieName string, value string, appConfig *config.AppConfig) {

	cookie, err := c.Cookie(cookieName)
	if err == nil {
		if cookie.Value == value {
			return
		}
	}

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

func clearCookies(c echo.Context) {
	_ = clearCookie(c, refreshTokenCookieName)
	_ = clearCookie(c, accessTokenCookieName)
}

func clearCookie(c echo.Context, cookieName string) error {
	refreshTokenCookie, err := c.Cookie(cookieName)
	if err != nil {
		return err
	}
	refreshTokenCookie.MaxAge = 0
	refreshTokenCookie.Expires = time.Unix(0, 0)
	c.SetCookie(refreshTokenCookie)
	return nil
}

func setSubject(e echo.Context, subject string) {
	e.Set(SubjectKey, subject)
}
func setUsername(e echo.Context, username string) {
	e.Set(SubjectUsernameKey, username)
}
func setEmail(e echo.Context, email string) {
	e.Set(SubjectEmailKey, email)
}
func SetIsAuthenticated(e echo.Context, isAuthenticated bool) {
	e.Set(IsAuthenticatedKey, isAuthenticated)
}

type MockAuthorizer struct {
	IsAuthorized       bool
	MockCurrentSession func() UserSession
}

func (a *MockAuthorizer) GetAuthUserSession(c echo.Context) UserSession {
	return a.MockCurrentSession()
}

var _ IAuth = &MockAuthorizer{}

func (a *MockAuthorizer) Authenticate(redirectOnError bool) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return handlerFunc(c)
		}
	}
}

func NewTestAuthorizer() *MockAuthorizer {
	return &MockAuthorizer{}
}
