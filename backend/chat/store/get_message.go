package store

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/store"
	"go.uber.org/zap"
)

func (cs *ChatStore) GetMessage(ctx context.Context, messageKey model.MessageKey) (*chat.Message, error) {

	ctx, l := store.GetCtx(ctx, "ChatStore", "GetMessage")

	var message store.Message
	err := cs.db.
		Model(store.Message{}).
		Where("id = ?", messageKey.String()).
		First(&message).
		Error

	if err != nil {
		l.Error("could not get message", zap.Error(err))
		return nil, err
	}

	returnMessage, err := mapMessage(ctx, &message)
	if err != nil {
		l.Error("could not map message", zap.Error(err))
		return nil, err
	}

	return returnMessage, nil
}
