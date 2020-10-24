package model

import uuid "github.com/satori/go.uuid"

type TopicKey struct {
	ID uuid.UUID
}

func (tk *TopicKey) String() string {
	return tk.ID.String()
}

func NewTopicKey(id uuid.UUID) TopicKey {
	return TopicKey{
		ID: id,
	}
}
