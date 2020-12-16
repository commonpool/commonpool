package chatmodel

import (
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

type Messages struct {
	Items []Message
}

func NewMessages(items []Message) Messages {
	return Messages{
		Items: items,
	}
}

func (m *Messages) GetAllAuthorKeys() *usermodel.UserKeys {
	var userKeys []usermodel.UserKey
	var userMap = map[usermodel.UserKey]bool{}
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
	return usermodel.NewUserKeys(userKeys)
}
