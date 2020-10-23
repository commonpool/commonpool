package model

import uuid "github.com/satori/go.uuid"

type ResourceTopicKey struct {
	ResourceId uuid.UUID
	TopicId    uuid.UUID
}
