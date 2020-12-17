package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
)

func (c ChatService) GetMessage(ctx context.Context, messageKey keys.MessageKey) (*chat.Message, error) {
	return c.chatStore.GetMessage(ctx, messageKey)
}
