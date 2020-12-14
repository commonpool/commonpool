package service

import (
	"context"
	"errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/commonpool/backend/pkg/user"
	"sort"
	"strings"
)

// getOrCreateConversationChannel Will retrieve or create the conversation channel for a given set of users
// If the channel is not already created, it will automatically be created and the users will be subscribed to it.
// It also sets up RabbitMQ routing so that messages to this channel will find the user's exchange.
func (c ChatService) getOrCreateConversationChannel(ctx context.Context, userKeys *model.UserKeys) (*chat.GetOrCreateConversationChannelResponse, error) {

	channelKey, err := c.GetConversationChannelKey(ctx, userKeys)
	if err != nil {
		return nil, err
	}

	channel, err := c.chatStore.GetChannel(ctx, channelKey)
	if err != nil && !errors.Is(err, chat.ErrChannelNotFound) {
		return nil, err
	}

	if errors.Is(err, chat.ErrChannelNotFound) {

		var title = ""

		newChannel := &chat.Channel{
			Key:   channelKey,
			Title: title,
			Type:  chat.ConversationChannel,
		}

		err := c.chatStore.CreateChannel(ctx, newChannel)
		if err != nil {
			return nil, err
		}

		amqpChannelKey, err := c.GetConversationChannelKey(ctx, userKeys)
		if err != nil {
			return nil, err
		}

		_, err = c.createSubscriptionsAndMqBindingsForUserConversation(ctx, userKeys)
		if err != nil {
			return nil, err
		}

		channel, err = c.chatStore.GetChannel(ctx, amqpChannelKey)
		if err != nil {
			return nil, err
		}

	}

	return &chat.GetOrCreateConversationChannelResponse{
		Channel: channel,
	}, nil

}

func (c ChatService) createSubscriptionsAndMqBindingsForUserConversation(ctx context.Context, userKeys *model.UserKeys) ([]chat.ChannelSubscription, error) {

	channelKey, err := c.GetConversationChannelKey(ctx, userKeys)
	if err != nil {
		return nil, err
	}

	users, err := c.us.GetByKeys(ctx, userKeys.Items)
	if err != nil {
		return nil, err
	}

	var subscriptions []chat.ChannelSubscription
	for _, u := range users.Items {
		subscription, err := c.createSubscriptionAndMqBindingForUserConversation(ctx, u, users, channelKey)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, *subscription)
	}
	return subscriptions, nil
}

func (c ChatService) createSubscriptionAndMqBindingForUserConversation(
	ctx context.Context,
	user *user.User,
	conversationUsers *user.Users,
	channelKey model.ChannelKey,
) (*chat.ChannelSubscription, error) {

	conversationName := c.getConversationNameForUser(ctx, conversationUsers, user)

	channelSubscriptionKey := model.NewChannelSubscriptionKey(channelKey, user.GetUserKey())
	subscription, err := c.chatStore.CreateSubscription(ctx, channelSubscriptionKey, conversationName)
	if err != nil {
		return nil, err
	}

	amqpChannel, err := c.amqpClient.GetChannel()
	if err != nil {
		return nil, err
	}
	defer amqpChannel.Close()

	headers := c.getChannelBindingHeaders(subscription.GetKey())
	userExchangeName, err := c.CreateUserExchange(ctx, user.GetUserKey())
	err = amqpChannel.ExchangeBind(ctx, userExchangeName, "", mq.WebsocketMessagesExchange, false, headers)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// getConversationNameForUser Gets the name of the conversation for a specific user.
// A user will see the name of a conversation as equal to the names of the other participants
// but not his own name in it.
//
// Joe  would see "Dana, Mark"
// Dana would see "Joe, Mark"
// Mark would see "Dana, Joe"
func (c ChatService) getConversationNameForUser(
	ctx context.Context,
	us *user.Users,
	u *user.User,
) string {

	// First, sort the user names
	userList := us.Items

	copied := make([]*user.User, len(userList))
	copy(copied, userList)
	sort.Slice(copied, func(i, j int) bool {
		return copied[i].Username > copied[j].Username
	})

	// Get the usernames, but omit the user that will see this name
	var userNames []string
	for _, otherUser := range copied {
		if u.GetUserKey() == otherUser.GetUserKey() {
			continue
		}
		userNames = append(userNames, otherUser.Username)
	}

	// Join the name with a space
	conversationName := strings.Join(userNames, " ")

	return conversationName
}
