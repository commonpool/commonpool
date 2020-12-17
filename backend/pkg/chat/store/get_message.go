package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
)

func (cs *ChatStore) GetMessage(ctx context.Context, messageKey keys.MessageKey) (*chat.Message, error) {

	var message Message
	err := cs.db.
		Model(Message{}).
		Where("id = ?", messageKey.String()).
		First(&message).
		Error

	if err != nil {
		return nil, err
	}

	returnMessage, err := mapMessage(ctx, &message)
	if err != nil {
		return nil, err
	}

	return returnMessage, nil
}
