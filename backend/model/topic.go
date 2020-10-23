package model

import uuid "github.com/satori/go.uuid"

type Topic struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key"`
}

