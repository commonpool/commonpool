package service

import (
	"context"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

// GetUserExchangeName will return the name of the exchange for a user key
func (c ChatService) GetUserExchangeName(ctx context.Context, userKey usermodel.UserKey) string {
	return "users." + userKey.String()
}
