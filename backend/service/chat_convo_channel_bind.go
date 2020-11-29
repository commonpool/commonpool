package service

import (
	"context"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"go.uber.org/zap"
)

// createSubscriptionsAndMqBindingsForUserConversation subscribes multiple users to a conversation channel
// It also generates the RabbitMQ bindings so that messages to the conversation
// are also properly routed to the user's exchange.
func (c ChatService) createSubscriptionsAndMqBindingsForUserConversation(ctx context.Context, userKeys *model.UserKeys) ([]chat.ChannelSubscription, error) {

	ctx, l := GetCtx(ctx, "ChatService", "createSubscriptionsAndMqBindingsForUserConversation")

	channelKey, err := c.GetConversationChannelKey(ctx, userKeys)
	if err != nil {
		l.Error("could not get conversation channel key")
		return nil, err
	}

	users, err := c.us.GetByKeys(ctx, userKeys.Items)
	if err != nil {
		l.Error("could not get users by keys")
		return nil, err
	}

	var subscriptions []chat.ChannelSubscription
	for _, user := range users.Items {

		l.Info("user", zap.String("user_id", user.ID))

		subscription, err := c.createSubscriptionAndMqBindingForUserConversation(ctx, user, users, channelKey)
		if err != nil {
			l.Error("could not create subscription and mq binding for user conversation", zap.Error(err))
			return nil, err
		}
		subscriptions = append(subscriptions, *subscription)
	}
	return subscriptions, nil
}

// createSubscriptionAndMqBindingForUserConversation subscribes the user to a conversation channel.
// It also generates the RabbitMQ bindings so that messages to the conversation
// are also properly routed to the user's exchange.
func (c ChatService) createSubscriptionAndMqBindingForUserConversation(
	ctx context.Context,
	user auth.User,
	conversationUsers auth.Users,
	channelKey model.ChannelKey,
) (*chat.ChannelSubscription, error) {

	ctx, l := GetCtx(ctx, "ChatService", "createSubscriptionAndMqBindingForUserConversation")

	// The name of the conversation
	conversationName := c.getConversationNameForUser(ctx, conversationUsers, user)

	channelSubscriptionKey := model.NewChannelSubscriptionKey(channelKey, user.GetUserKey())
	subscription, err := c.cs.CreateSubscription(ctx, channelSubscriptionKey, conversationName)
	if err != nil {
		l.Error("could not create subscription")
		return nil, err
	}

	amqpChannel, err := c.mq.GetChannel()
	if err != nil {
		l.Error("could not get amqp channel", zap.Error(err))
		return nil, err
	}
	defer amqpChannel.Close()

	headers := c.getChannelBindingHeaders(subscription.GetKey())
	userExchangeName, err := c.CreateUserExchange(ctx, user.GetUserKey())
	err = amqpChannel.ExchangeBind(ctx, userExchangeName, "", amqp.WebsocketMessagesExchange, false, headers)
	if err != nil {
		l.Error("could not register user channel binding", zap.Error(err))
		return nil, err
	}

	return subscription, nil
}
