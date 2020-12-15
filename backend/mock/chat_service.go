// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/chat"
	model2 "github.com/commonpool/backend/pkg/chat/model"
	"sync"
	"time"
)

// Ensure, that ChatService does implement chat.Service.
// If this is not the case, regenerate this file with moq.
var _ chat.Service = &ChatService{}

// ChatService is a mock implementation of chat.Service.
//
//     func TestSomethingThatUsesService(t *testing.T) {
//
//         // make and configure a mocked chat.Service
//         mockedService := &ChatService{
//             CreateChannelFunc: func(ctx context.Context, channelKey model.ChannelKey, channelType chat.ChannelType) (*chat.Channel, error) {
// 	               panic("mock out the CreateChannel method")
//             },
//             CreateUserExchangeFunc: func(ctx context.Context, userKey model.UserKey) (string, error) {
// 	               panic("mock out the CreateUserExchange method")
//             },
//             DeleteGroupChannelFunc: func(ctx context.Context, request *chat.DeleteGroupChannel) (*chat.DeleteGroupChannelResponse, error) {
// 	               panic("mock out the DeleteGroupChannel method")
//             },
//             GetChannelFunc: func(ctx context.Context, channelKey model.ChannelKey) (*chat.Channel, error) {
// 	               panic("mock out the GetChannel method")
//             },
//             GetGroupLinkFunc: func(groupKey model.GroupKey) string {
// 	               panic("mock out the GetGroupLink method")
//             },
//             GetMessageFunc: func(ctx context.Context, messageKey model.MessageKey) (*chat.Message, error) {
// 	               panic("mock out the GetMessage method")
//             },
//             GetMessagesFunc: func(ctx context.Context, channel model.ChannelKey, before time.Time, take int) (*chat.GetMessagesResponse, error) {
// 	               panic("mock out the GetMessages method")
//             },
//             GetResourceLinkFunc: func(resource model.ResourceKey) string {
// 	               panic("mock out the GetResourceLink method")
//             },
//             GetSubscriptionsForUserFunc: func(ctx context.Context, take int, skip int) (*chat.ChannelSubscriptions, error) {
// 	               panic("mock out the GetSubscriptionsForUser method")
//             },
//             GetUserExchangeNameFunc: func(ctx context.Context, userKey model.UserKey) string {
// 	               panic("mock out the GetUserExchangeName method")
//             },
//             GetUserLinkFunc: func(userKey model.UserKey) string {
// 	               panic("mock out the GetUserLink method")
//             },
//             NotifyUserInterestedAboutResourceFunc: func(ctx context.Context, request *chat.NotifyUserInterestedAboutResource) (*chat.NotifyUserInterestedAboutResourceResponse, error) {
// 	               panic("mock out the NotifyUserInterestedAboutResource method")
//             },
//             SendConversationMessageFunc: func(ctx context.Context, request *chat.SendConversationMessage) (*chat.SendConversationMessageResponse, error) {
// 	               panic("mock out the SendConversationMessage method")
//             },
//             SendGroupMessageFunc: func(ctx context.Context, request *chat.SendGroupMessage) (*chat.SendGroupMessageResponse, error) {
// 	               panic("mock out the SendGroupMessage method")
//             },
//             SendMessageFunc: func(ctx context.Context, message *chat.Message) error {
// 	               panic("mock out the SendMessage method")
//             },
//             SubscribeToChannelFunc: func(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey, name string) (*chat.ChannelSubscription, error) {
// 	               panic("mock out the SubscribeToChannel method")
//             },
//             UnsubscribeFromChannelFunc: func(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey) error {
// 	               panic("mock out the UnsubscribeFromChannel method")
//             },
//         }
//
//         // use mockedService in code that requires chat.Service
//         // and then make assertions.
//
//     }
type ChatService struct {
	// CreateChannelFunc mocks the CreateChannel method.
	CreateChannelFunc func(ctx context.Context, channelKey model2.ChannelKey, channelType model2.ChannelType) (*model2.Channel, error)

	// CreateUserExchangeFunc mocks the CreateUserExchange method.
	CreateUserExchangeFunc func(ctx context.Context, userKey model.UserKey) (string, error)

	// DeleteGroupChannelFunc mocks the DeleteGroupChannel method.
	DeleteGroupChannelFunc func(ctx context.Context, request *chat.DeleteGroupChannel) (*chat.DeleteGroupChannelResponse, error)

	// GetChannelFunc mocks the GetChannel method.
	GetChannelFunc func(ctx context.Context, channelKey model2.ChannelKey) (*model2.Channel, error)

	// GetGroupLinkFunc mocks the GetGroupLink method.
	GetGroupLinkFunc func(groupKey model.GroupKey) string

	// GetMessageFunc mocks the GetMessage method.
	GetMessageFunc func(ctx context.Context, messageKey model.MessageKey) (*model2.Message, error)

	// GetMessagesFunc mocks the GetMessages method.
	GetMessagesFunc func(ctx context.Context, channel model2.ChannelKey, before time.Time, take int) (*chat.GetMessagesResponse, error)

	// GetResourceLinkFunc mocks the GetResourceLink method.
	GetResourceLinkFunc func(resource model.ResourceKey) string

	// GetSubscriptionsForUserFunc mocks the GetSubscriptionsForUser method.
	GetSubscriptionsForUserFunc func(ctx context.Context, take int, skip int) (*model2.ChannelSubscriptions, error)

	// GetUserExchangeNameFunc mocks the GetUserExchangeName method.
	GetUserExchangeNameFunc func(ctx context.Context, userKey model.UserKey) string

	// GetUserLinkFunc mocks the GetUserLink method.
	GetUserLinkFunc func(userKey model.UserKey) string

	// NotifyUserInterestedAboutResourceFunc mocks the NotifyUserInterestedAboutResource method.
	NotifyUserInterestedAboutResourceFunc func(ctx context.Context, request *chat.NotifyUserInterestedAboutResource) (*chat.NotifyUserInterestedAboutResourceResponse, error)

	// SendConversationMessageFunc mocks the SendConversationMessage method.
	SendConversationMessageFunc func(ctx context.Context, request *chat.SendConversationMessage) (*chat.SendConversationMessageResponse, error)

	// SendGroupMessageFunc mocks the SendGroupMessage method.
	SendGroupMessageFunc func(ctx context.Context, request *chat.SendGroupMessage) (*chat.SendGroupMessageResponse, error)

	// SendMessageFunc mocks the SendMessage method.
	SendMessageFunc func(ctx context.Context, message *model2.Message) error

	// SubscribeToChannelFunc mocks the SubscribeToChannel method.
	SubscribeToChannelFunc func(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey, name string) (*model2.ChannelSubscription, error)

	// UnsubscribeFromChannelFunc mocks the UnsubscribeFromChannel method.
	UnsubscribeFromChannelFunc func(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey) error

	// calls tracks calls to the methods.
	calls struct {
		// CreateChannel holds details about calls to the CreateChannel method.
		CreateChannel []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ChannelKey is the channelKey argument value.
			ChannelKey model2.ChannelKey
			// ChannelType is the channelType argument value.
			ChannelType model2.ChannelType
		}
		// CreateUserExchange holds details about calls to the CreateUserExchange method.
		CreateUserExchange []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserKey is the userKey argument value.
			UserKey model.UserKey
		}
		// DeleteGroupChannel holds details about calls to the DeleteGroupChannel method.
		DeleteGroupChannel []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *chat.DeleteGroupChannel
		}
		// GetChannel holds details about calls to the GetChannel method.
		GetChannel []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ChannelKey is the channelKey argument value.
			ChannelKey model2.ChannelKey
		}
		// GetGroupLink holds details about calls to the GetGroupLink method.
		GetGroupLink []struct {
			// GroupKey is the groupKey argument value.
			GroupKey model.GroupKey
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
			// Channel is the channel argument value.
			Channel model2.ChannelKey
			// Before is the before argument value.
			Before time.Time
			// Take is the take argument value.
			Take int
		}
		// GetResourceLink holds details about calls to the GetResourceLink method.
		GetResourceLink []struct {
			// Resource is the resource argument value.
			Resource model.ResourceKey
		}
		// GetSubscriptionsForUser holds details about calls to the GetSubscriptionsForUser method.
		GetSubscriptionsForUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Take is the take argument value.
			Take int
			// Skip is the skip argument value.
			Skip int
		}
		// GetUserExchangeName holds details about calls to the GetUserExchangeName method.
		GetUserExchangeName []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserKey is the userKey argument value.
			UserKey model.UserKey
		}
		// GetUserLink holds details about calls to the GetUserLink method.
		GetUserLink []struct {
			// UserKey is the userKey argument value.
			UserKey model.UserKey
		}
		// NotifyUserInterestedAboutResource holds details about calls to the NotifyUserInterestedAboutResource method.
		NotifyUserInterestedAboutResource []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *chat.NotifyUserInterestedAboutResource
		}
		// SendConversationMessage holds details about calls to the SendConversationMessage method.
		SendConversationMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *chat.SendConversationMessage
		}
		// SendGroupMessage holds details about calls to the SendGroupMessage method.
		SendGroupMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *chat.SendGroupMessage
		}
		// SendMessage holds details about calls to the SendMessage method.
		SendMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Message is the message argument value.
			Message *model2.Message
		}
		// SubscribeToChannel holds details about calls to the SubscribeToChannel method.
		SubscribeToChannel []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ChannelSubscriptionKey is the channelSubscriptionKey argument value.
			ChannelSubscriptionKey model.ChannelSubscriptionKey
			// Name is the name argument value.
			Name string
		}
		// UnsubscribeFromChannel holds details about calls to the UnsubscribeFromChannel method.
		UnsubscribeFromChannel []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ChannelSubscriptionKey is the channelSubscriptionKey argument value.
			ChannelSubscriptionKey model.ChannelSubscriptionKey
		}
	}
	lockCreateChannel                     sync.RWMutex
	lockCreateUserExchange                sync.RWMutex
	lockDeleteGroupChannel                sync.RWMutex
	lockGetChannel                        sync.RWMutex
	lockGetGroupLink                      sync.RWMutex
	lockGetMessage                        sync.RWMutex
	lockGetMessages                       sync.RWMutex
	lockGetResourceLink                   sync.RWMutex
	lockGetSubscriptionsForUser           sync.RWMutex
	lockGetUserExchangeName               sync.RWMutex
	lockGetUserLink                       sync.RWMutex
	lockNotifyUserInterestedAboutResource sync.RWMutex
	lockSendConversationMessage           sync.RWMutex
	lockSendGroupMessage                  sync.RWMutex
	lockSendMessage                       sync.RWMutex
	lockSubscribeToChannel                sync.RWMutex
	lockUnsubscribeFromChannel            sync.RWMutex
}

