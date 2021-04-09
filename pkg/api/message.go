package api

import "time"

type Message struct {
	ID        string
	AuthorID  string
	Author    *User
	Content   string
	ThreadID  string
	CreatedAt time.Time
}
