package config

import (
	"fmt"
	"strconv"
	"strings"
)

type AppConfig struct {
	SecureCookies    bool
	BaseUri          string
	OidcDiscoveryUrl string
	OidcClientId     string
	OidcClientSecret string
	DbHost           string
	DbPort           int
	DbName           string
	DbUsername       string
	DbPassword       string
}

func GetAppConfig(readEnv EnvReader, readFile FileReader) (*AppConfig, error) {

	dbUser, hasDbUser := readEnv(dbUserEnv)
	dbUserFile, hasDbUserFile := readEnv(dbUserFileEnv)

	if !hasDbUser && !hasDbUserFile {
		panic(fmt.Errorf("%s or %s environment variable is required", dbUserFileEnv, dbUserEnv))
	}

	if hasDbUserFile {
		dbUserIo, err := readFile(dbUserFile)
		if err != nil {
			return nil, err
		}
		dbUser = string(dbUserIo)
	}

	dbPassword, hasDbPassword := readEnv(dbPasswordEnv)
	dbPasswordFile, hasDbPasswordFile := readEnv(dbPasswordFileEnv)

	if !hasDbPassword && !hasDbPasswordFile {
		return nil, fmt.Errorf("%s or %s environment variable is required", dbPasswordEnv, dbPasswordFileEnv)
	}

	if hasDbPasswordFile {
		dbUserIo, err := readFile(dbPasswordFile)
		if err != nil {
			return nil, err
		}
		dbPassword = string(dbUserIo)
	}

	dbName, ok := readEnv(dbNameEnv)
	if !ok {
		return nil, fmt.Errorf("%s environment variable is required", dbNameEnv)
	}

	dbPortStr, ok := readEnv(dbPortEnv)
	if !ok {
		return nil, fmt.Errorf("%s environment variable is required", dbPortEnv)
	}

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, err
	}

	dbHost, ok := readEnv(dbHostEnv)
	if !ok {
		return nil, fmt.Errorf("%s environment variable is required", dbHostEnv)
	}

	baseUri, hasBaseUrl := readEnv(baseUrlEnv)
	if !hasBaseUrl {
		panic(fmt.Errorf("%s env var is required", baseUrlEnv))
	}

	discoveryUrl, hasConfigUrl := readEnv(oidcDiscoveryUrlEnv)
	if !hasConfigUrl {
		panic(fmt.Errorf("%s is required", oidcDiscoveryUrlEnv))
	}

	clientIdFile, hasClientIdFile := readEnv(oidcClientIdFileEnv)
	clientId, hasClientId := readEnv(oidcClientIdEnv)
	if !hasClientIdFile && !hasClientId {
		panic(fmt.Errorf("%s or %s env var is required", oidcClientIdEnv, oidcClientIdFileEnv))
	}
	if hasClientIdFile {
		clientIdIo, err := readFile(clientIdFile)
		if err != nil {
			panic(err)
		}
		clientId = string(clientIdIo)
	}

	clientSecretFile, hasClientSecretFile := readEnv(oidcClientSecretFileEnv)
	clientSecret, hasClientSecret := readEnv(oidcClientSecretEnv)
	if !hasClientSecretFile && !hasClientSecret {
		panic(fmt.Errorf("%s or %s env var is required", oidcClientSecretEnv, oidcClientSecretFileEnv))
	}
	if hasClientSecretFile {
		clientSecretIo, err := readFile(clientSecretFile)
		if err != nil {
			panic(err)
		}
		clientSecret = string(clientSecretIo)
	}

	secureCookies := true
	secureCookiesStr, hasSecureCookies := readEnv(secureCookiesEnv)
	if hasSecureCookies {
		secureCookies = strings.ToLower(secureCookiesStr) == "true"
	}

	dbConfig := &AppConfig{
		BaseUri:          baseUri,
		OidcClientId:     clientId,
		OidcClientSecret: clientSecret,
		OidcDiscoveryUrl: discoveryUrl,
		DbHost:           dbHost,
		DbPort:           dbPort,
		DbName:           dbName,
		DbUsername:       dbUser,
		DbPassword:       dbPassword,
		SecureCookies:    secureCookies,
	}
	return dbConfig, nil
}

type EnvReader func(string) (string, bool)

const (
	dbUserEnv               = "DB_USER"
	dbUserFileEnv           = "DB_USER_FILE"
	dbPasswordEnv           = "DB_PASSWORD"
	dbPasswordFileEnv       = "DB_PASSWORD_FILE"
	dbNameEnv               = "DB_NAME"
	dbPortEnv               = "DB_PORT"
	dbHostEnv               = "DB_HOST"
	baseUrlEnv              = "BASE_URL"
	oidcDiscoveryUrlEnv     = "OIDC_DISCOVERY_URL"
	oidcClientIdFileEnv     = "OIDC_CLIENT_ID_FILE"
	oidcClientIdEnv         = "OIDC_CLIENT_ID"
	oidcClientSecretFileEnv = "OIDC_CLIENT_SECRET_FILE"
	oidcClientSecretEnv     = "OIDC_CLIENT_SECRET"
	secureCookiesEnv        = "SECURE_COOKIES"
)

type FileReader func(string) ([]byte, error)
