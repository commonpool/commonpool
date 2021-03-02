package auth

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

var user1 *user.User
var user2 *user.User
var user1Key keys.UserKey
var user2Key keys.UserKey

func setup() {
	user1 = &user.User{ID: "user1"}
	user2 = &user.User{ID: "user2"}
	user1Key = user1.GetUserKey()
	user2Key = user2.GetUserKey()
}

func TestNewUsersIgnoreDouble(t *testing.T) {
	setup()

	users := user.NewUsers([]*user.User{user1, user1})
	assert.Equal(t, 1, len(users.Items))
}

func TestUsersGetKey(t *testing.T) {
	setup()

	key := user1.GetUserKey()
	assert.Equal(t, "user1", key.String())
}

func TestUsersContains(t *testing.T) {
	setup()

	users := user.NewUsers([]*user.User{user1, user2})
	assert.True(t, users.Contains(user1Key))
	assert.True(t, users.Contains(user2Key))
	assert.False(t, users.Contains(keys.NewUserKey("abc")))
}

func TestUsersAppend(t *testing.T) {
	setup()

	users := user.NewUsers([]*user.User{user1})
	users = users.Append(user2)
	assert.True(t, users.Contains(user2Key))
}

func TestUsersAppendAll(t *testing.T) {
	setup()

	users := user.NewUsers([]*user.User{})
	users = users.AppendAll(user.NewUsers([]*user.User{user1, user2}))
	assert.True(t, users.Contains(user1Key))
	assert.True(t, users.Contains(user2Key))
}

func TestUsersAppendAllDouble(t *testing.T) {
	setup()

	users := user.NewUsers([]*user.User{user1})
	users = users.AppendAll(user.NewUsers([]*user.User{user1, user2}))
	assert.Equal(t, 2, len(users.Items))
}

func TestUsersGetUser(t *testing.T) {
	setup()

	users := user.NewUsers([]*user.User{user1})
	user, err := users.GetUser(user1.GetUserKey())
	assert.NoError(t, err)
	assert.Equal(t, user1.GetUserKey(), user.GetUserKey())
}

func TestUsersGetUserKeys(t *testing.T) {
	setup()

	users := user.NewUsers([]*user.User{user1, user2})
	userKeys := users.GetUserKeys()

	assert.True(t, userKeys.Contains(user1.GetUserKey()))
	assert.True(t, userKeys.Contains(user2.GetUserKey()))
}

func TestUserKeysContains(t *testing.T) {
	setup()

	userKeys := keys.NewUserKeys([]keys.UserKey{user1Key, user2Key})
	assert.True(t, userKeys.Contains(user1Key))
	assert.False(t, userKeys.Contains(keys.NewUserKey("abc")))
}

func TestNewUserKeysIgnoreDoubles(t *testing.T) {
	setup()

	userKeys := keys.NewUserKeys([]keys.UserKey{user1Key, user1Key})
	assert.True(t, userKeys.Contains(user1Key))
	assert.Equal(t, 1, len(userKeys.Items))
}
