package model

type Thread struct {
	Id                string `gorm:"type:uuid;primary_key"`
	HasUnreadMessages bool
}