// CreateChannel calls CreateChannelFunc.
func (mock *ChatService) CreateChannel(ctx context.Context, channelKey model2.ChannelKey, channelType model2.ChannelType) (*model2.Channel, error) {
	if mock.CreateChannelFunc == nil {
		panic("ChatService.CreateChannelFunc: method is nil but Service.CreateChannel was just called")
	}
	callInfo := struct {
		Ctx         context.Context
		ChannelKey  model2.ChannelKey
		ChannelType model2.ChannelType
	}{
		Ctx:         ctx,
		ChannelKey:  channelKey,
		ChannelType: channelType,
	}
	mock.lockCreateChannel.Lock()
	mock.calls.CreateChannel = append(mock.calls.CreateChannel, callInfo)
	mock.lockCreateChannel.Unlock()
	return mock.CreateChannelFunc(ctx, channelKey, channelType)
}

// CreateChannelCalls gets all the calls that were made to CreateChannel.
// Check the length with:
//     len(mockedService.CreateChannelCalls())
func (mock *ChatService) CreateChannelCalls() []struct {
	Ctx         context.Context
	ChannelKey  model2.ChannelKey
	ChannelType model2.ChannelType
} {
	var calls []struct {
		Ctx         context.Context
		ChannelKey  model2.ChannelKey
		ChannelType model2.ChannelType
	}
	mock.lockCreateChannel.RLock()
	calls = mock.calls.CreateChannel
	mock.lockCreateChannel.RUnlock()
	return calls
}

