package service

import (
	"context"
	"github.com/commonpool/backend/model"
)

// GetUserExchangeName will return the name of the exchange for a user key
func (c ChatService) GetUserExchangeName(ctx context.Context, userKey model.UserKey) string {
	return "users." + userKey.String()
}
