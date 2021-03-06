package oidc

import (
	"github.com/commonpool/backend/pkg/config"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

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
