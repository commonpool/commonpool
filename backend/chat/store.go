package chat

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
)

type Store interface {
	GetLatestThreads(key model.UserKey, take int, skip int) ([]model.Thread, error)
	GetThreadMessages(key model.ThreadKey, take int, skip int) ([]model.Message, error)
	SendMessage(sendMessageRequest *SendMessageRequest) *SendMessageResponse
	SendMessageToThread(sendMessageRequest *SendMessageToThreadRequest) *SendMessageToThreadResponse
	GetOrCreateResourceTopicMapping(rk model.ResourceKey, uk model.UserKey, rs resource.Store) (*model.ResourceTopic, error)
	GetThread(threadKey model.ThreadKey) (*model.Thread, error)
	GetTopic(key model.TopicKey) (*model.Topic, error)
	GetOrCreateConversationTopic(request *GetOrCreateConversationTopicRequest) *GetOrCreateConversationTopicResponse
}

type GetOrCreateConversationTopicRequest struct {
	ParticipantList model.UserKeys
}

type GetOrCreateConversationTopicResponse struct {
	ParticipantList model.UserKeys
	TopicKey        model.TopicKey
	Error           error
}

func NewGetOrCreateConversationTopicRequest(participantList model.UserKeys) GetOrCreateConversationTopicRequest {
	return GetOrCreateConversationTopicRequest{
		ParticipantList: participantList,
	}
}

type SendMessageRequest struct {
	TopicKey     model.TopicKey
	Text         string
	Attachments  []model.Attachment
	Blocks       []model.Block
	FromUser     model.UserKey
	FromUserName string
}

type SendMessageToThreadRequest struct {
	ThreadKey    model.ThreadKey
	Text         string
	Attachments  []model.Attachment
	Blocks       []model.Block
	FromUser     model.UserKey
	FromUserName string
}

type SendMessageResponse struct {
	Error error
}

type SendMessageToThreadResponse struct {
	Error error
}

func NewSendMessageRequest(
	topicKey model.TopicKey,
	fromUser model.UserKey,
	fromUserName string,
	text string,
	blocks []model.Block,
	attachments []model.Attachment,
) SendMessageRequest {
	return SendMessageRequest{
		TopicKey:     topicKey,
		Text:         text,
		Attachments:  attachments,
		Blocks:       blocks,
		FromUser:     fromUser,
		FromUserName: fromUserName,
	}
}

func NewSendMessageToThreadRequest(
	threadKey model.ThreadKey,
	fromUser model.UserKey,
	fromUserName string,
	text string,
	blocks []model.Block,
	attachments []model.Attachment,
) SendMessageToThreadRequest {
	return SendMessageToThreadRequest{
		ThreadKey:    threadKey,
		Text:         text,
		Attachments:  attachments,
		Blocks:       blocks,
		FromUser:     fromUser,
		FromUserName: fromUserName,
	}
}
