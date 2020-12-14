package service

import (
	"fmt"
	"github.com/commonpool/backend/model"
)

// GetUserLink Gets the markdown representing the link to a user profile
func (c ChatService) GetUserLink(userKey model.UserKey) string {
	return fmt.Sprintf("<commonpool-user id='%s'></commonpool-user>", userKey.String())
}
