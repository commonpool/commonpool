package handler

import (
	"context"
	"cp/pkg/api"
	"fmt"
	oidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"os"
)

func (h *Handler) getSession(c echo.Context) (*sessions.Session, error) {
	return h.cookieStore.Get(c.Request(), "session")
}

func GetSession(store sessions.Store, c echo.Context) (*sessions.Session, error) {
	return store.Get(c.Request(), "session")
}

func GetProfile(store sessions.Store, c echo.Context) (*api.Profile, error) {
	session, err := GetSession(store, c)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	idIntf, hasID := session.Values["id"]
	if !hasID {
		return nil, nil
	}

	emailIntf, hasEmail := session.Values["email"]
	if !hasEmail {
		return nil, nil
	}

	usernameIntf, hasUsername := session.Values["username"]
	if !hasUsername {
		return nil, nil
	}

	groupsIntf, hasGroups := session.Values["groups"]
	if !hasGroups {
		return nil, nil
	}

	id, ok := idIntf.(string)
	if !ok {
		return nil, nil
	}

	email, ok := emailIntf.(string)
	if !ok {
		return nil, nil
	}

	username, ok := usernameIntf.(string)
	if !ok {
		return nil, nil
	}

	groups, ok := groupsIntf.([]string)
	if !ok {
		return nil, nil
	}

	var profile = &api.Profile{
		Email:    email,
		Username: username,
		ID:       id,
		Groups:   groups,
	}
	return profile, nil
}

type Authenticator struct {
	provider *oidc.Provider
	Config   oauth2.Config
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, os.Getenv("OIDC_DISCOVERY_URL"))
	if err != nil {
		err := fmt.Errorf("failed to get oidc provider: %w", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("OIDC_CLIENT_ID"),
		ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OIDC_REDIRECT_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Authenticator{
		provider: provider,
		Config:   conf,
	}, nil
}
