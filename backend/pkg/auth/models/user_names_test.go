package models

import (
	"github.com/commonpool/backend/pkg/keys"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestUserNames(t *testing.T) {
	userKey1 := keys.NewUserKey("user")
	userKey2 := keys.NewUserKey("user2")

	userNames := UserNames{
		userKey1: "name",
	}

	name, err := userNames.GetName(userKey1)
	assert.NoError(t, err)
	assert.Equal(t, "name", name)

	name, err = userNames.GetName(userKey2)
	assert.Error(t, err)

}
