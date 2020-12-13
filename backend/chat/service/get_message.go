package service

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
)

func (c ChatService) GetMessage(ctx context.Context, messageKey model.MessageKey) (*chat.Message, error) {
	return c.cs.GetMessage(ctx, messageKey)
}
