package messages

import (
	"cp/pkg/api"
	"gorm.io/gorm"
)

type Store interface {
	SendMessage(message *api.Message) error
	GetMessages(threadID string) ([]*api.Message, error)
	DeleteThread(threadID string) error
	FindUserIdsInThread(threadID string) ([]string, error)
}

type MessageStore struct {
	db *gorm.DB
}

func (m MessageStore) SendMessage(message *api.Message) error {
	return m.db.Create(message).Error
}

func (m MessageStore) GetMessages(threadID string) ([]*api.Message, error) {
	var result []*api.Message
	if err := m.db.Preload("Author").Model(&api.Message{}).Find(&result, "thread_id = ?", threadID).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (m MessageStore) DeleteThread(threadID string) error {
	return m.db.Delete(&api.Message{}, "thread_id = ?", threadID).Error
}

func (m MessageStore) FindUserIdsInThread(threadID string) ([]string, error) {
	var result []string
	if err := m.db.Raw("select messages.author_id from messages where thread_id = ? group by messages.author_id", threadID).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func NewMessageStore(db *gorm.DB) *MessageStore {
	return &MessageStore{db: db}
}
