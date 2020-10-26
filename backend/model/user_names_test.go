package model

import "testing"
import "github.com/stretchr/testify/assert"

func TestUserNames(t *testing.T) {
	userKey1 := NewUserKey("user")
	userKey2 := NewUserKey("user2")

	userNames := UserNames{
		userKey1: "name",
	}

	name, err := userNames.GetName(userKey1)
	assert.NoError(t, err)
	assert.Equal(t, "name", name)

	name, err = userNames.GetName(userKey2)
	assert.Error(t, err)

}
