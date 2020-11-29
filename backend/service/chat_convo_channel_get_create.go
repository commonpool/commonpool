package service

import (
	"context"
	"errors"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"go.uber.org/zap"
)

// getOrCreateConversationChannel Will retrieve or create the conversation channel for a given set of users
// If the channel is not already created, it will automatically be created and the users will be subscribed to it.
// It also sets up RabbitMQ routing so that messages to this channel will find the user's exchange.
func (c ChatService) getOrCreateConversationChannel(ctx context.Context, userKeys *model.UserKeys) (*chat.GetOrCreateConversationChannelResponse, error) {

	ctx, l := GetCtx(ctx, "ChatService", "getOrCreateConversationChannel")

	// Retrieve the channel key for that conversation
	channelKey, err := c.GetConversationChannelKey(ctx, userKeys)
	if err != nil {
		l.Error("could not get conversation channel key", zap.Error(err))
		return nil, err
	}

	channel, err := c.cs.GetChannel(nil, channelKey)
	if err != nil && !errors.Is(err, chat.ErrChannelNotFound) {
		l.Error("could not get channel", zap.Error(err))
		return nil, err
	}

	// If the channel doesn't exist
	if errors.Is(err, chat.ErrChannelNotFound) {

		l.Debug("channel doesn't exist. Creating channel")

		var title = ""

		newChannel := &chat.Channel{
			ID:    channelKey.String(),
			Title: title,
			Type:  chat.ConversationChannel,
		}

		err := c.cs.CreateChannel(ctx, newChannel)
		if err != nil {
			l.Error("could not create channel", zap.Error(err))
			return nil, err
		}

		amqpChannelKey, err := c.GetConversationChannelKey(ctx, userKeys)
		if err != nil {
			l.Error("could not get conversation channel key")
			return nil, err
		}

		_, err = c.createSubscriptionsAndMqBindingsForUserConversation(ctx, userKeys)
		if err != nil {
			l.Error("could not create subscriptions and mq bindings", zap.Error(err))
			return nil, err
		}

		channel, err = c.cs.GetChannel(ctx, amqpChannelKey)
		if err != nil {
			l.Error("could not get channel", zap.Error(err))
			return nil, err
		}

	}

	return &chat.GetOrCreateConversationChannelResponse{
		Channel: channel,
	}, nil

}
