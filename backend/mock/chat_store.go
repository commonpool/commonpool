// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/model"
	"sync"
)

// Ensure, that ChatStore does implement chat.Store.
// If this is not the case, regenerate this file with moq.
var _ chat.Store = &ChatStore{}

// ChatStore is a mock implementation of chat.Store.
//
//     func TestSomethingThatUsesStore(t *testing.T) {
//
//         // make and configure a mocked chat.Store
//         mockedStore := &ChatStore{
//             CreateChannelFunc: func(ctx context.Context, channel *model.Channel) error {
// 	               panic("mock out the CreateChannel method")
//             },
//             CreateSubscriptionFunc: func(ctx context.Context, key model.ChannelSubscriptionKey, name string) (*model.ChannelSubscription, error) {
// 	               panic("mock out the CreateSubscription method")
//             },
//             DeleteSubscriptionFunc: func(ctx context.Context, key model.ChannelSubscriptionKey) error {
// 	               panic("mock out the DeleteSubscription method")
//             },
//             GetChannelFunc: func(ctx context.Context, channelKey model.ChannelKey) (*model.Channel, error) {
// 	               panic("mock out the GetChannel method")
//             },
//             GetMessageFunc: func(ctx context.Context, messageKey model.MessageKey) (*model.Message, error) {
// 	               panic("mock out the GetMessage method")
//             },
//             GetMessagesFunc: func(ctx context.Context, request *chat.GetMessages) (*chat.GetMessagesResponse, error) {
// 	               panic("mock out the GetMessages method")
//             },
//             GetSubscriptionFunc: func(ctx context.Context, request *chat.GetSubscription) (*model.ChannelSubscription, error) {
// 	               panic("mock out the GetSubscription method")
//             },
//             GetSubscriptionsForChannelFunc: func(ctx context.Context, channelKey model.ChannelKey) ([]*model.ChannelSubscription, error) {
// 	               panic("mock out the GetSubscriptionsForChannel method")
//             },
//             GetSubscriptionsForUserFunc: func(ctx context.Context, request *chat.GetSubscriptions) (*model.ChannelSubscriptions, error) {
// 	               panic("mock out the GetSubscriptionsForUser method")
//             },
//             SaveMessageFunc: func(ctx context.Context, request *model.Message) error {
// 	               panic("mock out the SaveMessage method")
//             },
//         }
//
//         // use mockedStore in code that requires chat.Store
//         // and then make assertions.
//
//     }
type ChatStore struct {
	// CreateChannelFunc mocks the CreateChannel method.
	CreateChannelFunc func(ctx context.Context, channel *model.Channel) error

	// CreateSubscriptionFunc mocks the CreateSubscription method.
	CreateSubscriptionFunc func(ctx context.Context, key model.ChannelSubscriptionKey, name string) (*model.ChannelSubscription, error)

	// DeleteSubscriptionFunc mocks the DeleteSubscription method.
	DeleteSubscriptionFunc func(ctx context.Context, key model.ChannelSubscriptionKey) error

	// GetChannelFunc mocks the GetChannel method.
	GetChannelFunc func(ctx context.Context, channelKey model.ChannelKey) (*model.Channel, error)

	// GetMessageFunc mocks the GetMessage method.
	GetMessageFunc func(ctx context.Context, messageKey model.MessageKey) (*model.Message, error)

	// GetMessagesFunc mocks the GetMessages method.
	GetMessagesFunc func(ctx context.Context, request *chat.GetMessages) (*chat.GetMessagesResponse, error)

	// GetSubscriptionFunc mocks the GetSubscription method.
	GetSubscriptionFunc func(ctx context.Context, request *chat.GetSubscription) (*model.ChannelSubscription, error)

	// GetSubscriptionsForChannelFunc mocks the GetSubscriptionsForChannel method.
	GetSubscriptionsForChannelFunc func(ctx context.Context, channelKey model.ChannelKey) ([]*model.ChannelSubscription, error)

	// GetSubscriptionsForUserFunc mocks the GetSubscriptionsForUser method.
	GetSubscriptionsForUserFunc func(ctx context.Context, request *chat.GetSubscriptions) (*model.ChannelSubscriptions, error)

	// SaveMessageFunc mocks the SaveMessage method.
	SaveMessageFunc func(ctx context.Context, request *model.Message) error

	// calls tracks calls to the methods.
	calls struct {
		// CreateChannel holds details about calls to the CreateChannel method.
		CreateChannel []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Channel is the channel argument value.
			Channel *model.Channel
		}
		// CreateSubscription holds details about calls to the CreateSubscription method.
		CreateSubscription []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Key is the key argument value.
			Key model.ChannelSubscriptionKey
			// Name is the name argument value.
			Name string
		}
		// DeleteSubscription holds details about calls to the DeleteSubscription method.
		DeleteSubscription []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Key is the key argument value.
			Key model.ChannelSubscriptionKey
		}
		// GetChannel holds details about calls to the GetChannel method.
		GetChannel []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ChannelKey is the channelKey argument value.
			ChannelKey model.ChannelKey
		}
		// GetMessage holds details about calls to the GetMessage method.
		GetMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// MessageKey is the messageKey argument value.
			MessageKey model.MessageKey
		}
		// GetMessages holds details about calls to the GetMessages method.
		GetMessages []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *chat.GetMessages
		}
		// GetSubscription holds details about calls to the GetSubscription method.
		GetSubscription []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *chat.GetSubscription
		}
		// GetSubscriptionsForChannel holds details about calls to the GetSubscriptionsForChannel method.
		GetSubscriptionsForChannel []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ChannelKey is the channelKey argument value.
			ChannelKey model.ChannelKey
		}
		// GetSubscriptionsForUser holds details about calls to the GetSubscriptionsForUser method.
		GetSubscriptionsForUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *chat.GetSubscriptions
		}
		// SaveMessage holds details about calls to the SaveMessage method.
		SaveMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *model.Message
		}
	}
	lockCreateChannel              sync.RWMutex
	lockCreateSubscription         sync.RWMutex
	lockDeleteSubscription         sync.RWMutex
	lockGetChannel                 sync.RWMutex
	lockGetMessage                 sync.RWMutex
	lockGetMessages                sync.RWMutex
	lockGetSubscription            sync.RWMutex
	lockGetSubscriptionsForChannel sync.RWMutex
	lockGetSubscriptionsForUser    sync.RWMutex
	lockSaveMessage                sync.RWMutex
}

