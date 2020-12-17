package service

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

// GetUserExchangeName will return the name of the exchange for a user key
func (c ChatService) GetUserExchangeName(ctx context.Context, userKey keys.UserKey) string {
	return "users." + userKey.String()
}
