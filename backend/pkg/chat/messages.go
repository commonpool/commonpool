package chat

import (
	"github.com/commonpool/backend/pkg/keys"
)

type Messages struct {
	Items []Message
}

func NewMessages(items []Message) Messages {
	return Messages{
		Items: items,
	}
}

func (m *Messages) GetAllAuthorKeys() *keys.UserKeys {
	var userKeys []keys.UserKey
	var userMap = map[keys.UserKey]bool{}
	for _, item := range m.Items {
		if item.MessageType != NormalMessage || item.MessageSubType != UserMessage {
			continue
		}
		authorKey := item.SentBy.UserKey
		if _, ok := userMap[authorKey]; !ok {
			userKeys = append(userKeys, authorKey)
			userMap[authorKey] = true
		}
	}
	return keys.NewUserKeys(userKeys)
}
