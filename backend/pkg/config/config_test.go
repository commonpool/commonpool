package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	env  = map[string]string{}
	file = map[string]string{}
)

func getEnv(key string) (string, bool) {
	value, ok := env[key]
	return value, ok
}

func getFile(key string) ([]byte, error) {
	value, ok := file[key]
	if !ok {
		return nil, fmt.Errorf("error")
	}
	return []byte(value), nil
}

func setup() {
	env = map[string]string{}
	file = map[string]string{}
	env["DB_USER"] = "dbUser"
	env["DB_NAME"] = "dbName"
	env["DB_PASSWORD"] = "dbPassword"
	env["DB_PORT"] = "1234"
	env["DB_HOST"] = "dbHost"
	env["BASE_URL"] = "https://example.com"
	env["OIDC_DISCOVERY_URL"] = "https://auth.example.com"
	env["OIDC_CLIENT_ID"] = "oidc-client-id"
	env["OIDC_CLIENT_SECRET"] = "oidc-client-secret"
}

func TestGetDbConfigShouldNotFail(t *testing.T) {

	setup()

	c, err := GetAppConfig(getEnv, getFile)

	assert.NoError(t, err)
	assert.Equal(t, "dbUser", c.DbUsername)
	assert.Equal(t, "dbName", c.DbName)
	assert.Equal(t, "dbPassword", c.DbPassword)
	assert.Equal(t, 1234, c.DbPort)
	assert.Equal(t, "dbHost", c.DbHost)
	assert.Equal(t, "https://example.com", c.BaseUri)
	assert.Equal(t, "https://auth.example.com", c.OidcDiscoveryUrl)
	assert.Equal(t, "oidc-client-id", c.OidcClientId)
	assert.Equal(t, "oidc-client-secret", c.OidcClientSecret)
	assert.Equal(t, "oidc-client-secret", c.OidcClientSecret)

}

func TestDbUserFile(t *testing.T) {

	setup()

	env["DB_USER_FILE"] = "./db-user-file"
	file["./db-user-file"] = "dbUserFileContent"

	c, err := GetAppConfig(getEnv, getFile)

	assert.NoError(t, err)
	assert.Equal(t, "dbUserFileContent", c.DbUsername)

}

func TestDbPasswordFile(t *testing.T) {

	setup()

	env["DB_PASSWORD_FILE"] = "./db-password"
	file["./db-password"] = "dbPasswordFileContent"

	c, err := GetAppConfig(getEnv, getFile)

	assert.NoError(t, err)
	assert.Equal(t, "dbPasswordFileContent", c.DbPassword)
}

func TestOidcClientIdFile(t *testing.T) {

	setup()

	env["OIDC_CLIENT_ID_FILE"] = "./client-id"
	file["./client-id"] = "client-id"

	c, err := GetAppConfig(getEnv, getFile)

	assert.NoError(t, err)
	assert.Equal(t, "client-id", c.OidcClientId)
}

func TestOidcClientSecretFile(t *testing.T) {

	setup()

	env["OIDC_CLIENT_SECRET_FILE"] = "./client-secret"
	file["./client-secret"] = "client-secret"

	c, err := GetAppConfig(getEnv, getFile)

	assert.NoError(t, err)
	assert.Equal(t, "client-secret", c.OidcClientSecret)
}
