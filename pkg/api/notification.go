package api

import "time"

type Notification struct {
	ID        string
	UserID    string
	User      *User
	Title     string
	Message   string
	Link      string
	CreatedAt time.Time
}
