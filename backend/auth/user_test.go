package auth

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

var user1 exceptions.User
var user2 exceptions.User
var user1Key model.UserKey
var user2Key model.UserKey

func setup() {
	user1 = exceptions.User{ID: "user1"}
	user2 = exceptions.User{ID: "user2"}
	user1Key = user1.GetUserKey()
	user2Key = user2.GetUserKey()
}

func TestNewUsersIgnoreDouble(t *testing.T) {
	setup()

	users := exceptions.NewUsers([]User{user1, user1})
	assert.Equal(t, 1, len(users.Items))
}

func TestUsersGetKey(t *testing.T) {
	setup()

	key := user1.GetUserKey()
	assert.Equal(t, "user1", key.String())
}

func TestUsersContains(t *testing.T) {
	setup()

	users := exceptions.NewUsers([]exceptions.User{user1, user2})
	assert.True(t, users.Contains(user1Key))
	assert.True(t, users.Contains(user2Key))
	assert.False(t, users.Contains(model.NewUserKey("abc")))
}

func TestUsersAppend(t *testing.T) {
	setup()

	users := exceptions.NewUsers([]exceptions.User{user1})
	users = users.Append(user2)
	assert.True(t, users.Contains(user2Key))
}

func TestUsersAppendAll(t *testing.T) {
	setup()

	users := exceptions.NewUsers([]exceptions.User{})
	users = users.AppendAll(exceptions.NewUsers([]exceptions.User{user1, user2}))
	assert.True(t, users.Contains(user1Key))
	assert.True(t, users.Contains(user2Key))
}

func TestUsersAppendAllDouble(t *testing.T) {
	setup()

	users := exceptions.NewUsers([]exceptions.User{user1})
	users = users.AppendAll(exceptions.NewUsers([]exceptions.User{user1, user2}))
	assert.Equal(t, 2, len(users.Items))
}

func TestUsersGetUser(t *testing.T) {
	setup()

	users := exceptions.NewUsers([]exceptions.User{user1})
	user, err := users.GetUser(user1.GetUserKey())
	assert.NoError(t, err)
	assert.Equal(t, user1.GetUserKey(), user.GetUserKey())
}

func TestUsersGetUserKeys(t *testing.T) {
	setup()

	users := exceptions.NewUsers([]exceptions.User{user1, user2})
	userKeys := users.GetUserKeys()

	assert.True(t, userKeys.Contains(user1.GetUserKey()))
	assert.True(t, userKeys.Contains(user2.GetUserKey()))
}

func TestUserKeysContains(t *testing.T) {
	setup()

	userKeys := model.NewUserKeys([]model.UserKey{user1Key, user2Key})
	assert.True(t, userKeys.Contains(user1Key))
	assert.False(t, userKeys.Contains(model.NewUserKey("abc")))
}

func TestNewUserKeysIgnoreDoubles(t *testing.T) {
	setup()

	userKeys := model.NewUserKeys([]model.UserKey{user1Key, user1Key})
	assert.True(t, userKeys.Contains(user1Key))
	assert.Equal(t, 1, len(userKeys.Items))
}
