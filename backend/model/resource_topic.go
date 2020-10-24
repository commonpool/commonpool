package model

import uuid "github.com/satori/go.uuid"

type ResourceTopic struct {
	ResourceId uuid.UUID `gorm:"type:uuid;primary_key"`
	UserId     string    `gorm:"primary_key"`
	TopicId    uuid.UUID `gorm:"type:uuid"`
}

func (r *ResourceTopic) GetTopicKey() TopicKey {
	return NewTopicKey(r.TopicId)
}