// CreateChannel calls CreateChannelFunc.
func (mock *ChatStore) CreateChannel(ctx context.Context, channel *model.Channel) error {
	if mock.CreateChannelFunc == nil {
		panic("ChatStore.CreateChannelFunc: method is nil but Store.CreateChannel was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Channel *model.Channel
	}{
		Ctx:     ctx,
		Channel: channel,
	}
	mock.lockCreateChannel.Lock()
	mock.calls.CreateChannel = append(mock.calls.CreateChannel, callInfo)
	mock.lockCreateChannel.Unlock()
	return mock.CreateChannelFunc(ctx, channel)
}

// CreateChannelCalls gets all the calls that were made to CreateChannel.
// Check the length with:
//     len(mockedStore.CreateChannelCalls())
func (mock *ChatStore) CreateChannelCalls() []struct {
	Ctx     context.Context
	Channel *model.Channel
} {
	var calls []struct {
		Ctx     context.Context
		Channel *model.Channel
	}
	mock.lockCreateChannel.RLock()
	calls = mock.calls.CreateChannel
	mock.lockCreateChannel.RUnlock()
	return calls
}

// CreateSubscription calls CreateSubscriptionFunc.
func (mock *ChatStore) CreateSubscription(ctx context.Context, key model.ChannelSubscriptionKey, name string) (*model.ChannelSubscription, error) {
	if mock.CreateSubscriptionFunc == nil {
		panic("ChatStore.CreateSubscriptionFunc: method is nil but Store.CreateSubscription was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Key  model.ChannelSubscriptionKey
		Name string
	}{
		Ctx:  ctx,
		Key:  key,
		Name: name,
	}
	mock.lockCreateSubscription.Lock()
	mock.calls.CreateSubscription = append(mock.calls.CreateSubscription, callInfo)
	mock.lockCreateSubscription.Unlock()
	return mock.CreateSubscriptionFunc(ctx, key, name)
}

// CreateSubscriptionCalls gets all the calls that were made to CreateSubscription.
// Check the length with:
//     len(mockedStore.CreateSubscriptionCalls())
func (mock *ChatStore) CreateSubscriptionCalls() []struct {
	Ctx  context.Context
	Key  model.ChannelSubscriptionKey
	Name string
} {
	var calls []struct {
		Ctx  context.Context
		Key  model.ChannelSubscriptionKey
		Name string
	}
	mock.lockCreateSubscription.RLock()
	calls = mock.calls.CreateSubscription
	mock.lockCreateSubscription.RUnlock()
	return calls
}

// DeleteSubscription calls DeleteSubscriptionFunc.
func (mock *ChatStore) DeleteSubscription(ctx context.Context, key model.ChannelSubscriptionKey) error {
	if mock.DeleteSubscriptionFunc == nil {
		panic("ChatStore.DeleteSubscriptionFunc: method is nil but Store.DeleteSubscription was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Key model.ChannelSubscriptionKey
	}{
		Ctx: ctx,
		Key: key,
	}
	mock.lockDeleteSubscription.Lock()
	mock.calls.DeleteSubscription = append(mock.calls.DeleteSubscription, callInfo)
	mock.lockDeleteSubscription.Unlock()
	return mock.DeleteSubscriptionFunc(ctx, key)
}

// DeleteSubscriptionCalls gets all the calls that were made to DeleteSubscription.
// Check the length with:
//     len(mockedStore.DeleteSubscriptionCalls())
func (mock *ChatStore) DeleteSubscriptionCalls() []struct {
	Ctx context.Context
	Key model.ChannelSubscriptionKey
} {
	var calls []struct {
		Ctx context.Context
		Key model.ChannelSubscriptionKey
	}
	mock.lockDeleteSubscription.RLock()
	calls = mock.calls.DeleteSubscription
	mock.lockDeleteSubscription.RUnlock()
	return calls
}

// GetChannel calls GetChannelFunc.
func (mock *ChatStore) GetChannel(ctx context.Context, channelKey model.ChannelKey) (*model.Channel, error) {
	if mock.GetChannelFunc == nil {
		panic("ChatStore.GetChannelFunc: method is nil but Store.GetChannel was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		ChannelKey model.ChannelKey
	}{
		Ctx:        ctx,
		ChannelKey: channelKey,
	}
	mock.lockGetChannel.Lock()
	mock.calls.GetChannel = append(mock.calls.GetChannel, callInfo)
	mock.lockGetChannel.Unlock()
	return mock.GetChannelFunc(ctx, channelKey)
}

// GetChannelCalls gets all the calls that were made to GetChannel.
// Check the length with:
//     len(mockedStore.GetChannelCalls())
func (mock *ChatStore) GetChannelCalls() []struct {
	Ctx        context.Context
	ChannelKey model.ChannelKey
} {
	var calls []struct {
		Ctx        context.Context
		ChannelKey model.ChannelKey
	}
	mock.lockGetChannel.RLock()
	calls = mock.calls.GetChannel
	mock.lockGetChannel.RUnlock()
	return calls
}

// GetMessage calls GetMessageFunc.
func (mock *ChatStore) GetMessage(ctx context.Context, messageKey model.MessageKey) (*model.Message, error) {
	if mock.GetMessageFunc == nil {
		panic("ChatStore.GetMessageFunc: method is nil but Store.GetMessage was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		MessageKey model.MessageKey
	}{
		Ctx:        ctx,
		MessageKey: messageKey,
	}
	mock.lockGetMessage.Lock()
	mock.calls.GetMessage = append(mock.calls.GetMessage, callInfo)
	mock.lockGetMessage.Unlock()
	return mock.GetMessageFunc(ctx, messageKey)
}

// GetMessageCalls gets all the calls that were made to GetMessage.
// Check the length with:
//     len(mockedStore.GetMessageCalls())
func (mock *ChatStore) GetMessageCalls() []struct {
	Ctx        context.Context
	MessageKey model.MessageKey
} {
	var calls []struct {
		Ctx        context.Context
		MessageKey model.MessageKey
	}
	mock.lockGetMessage.RLock()
	calls = mock.calls.GetMessage
	mock.lockGetMessage.RUnlock()
	return calls
}

// GetMessages calls GetMessagesFunc.
func (mock *ChatStore) GetMessages(ctx context.Context, request *chat.GetMessages) (*chat.GetMessagesResponse, error) {
	if mock.GetMessagesFunc == nil {
		panic("ChatStore.GetMessagesFunc: method is nil but Store.GetMessages was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *chat.GetMessages
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockGetMessages.Lock()
	mock.calls.GetMessages = append(mock.calls.GetMessages, callInfo)
	mock.lockGetMessages.Unlock()
	return mock.GetMessagesFunc(ctx, request)
}

// GetMessagesCalls gets all the calls that were made to GetMessages.
// Check the length with:
//     len(mockedStore.GetMessagesCalls())
func (mock *ChatStore) GetMessagesCalls() []struct {
	Ctx     context.Context
	Request *chat.GetMessages
} {
	var calls []struct {
		Ctx     context.Context
		Request *chat.GetMessages
	}
	mock.lockGetMessages.RLock()
	calls = mock.calls.GetMessages
	mock.lockGetMessages.RUnlock()
	return calls
}

// GetSubscription calls GetSubscriptionFunc.
func (mock *ChatStore) GetSubscription(ctx context.Context, request *chat.GetSubscription) (*model.ChannelSubscription, error) {
	if mock.GetSubscriptionFunc == nil {
		panic("ChatStore.GetSubscriptionFunc: method is nil but Store.GetSubscription was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *chat.GetSubscription
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockGetSubscription.Lock()
	mock.calls.GetSubscription = append(mock.calls.GetSubscription, callInfo)
	mock.lockGetSubscription.Unlock()
	return mock.GetSubscriptionFunc(ctx, request)
}

// GetSubscriptionCalls gets all the calls that were made to GetSubscription.
// Check the length with:
//     len(mockedStore.GetSubscriptionCalls())
func (mock *ChatStore) GetSubscriptionCalls() []struct {
	Ctx     context.Context
	Request *chat.GetSubscription
} {
	var calls []struct {
		Ctx     context.Context
		Request *chat.GetSubscription
	}
	mock.lockGetSubscription.RLock()
	calls = mock.calls.GetSubscription
	mock.lockGetSubscription.RUnlock()
	return calls
}

// GetSubscriptionsForChannel calls GetSubscriptionsForChannelFunc.
func (mock *ChatStore) GetSubscriptionsForChannel(ctx context.Context, channelKey model.ChannelKey) ([]*model.ChannelSubscription, error) {
	if mock.GetSubscriptionsForChannelFunc == nil {
		panic("ChatStore.GetSubscriptionsForChannelFunc: method is nil but Store.GetSubscriptionsForChannel was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		ChannelKey model.ChannelKey
	}{
		Ctx:        ctx,
		ChannelKey: channelKey,
	}
	mock.lockGetSubscriptionsForChannel.Lock()
	mock.calls.GetSubscriptionsForChannel = append(mock.calls.GetSubscriptionsForChannel, callInfo)
	mock.lockGetSubscriptionsForChannel.Unlock()
	return mock.GetSubscriptionsForChannelFunc(ctx, channelKey)
}

// GetSubscriptionsForChannelCalls gets all the calls that were made to GetSubscriptionsForChannel.
// Check the length with:
//     len(mockedStore.GetSubscriptionsForChannelCalls())
func (mock *ChatStore) GetSubscriptionsForChannelCalls() []struct {
	Ctx        context.Context
	ChannelKey model.ChannelKey
} {
	var calls []struct {
		Ctx        context.Context
		ChannelKey model.ChannelKey
	}
	mock.lockGetSubscriptionsForChannel.RLock()
	calls = mock.calls.GetSubscriptionsForChannel
	mock.lockGetSubscriptionsForChannel.RUnlock()
	return calls
}

// GetSubscriptionsForUser calls GetSubscriptionsForUserFunc.
func (mock *ChatStore) GetSubscriptionsForUser(ctx context.Context, request *chat.GetSubscriptions) (*model.ChannelSubscriptions, error) {
	if mock.GetSubscriptionsForUserFunc == nil {
		panic("ChatStore.GetSubscriptionsForUserFunc: method is nil but Store.GetSubscriptionsForUser was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *chat.GetSubscriptions
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockGetSubscriptionsForUser.Lock()
	mock.calls.GetSubscriptionsForUser = append(mock.calls.GetSubscriptionsForUser, callInfo)
	mock.lockGetSubscriptionsForUser.Unlock()
	return mock.GetSubscriptionsForUserFunc(ctx, request)
}

// GetSubscriptionsForUserCalls gets all the calls that were made to GetSubscriptionsForUser.
// Check the length with:
//     len(mockedStore.GetSubscriptionsForUserCalls())
func (mock *ChatStore) GetSubscriptionsForUserCalls() []struct {
	Ctx     context.Context
	Request *chat.GetSubscriptions
} {
	var calls []struct {
		Ctx     context.Context
		Request *chat.GetSubscriptions
	}
	mock.lockGetSubscriptionsForUser.RLock()
	calls = mock.calls.GetSubscriptionsForUser
	mock.lockGetSubscriptionsForUser.RUnlock()
	return calls
}

// SaveMessage calls SaveMessageFunc.
func (mock *ChatStore) SaveMessage(ctx context.Context, request *model.Message) error {
	if mock.SaveMessageFunc == nil {
		panic("ChatStore.SaveMessageFunc: method is nil but Store.SaveMessage was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *model.Message
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockSaveMessage.Lock()
	mock.calls.SaveMessage = append(mock.calls.SaveMessage, callInfo)
	mock.lockSaveMessage.Unlock()
	return mock.SaveMessageFunc(ctx, request)
}

// SaveMessageCalls gets all the calls that were made to SaveMessage.
// Check the length with:
//     len(mockedStore.SaveMessageCalls())
func (mock *ChatStore) SaveMessageCalls() []struct {
	Ctx     context.Context
	Request *model.Message
} {
	var calls []struct {
		Ctx     context.Context
		Request *model.Message
	}
	mock.lockSaveMessage.RLock()
	calls = mock.calls.SaveMessage
	mock.lockSaveMessage.RUnlock()
	return calls
}
