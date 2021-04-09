package api

import "time"

type Image struct {
	ID        string
	PostID    string
	Post      *Post
	GroupID   string
	Group     *Group
	CreatedAt time.Time
}
