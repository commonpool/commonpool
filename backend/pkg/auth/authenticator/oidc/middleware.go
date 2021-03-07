package oidc

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/auth/store"
	"github.com/commonpool/backend/pkg/config"
	"github.com/coreos/go-oidc"
)

const oauthCallbackPath = "/oauth2/callback"

// NewAuth Setup auth middleware
func NewAuth(appConfig *config.AppConfig, groupPrefix string, as store.Store, userRepo domain.UserRepository) *OidcAuthenticator {

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
		userRepo:     userRepo,
	}

	return &authz
}
