package config

import (
	"fmt"
	"strconv"
)

type AppConfig struct {
	SecureCookies      bool
	BaseUri            string
	OidcDiscoveryUrl   string
	OidcClientId       string
	OidcClientSecret   string
	DbHost             string
	DbPort             int
	DbName             string
	DbUsername         string
	DbPassword         string
	CallbackToken      string
	AmqpUrl            string
	RedisHost          string
	RedisPort          string
	RedisPassword      string
	RedisTlsEnabled    bool
	RedisTlsSkipVerify bool
}

func GetAppConfig(readEnv EnvReader, readFile FileReader) (*AppConfig, error) {

	var (
		err                error
		dbUser             string
		dbPassword         string
		dbName             string
		dbPort             string
		callbackToken      string
		amqpUrl            string
		dbHost             string
		baseUri            string
		discoveryUrl       string
		clientId           string
		clientSecret       string
		secureCookies      string
		redisHost          string
		redisPort          string
		redisPassword      string
		redisTlsEnabled    string
		redisTlsSkipVerify string
	)

	if dbName, err = readEnvVarOrFile(readFile, readEnv, dbNameEnv); err != nil {
		return nil, err
	}

	if dbPort, err = readEnvVarOrFile(readFile, readEnv, dbPortEnv); err != nil {
		return nil, err
	}

	if dbHost, err = readEnvVarOrFile(readFile, readEnv, dbHostEnv); err != nil {
		return nil, err
	}

	if dbUser, err = readEnvVarOrFile(readFile, readEnv, dbUserEnv); err != nil {
		return nil, err
	}

	if dbPassword, err = readEnvVarOrFile(readFile, readEnv, dbPasswordEnv); err != nil {
		return nil, err
	}

	if callbackToken, err = readEnvVarOrFile(readFile, readEnv, callbackTokenEnv); err != nil {
		return nil, err
	}

	if amqpUrl, err = readEnvVarOrFile(readFile, readEnv, amqpUrlEnv); err != nil {
		return nil, err
	}

	if baseUri, err = readEnvVarOrFile(readFile, readEnv, baseUrlEnv); err != nil {
		return nil, err
	}

	if discoveryUrl, err = readEnvVarOrFile(readFile, readEnv, oidcDiscoveryUrlEnv); err != nil {
		return nil, err
	}

	if clientId, err = readEnvVarOrFile(readFile, readEnv, oidcClientIdEnv); err != nil {
		return nil, err
	}

	if clientSecret, err = readEnvVarOrFile(readFile, readEnv, oidcClientSecretEnv); err != nil {
		return nil, err
	}

	if secureCookies, err = readEnvVarOrFile(readFile, readEnv, secureCookiesEnv); err != nil {
		return nil, err
	}

	if redisHost, err = readEnvVarOrFile(readFile, readEnv, redisHostEnv); err != nil {
		return nil, err
	}
	if redisPort, err = readEnvVarOrFile(readFile, readEnv, redisPortEnv); err != nil {
		return nil, err
	}
	if redisPassword, err = readEnvVarOrFile(readFile, readEnv, redisPasswordEnv); err != nil {
		return nil, err
	}
	if redisTlsEnabled, err = readEnvVarOrFile(readFile, readEnv, redisTlsEnabledEnv); err != nil {
		return nil, err
	}
	if redisTlsSkipVerify, err = readEnvVarOrFile(readFile, readEnv, redisTlsSkipVerifyEnv); err != nil {
		return nil, err
	}

	dbPortValue, err := strconv.Atoi(dbPort)
	if err != nil {
		return nil, err
	}

	appConfig := &AppConfig{
		BaseUri:            baseUri,
		OidcClientId:       clientId,
		OidcClientSecret:   clientSecret,
		OidcDiscoveryUrl:   discoveryUrl,
		DbHost:             dbHost,
		DbPort:             dbPortValue,
		DbName:             dbName,
		DbUsername:         dbUser,
		DbPassword:         dbPassword,
		SecureCookies:      secureCookies == "true",
		CallbackToken:      callbackToken,
		AmqpUrl:            amqpUrl,
		RedisHost:          redisHost,
		RedisPort:          redisPort,
		RedisPassword:      redisPassword,
		RedisTlsEnabled:    redisTlsEnabled == "true",
		RedisTlsSkipVerify: redisTlsSkipVerify == "true",
	}
	return appConfig, nil
}

func readFileValue(readFile FileReader, filePath string) (string, error) {
	fileBytes, err := readFile(filePath)
	if err != nil {
		return "", err
	}
	fileValue := string(fileBytes)
	return fileValue, nil
}

func readEnvVarOrFile(readFile FileReader, readEnv EnvReader, envValueName string) (string, error) {
	var err error
	value, hasValue := readEnv(envValueName)
	file, hasFile := readEnv(envValueName + "_FILE")
	if !hasValue && !hasFile {
		return "", fmt.Errorf("%s or %s environment variable is required", envValueName, envValueName+"_FILE")
	}
	if hasFile {
		if value, err = readFileValue(readFile, file); err != nil {
			return "", err
		}
	}
	return value, nil
}

type EnvReader func(string) (string, bool)

const (
	dbUserEnv             = "DB_USER"
	dbPasswordEnv         = "DB_PASSWORD"
	dbNameEnv             = "DB_NAME"
	dbPortEnv             = "DB_PORT"
	dbHostEnv             = "DB_HOST"
	baseUrlEnv            = "BASE_URL"
	oidcDiscoveryUrlEnv   = "OIDC_DISCOVERY_URL"
	oidcClientIdEnv       = "OIDC_CLIENT_ID"
	oidcClientSecretEnv   = "OIDC_CLIENT_SECRET"
	secureCookiesEnv      = "SECURE_COOKIES"
	redisHostEnv          = "REDIS_HOST"
	redisPortEnv          = "REDIS_PORT"
	redisPasswordEnv      = "REDIS_PASSWORD"
	redisTlsEnabledEnv    = "REDIS_ENABLE_TLS"
	redisTlsSkipVerifyEnv = "REDIS_TLS_SKIP_VERIFY"
	callbackTokenEnv      = "CALLBACK_TOKEN"
	amqpUrlEnv            = "AMQP_URL"
)

type FileReader func(string) ([]byte, error)
