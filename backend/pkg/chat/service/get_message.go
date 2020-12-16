package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
)

func (c ChatService) GetMessage(ctx context.Context, messageKey chat.MessageKey) (*chat.Message, error) {
	return c.chatStore.GetMessage(ctx, messageKey)
}
