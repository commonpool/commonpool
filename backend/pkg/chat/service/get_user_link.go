package service

import (
	"fmt"
	usermodel "github.com/commonpool/backend/pkg/user/model"
)

// GetUserLink Gets the markdown representing the link to a user profile
func (c ChatService) GetUserLink(userKey usermodel.UserKey) string {
	return fmt.Sprintf("<commonpool-user id='%s'></commonpool-user>", userKey.String())
}