// CreateUserExchange calls CreateUserExchangeFunc.
func (mock *ChatService) CreateUserExchange(ctx context.Context, userKey model.UserKey) (string, error) {
	if mock.CreateUserExchangeFunc == nil {
		panic("ChatService.CreateUserExchangeFunc: method is nil but Service.CreateUserExchange was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		UserKey model.UserKey
	}{
		Ctx:     ctx,
		UserKey: userKey,
	}
	mock.lockCreateUserExchange.Lock()
	mock.calls.CreateUserExchange = append(mock.calls.CreateUserExchange, callInfo)
	mock.lockCreateUserExchange.Unlock()
	return mock.CreateUserExchangeFunc(ctx, userKey)
}

// CreateUserExchangeCalls gets all the calls that were made to CreateUserExchange.
// Check the length with:
//     len(mockedService.CreateUserExchangeCalls())
func (mock *ChatService) CreateUserExchangeCalls() []struct {
	Ctx     context.Context
	UserKey model.UserKey
} {
	var calls []struct {
		Ctx     context.Context
		UserKey model.UserKey
	}
	mock.lockCreateUserExchange.RLock()
	calls = mock.calls.CreateUserExchange
	mock.lockCreateUserExchange.RUnlock()
	return calls
}

// DeleteGroupChannel calls DeleteGroupChannelFunc.
func (mock *ChatService) DeleteGroupChannel(ctx context.Context, request *chat.DeleteGroupChannel) (*chat.DeleteGroupChannelResponse, error) {
	if mock.DeleteGroupChannelFunc == nil {
		panic("ChatService.DeleteGroupChannelFunc: method is nil but Service.DeleteGroupChannel was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *chat.DeleteGroupChannel
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockDeleteGroupChannel.Lock()
	mock.calls.DeleteGroupChannel = append(mock.calls.DeleteGroupChannel, callInfo)
	mock.lockDeleteGroupChannel.Unlock()
	return mock.DeleteGroupChannelFunc(ctx, request)
}

// DeleteGroupChannelCalls gets all the calls that were made to DeleteGroupChannel.
// Check the length with:
//     len(mockedService.DeleteGroupChannelCalls())
func (mock *ChatService) DeleteGroupChannelCalls() []struct {
	Ctx     context.Context
	Request *chat.DeleteGroupChannel
} {
	var calls []struct {
		Ctx     context.Context
		Request *chat.DeleteGroupChannel
	}
	mock.lockDeleteGroupChannel.RLock()
	calls = mock.calls.DeleteGroupChannel
	mock.lockDeleteGroupChannel.RUnlock()
	return calls
}

// GetChannel calls GetChannelFunc.
func (mock *ChatService) GetChannel(ctx context.Context, channelKey model2.ChannelKey) (*model2.Channel, error) {
	if mock.GetChannelFunc == nil {
		panic("ChatService.GetChannelFunc: method is nil but Service.GetChannel was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		ChannelKey model2.ChannelKey
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
//     len(mockedService.GetChannelCalls())
func (mock *ChatService) GetChannelCalls() []struct {
	Ctx        context.Context
	ChannelKey model2.ChannelKey
} {
	var calls []struct {
		Ctx        context.Context
		ChannelKey model2.ChannelKey
	}
	mock.lockGetChannel.RLock()
	calls = mock.calls.GetChannel
	mock.lockGetChannel.RUnlock()
	return calls
}

// GetGroupLink calls GetGroupLinkFunc.
func (mock *ChatService) GetGroupLink(groupKey model.GroupKey) string {
	if mock.GetGroupLinkFunc == nil {
		panic("ChatService.GetGroupLinkFunc: method is nil but Service.GetGroupLink was just called")
	}
	callInfo := struct {
		GroupKey model.GroupKey
	}{
		GroupKey: groupKey,
	}
	mock.lockGetGroupLink.Lock()
	mock.calls.GetGroupLink = append(mock.calls.GetGroupLink, callInfo)
	mock.lockGetGroupLink.Unlock()
	return mock.GetGroupLinkFunc(groupKey)
}

// GetGroupLinkCalls gets all the calls that were made to GetGroupLink.
// Check the length with:
//     len(mockedService.GetGroupLinkCalls())
func (mock *ChatService) GetGroupLinkCalls() []struct {
	GroupKey model.GroupKey
} {
	var calls []struct {
		GroupKey model.GroupKey
	}
	mock.lockGetGroupLink.RLock()
	calls = mock.calls.GetGroupLink
	mock.lockGetGroupLink.RUnlock()
	return calls
}

// GetMessage calls GetMessageFunc.
func (mock *ChatService) GetMessage(ctx context.Context, messageKey model.MessageKey) (*model2.Message, error) {
	if mock.GetMessageFunc == nil {
		panic("ChatService.GetMessageFunc: method is nil but Service.GetMessage was just called")
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
//     len(mockedService.GetMessageCalls())
func (mock *ChatService) GetMessageCalls() []struct {
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
func (mock *ChatService) GetMessages(ctx context.Context, channel model2.ChannelKey, before time.Time, take int) (*chat.GetMessagesResponse, error) {
	if mock.GetMessagesFunc == nil {
		panic("ChatService.GetMessagesFunc: method is nil but Service.GetMessages was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Channel model2.ChannelKey
		Before  time.Time
		Take    int
	}{
		Ctx:     ctx,
		Channel: channel,
		Before:  before,
		Take:    take,
	}
	mock.lockGetMessages.Lock()
	mock.calls.GetMessages = append(mock.calls.GetMessages, callInfo)
	mock.lockGetMessages.Unlock()
	return mock.GetMessagesFunc(ctx, channel, before, take)
}

// GetMessagesCalls gets all the calls that were made to GetMessages.
// Check the length with:
//     len(mockedService.GetMessagesCalls())
func (mock *ChatService) GetMessagesCalls() []struct {
	Ctx     context.Context
	Channel model2.ChannelKey
	Before  time.Time
	Take    int
} {
	var calls []struct {
		Ctx     context.Context
		Channel model2.ChannelKey
		Before  time.Time
		Take    int
	}
	mock.lockGetMessages.RLock()
	calls = mock.calls.GetMessages
	mock.lockGetMessages.RUnlock()
	return calls
}

// GetResourceLink calls GetResourceLinkFunc.
func (mock *ChatService) GetResourceLink(resource model.ResourceKey) string {
	if mock.GetResourceLinkFunc == nil {
		panic("ChatService.GetResourceLinkFunc: method is nil but Service.GetResourceLink was just called")
	}
	callInfo := struct {
		Resource model.ResourceKey
	}{
		Resource: resource,
	}
	mock.lockGetResourceLink.Lock()
	mock.calls.GetResourceLink = append(mock.calls.GetResourceLink, callInfo)
	mock.lockGetResourceLink.Unlock()
	return mock.GetResourceLinkFunc(resource)
}

// GetResourceLinkCalls gets all the calls that were made to GetResourceLink.
// Check the length with:
//     len(mockedService.GetResourceLinkCalls())
func (mock *ChatService) GetResourceLinkCalls() []struct {
	Resource model.ResourceKey
} {
	var calls []struct {
		Resource model.ResourceKey
	}
	mock.lockGetResourceLink.RLock()
	calls = mock.calls.GetResourceLink
	mock.lockGetResourceLink.RUnlock()
	return calls
}

// GetSubscriptionsForUser calls GetSubscriptionsForUserFunc.
func (mock *ChatService) GetSubscriptionsForUser(ctx context.Context, take int, skip int) (*model2.ChannelSubscriptions, error) {
	if mock.GetSubscriptionsForUserFunc == nil {
		panic("ChatService.GetSubscriptionsForUserFunc: method is nil but Service.GetSubscriptionsForUser was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Take int
		Skip int
	}{
		Ctx:  ctx,
		Take: take,
		Skip: skip,
	}
	mock.lockGetSubscriptionsForUser.Lock()
	mock.calls.GetSubscriptionsForUser = append(mock.calls.GetSubscriptionsForUser, callInfo)
	mock.lockGetSubscriptionsForUser.Unlock()
	return mock.GetSubscriptionsForUserFunc(ctx, take, skip)
}

// GetSubscriptionsForUserCalls gets all the calls that were made to GetSubscriptionsForUser.
// Check the length with:
//     len(mockedService.GetSubscriptionsForUserCalls())
func (mock *ChatService) GetSubscriptionsForUserCalls() []struct {
	Ctx  context.Context
	Take int
	Skip int
} {
	var calls []struct {
		Ctx  context.Context
		Take int
		Skip int
	}
	mock.lockGetSubscriptionsForUser.RLock()
	calls = mock.calls.GetSubscriptionsForUser
	mock.lockGetSubscriptionsForUser.RUnlock()
	return calls
}

// GetUserExchangeName calls GetUserExchangeNameFunc.
func (mock *ChatService) GetUserExchangeName(ctx context.Context, userKey model.UserKey) string {
	if mock.GetUserExchangeNameFunc == nil {
		panic("ChatService.GetUserExchangeNameFunc: method is nil but Service.GetUserExchangeName was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		UserKey model.UserKey
	}{
		Ctx:     ctx,
		UserKey: userKey,
	}
	mock.lockGetUserExchangeName.Lock()
	mock.calls.GetUserExchangeName = append(mock.calls.GetUserExchangeName, callInfo)
	mock.lockGetUserExchangeName.Unlock()
	return mock.GetUserExchangeNameFunc(ctx, userKey)
}

// GetUserExchangeNameCalls gets all the calls that were made to GetUserExchangeName.
// Check the length with:
//     len(mockedService.GetUserExchangeNameCalls())
func (mock *ChatService) GetUserExchangeNameCalls() []struct {
	Ctx     context.Context
	UserKey model.UserKey
} {
	var calls []struct {
		Ctx     context.Context
		UserKey model.UserKey
	}
	mock.lockGetUserExchangeName.RLock()
	calls = mock.calls.GetUserExchangeName
	mock.lockGetUserExchangeName.RUnlock()
	return calls
}

// GetUserLink calls GetUserLinkFunc.
func (mock *ChatService) GetUserLink(userKey model.UserKey) string {
	if mock.GetUserLinkFunc == nil {
		panic("ChatService.GetUserLinkFunc: method is nil but Service.GetUserLink was just called")
	}
	callInfo := struct {
		UserKey model.UserKey
	}{
		UserKey: userKey,
	}
	mock.lockGetUserLink.Lock()
	mock.calls.GetUserLink = append(mock.calls.GetUserLink, callInfo)
	mock.lockGetUserLink.Unlock()
	return mock.GetUserLinkFunc(userKey)
}

// GetUserLinkCalls gets all the calls that were made to GetUserLink.
// Check the length with:
//     len(mockedService.GetUserLinkCalls())
func (mock *ChatService) GetUserLinkCalls() []struct {
	UserKey model.UserKey
} {
	var calls []struct {
		UserKey model.UserKey
	}
	mock.lockGetUserLink.RLock()
	calls = mock.calls.GetUserLink
	mock.lockGetUserLink.RUnlock()
	return calls
}

// NotifyUserInterestedAboutResource calls NotifyUserInterestedAboutResourceFunc.
func (mock *ChatService) NotifyUserInterestedAboutResource(ctx context.Context, request *chat.NotifyUserInterestedAboutResource) (*chat.NotifyUserInterestedAboutResourceResponse, error) {
	if mock.NotifyUserInterestedAboutResourceFunc == nil {
		panic("ChatService.NotifyUserInterestedAboutResourceFunc: method is nil but Service.NotifyUserInterestedAboutResource was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *chat.NotifyUserInterestedAboutResource
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockNotifyUserInterestedAboutResource.Lock()
	mock.calls.NotifyUserInterestedAboutResource = append(mock.calls.NotifyUserInterestedAboutResource, callInfo)
	mock.lockNotifyUserInterestedAboutResource.Unlock()
	return mock.NotifyUserInterestedAboutResourceFunc(ctx, request)
}

// NotifyUserInterestedAboutResourceCalls gets all the calls that were made to NotifyUserInterestedAboutResource.
// Check the length with:
//     len(mockedService.NotifyUserInterestedAboutResourceCalls())
func (mock *ChatService) NotifyUserInterestedAboutResourceCalls() []struct {
	Ctx     context.Context
	Request *chat.NotifyUserInterestedAboutResource
} {
	var calls []struct {
		Ctx     context.Context
		Request *chat.NotifyUserInterestedAboutResource
	}
	mock.lockNotifyUserInterestedAboutResource.RLock()
	calls = mock.calls.NotifyUserInterestedAboutResource
	mock.lockNotifyUserInterestedAboutResource.RUnlock()
	return calls
}

// SendConversationMessage calls SendConversationMessageFunc.
func (mock *ChatService) SendConversationMessage(ctx context.Context, request *chat.SendConversationMessage) (*chat.SendConversationMessageResponse, error) {
	if mock.SendConversationMessageFunc == nil {
		panic("ChatService.SendConversationMessageFunc: method is nil but Service.SendConversationMessage was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *chat.SendConversationMessage
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockSendConversationMessage.Lock()
	mock.calls.SendConversationMessage = append(mock.calls.SendConversationMessage, callInfo)
	mock.lockSendConversationMessage.Unlock()
	return mock.SendConversationMessageFunc(ctx, request)
}

// SendConversationMessageCalls gets all the calls that were made to SendConversationMessage.
// Check the length with:
//     len(mockedService.SendConversationMessageCalls())
func (mock *ChatService) SendConversationMessageCalls() []struct {
	Ctx     context.Context
	Request *chat.SendConversationMessage
} {
	var calls []struct {
		Ctx     context.Context
		Request *chat.SendConversationMessage
	}
	mock.lockSendConversationMessage.RLock()
	calls = mock.calls.SendConversationMessage
	mock.lockSendConversationMessage.RUnlock()
	return calls
}

// SendGroupMessage calls SendGroupMessageFunc.
func (mock *ChatService) SendGroupMessage(ctx context.Context, request *chat.SendGroupMessage) (*chat.SendGroupMessageResponse, error) {
	if mock.SendGroupMessageFunc == nil {
		panic("ChatService.SendGroupMessageFunc: method is nil but Service.SendGroupMessage was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *chat.SendGroupMessage
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockSendGroupMessage.Lock()
	mock.calls.SendGroupMessage = append(mock.calls.SendGroupMessage, callInfo)
	mock.lockSendGroupMessage.Unlock()
	return mock.SendGroupMessageFunc(ctx, request)
}

// SendGroupMessageCalls gets all the calls that were made to SendGroupMessage.
// Check the length with:
//     len(mockedService.SendGroupMessageCalls())
func (mock *ChatService) SendGroupMessageCalls() []struct {
	Ctx     context.Context
	Request *chat.SendGroupMessage
} {
	var calls []struct {
		Ctx     context.Context
		Request *chat.SendGroupMessage
	}
	mock.lockSendGroupMessage.RLock()
	calls = mock.calls.SendGroupMessage
	mock.lockSendGroupMessage.RUnlock()
	return calls
}

// SendMessage calls SendMessageFunc.
func (mock *ChatService) SendMessage(ctx context.Context, message *model2.Message) error {
	if mock.SendMessageFunc == nil {
		panic("ChatService.SendMessageFunc: method is nil but Service.SendMessage was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Message *model2.Message
	}{
		Ctx:     ctx,
		Message: message,
	}
	mock.lockSendMessage.Lock()
	mock.calls.SendMessage = append(mock.calls.SendMessage, callInfo)
	mock.lockSendMessage.Unlock()
	return mock.SendMessageFunc(ctx, message)
}

// SendMessageCalls gets all the calls that were made to SendMessage.
// Check the length with:
//     len(mockedService.SendMessageCalls())
func (mock *ChatService) SendMessageCalls() []struct {
	Ctx     context.Context
	Message *model2.Message
} {
	var calls []struct {
		Ctx     context.Context
		Message *model2.Message
	}
	mock.lockSendMessage.RLock()
	calls = mock.calls.SendMessage
	mock.lockSendMessage.RUnlock()
	return calls
}

// SubscribeToChannel calls SubscribeToChannelFunc.
func (mock *ChatService) SubscribeToChannel(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey, name string) (*model2.ChannelSubscription, error) {
	if mock.SubscribeToChannelFunc == nil {
		panic("ChatService.SubscribeToChannelFunc: method is nil but Service.SubscribeToChannel was just called")
	}
	callInfo := struct {
		Ctx                    context.Context
		ChannelSubscriptionKey model.ChannelSubscriptionKey
		Name                   string
	}{
		Ctx:                    ctx,
		ChannelSubscriptionKey: channelSubscriptionKey,
		Name:                   name,
	}
	mock.lockSubscribeToChannel.Lock()
	mock.calls.SubscribeToChannel = append(mock.calls.SubscribeToChannel, callInfo)
	mock.lockSubscribeToChannel.Unlock()
	return mock.SubscribeToChannelFunc(ctx, channelSubscriptionKey, name)
}

// SubscribeToChannelCalls gets all the calls that were made to SubscribeToChannel.
// Check the length with:
//     len(mockedService.SubscribeToChannelCalls())
func (mock *ChatService) SubscribeToChannelCalls() []struct {
	Ctx                    context.Context
	ChannelSubscriptionKey model.ChannelSubscriptionKey
	Name                   string
} {
	var calls []struct {
		Ctx                    context.Context
		ChannelSubscriptionKey model.ChannelSubscriptionKey
		Name                   string
	}
	mock.lockSubscribeToChannel.RLock()
	calls = mock.calls.SubscribeToChannel
	mock.lockSubscribeToChannel.RUnlock()
	return calls
}

// UnsubscribeFromChannel calls UnsubscribeFromChannelFunc.
func (mock *ChatService) UnsubscribeFromChannel(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey) error {
	if mock.UnsubscribeFromChannelFunc == nil {
		panic("ChatService.UnsubscribeFromChannelFunc: method is nil but Service.UnsubscribeFromChannel was just called")
	}
	callInfo := struct {
		Ctx                    context.Context
		ChannelSubscriptionKey model.ChannelSubscriptionKey
	}{
		Ctx:                    ctx,
		ChannelSubscriptionKey: channelSubscriptionKey,
	}
	mock.lockUnsubscribeFromChannel.Lock()
	mock.calls.UnsubscribeFromChannel = append(mock.calls.UnsubscribeFromChannel, callInfo)
	mock.lockUnsubscribeFromChannel.Unlock()
	return mock.UnsubscribeFromChannelFunc(ctx, channelSubscriptionKey)
}

// UnsubscribeFromChannelCalls gets all the calls that were made to UnsubscribeFromChannel.
// Check the length with:
//     len(mockedService.UnsubscribeFromChannelCalls())
func (mock *ChatService) UnsubscribeFromChannelCalls() []struct {
	Ctx                    context.Context
	ChannelSubscriptionKey model.ChannelSubscriptionKey
} {
	var calls []struct {
		Ctx                    context.Context
		ChannelSubscriptionKey model.ChannelSubscriptionKey
	}
	mock.lockUnsubscribeFromChannel.RLock()
	calls = mock.calls.UnsubscribeFromChannel
	mock.lockUnsubscribeFromChannel.RUnlock()
	return calls
}
