package main

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
		return nil, fmt.Errorf("Error")
	}
	return []byte(value), nil
}

func setup() {
	env = map[string]string{}
	file = map[string]string{}
}

func TestGetDbConfigShouldNotFail(t *testing.T) {

	setup()

	env["DB_USER"] = "dbUser"
	env["DB_NAME"] = "dbName"
	env["DB_PASSWORD"] = "dbPassword"
	env["DB_PORT"] = "1234"
	env["DB_HOST"] = "dbHost"

	c, err := getDbConfig(getEnv, getFile)

	assert.NoError(t, err)
	assert.Equal(t, "dbUser", c.Username)
	assert.Equal(t, "dbName", c.Name)
	assert.Equal(t, "dbPassword", c.Password)
	assert.Equal(t, 1234, c.Port)
	assert.Equal(t, "dbHost", c.Host)

}

func TestDbUserFile(t *testing.T) {

	setup()

	env["DB_USER"] = "dbUser"
	env["DB_USER_FILE"] = "./db-user-file"
	env["DB_NAME"] = "dbName"
	env["DB_PASSWORD"] = "dbPassword"
	env["DB_PORT"] = "1234"
	env["DB_HOST"] = "dbHost"

	file["./db-user-file"] = "dbUserFileContent"

	c, err := getDbConfig(getEnv, getFile)

	assert.NoError(t, err)
	assert.Equal(t, "dbUserFileContent", c.Username)
	assert.Equal(t, "dbName", c.Name)
	assert.Equal(t, "dbPassword", c.Password)
	assert.Equal(t, 1234, c.Port)
	assert.Equal(t, "dbHost", c.Host)

}

func TestDbPasswordFile(t *testing.T) {

	setup()

	env["DB_USER"] = "dbUser"
	env["DB_NAME"] = "dbName"
	env["DB_PASSWORD"] = "dbPassword"
	env["DB_PASSWORD_FILE"] = "./db-password"
	env["DB_PORT"] = "1234"
	env["DB_HOST"] = "dbHost"

	file["./db-password"] = "dbPasswordFileContent"

	c, err := getDbConfig(getEnv, getFile)

	assert.NoError(t, err)
	assert.Equal(t, "dbUser", c.Username)
	assert.Equal(t, "dbName", c.Name)
	assert.Equal(t, "dbPasswordFileContent", c.Password)
	assert.Equal(t, 1234, c.Port)
	assert.Equal(t, "dbHost", c.Host)

}
