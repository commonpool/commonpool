package service

import (
	"context"
	chatmodel "github.com/commonpool/backend/pkg/chat/model"
)

func (c ChatService) GetMessage(ctx context.Context, messageKey chatmodel.MessageKey) (*chatmodel.Message, error) {
	return c.chatStore.GetMessage(ctx, messageKey)
}