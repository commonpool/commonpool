package api

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	gorm.Model
	ID           string
	GroupID      string
	Group        *Group
	AuthorID     string
	Author       *User
	Title        string
	Description  string
	Type         PostType
	ValueFrom    *time.Duration
	ValueTo      *time.Duration
	CreatedAt    time.Time
	DeletedAt    *time.Time
	MessageCount int `gorm:"-"`
	Images       []*Image
}

func (p Post) HTMLLink() string {
	return fmt.Sprintf(`<a href="/groups/%s/posts/%s">%s</a>`, p.GroupID, p.ID, p.Title)
}

type PostType string

func (p PostType) IsOffer() bool {
	return p == OfferPost
}

func (p PostType) IsRequest() bool {
	return p == RequestPost
}

const (
	OfferPost   = "offer"
	RequestPost = "request"
	CommentPost = "comment"
)
