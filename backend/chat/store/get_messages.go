package store

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/store"
	"go.uber.org/zap"
)

func (cs *ChatStore) GetMessages(ctx context.Context, request *chat.GetMessages) (*chat.GetMessagesResponse, error) {

	ctx, l := store.GetCtx(ctx, "ChatStore", "GetMessages")

	var messages []store.Message

	err := cs.db.
		Model(store.Message{}).
		Where("channel_id = ? AND (visible_to_user IS NULL OR visible_to_user = ?) AND sent_at < ?",
			request.Channel.String(),
			request.UserKey.String(),
			request.Before).
		Order("sent_at desc").
		Limit(request.Take + 1).
		Find(&messages).
		Error

	if err != nil {
		l.Error("could not get messages", zap.Error(err))
		return nil, err
	}

	messageCount := len(messages)
	if messageCount > 0 {
		lastMessageTs := messages[0].SentAt
		err = cs.db.Model(&chat.ChannelSubscription{}).
			Where("channel_id = ? AND user_id = ?",
				request.Channel.String(),
				request.UserKey.String(),
				lastMessageTs).
			Update("last_time_read", lastMessageTs).
			Error
		if err != nil {
			l.Error("could not update subscription", zap.Error(err))
			return nil, err
		}
	}

	if messageCount > request.Take && request.Take > 0 {
		messages = messages[:messageCount-1]
	}

	var mappedMessages []chat.Message
	for _, message := range messages {
		mappedMessage, err := mapMessage(ctx, &message)
		if err != nil {
			l.Error("could not map message", zap.Error(err))
			return nil, err
		}
		mappedMessages = append(mappedMessages, *mappedMessage)
	}

	messageLst := chat.NewMessages(mappedMessages)

	return &chat.GetMessagesResponse{
		Messages: messageLst,
		HasMore:  messageCount > request.Take,
	}, nil
}
