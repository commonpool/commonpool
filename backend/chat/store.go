package chat

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
)

type Store interface {
	GetLatestThreads(key model.UserKey, take int, skip int) ([]model.Thread, error)
	GetThreadMessages(key model.ThreadKey, take int, skip int) ([]model.Message, error)
	SendMessage(author model.UserKey, authorUserName string, topic model.TopicKey, content string) error
	GetOrCreateResourceTopicMapping(rk model.ResourceKey, uk model.UserKey, rs resource.Store) (*model.ResourceTopic, error)
	GetByKey(threadKey model.ThreadKey) (*model.Thread, error)
}
