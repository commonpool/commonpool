package chat

import "github.com/commonpool/backend/model"

type Store interface {
	GetLatestThreads(key model.UserKey, take int64, skip int64, r []*model.Thread) error
	GetThreadMessages(key model.ThreadKey, r []*model.Message) error
}
