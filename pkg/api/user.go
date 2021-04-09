package api

import (
	"fmt"
	"time"
)

type User struct {
	ID               string
	Username         string
	Email            string
	Memberships      []*Membership
	Name             string
	ContactInfo      string
	About            string
	ProfilePictureID string
	CreatedAt        time.Time
}

func (u User) HTMLLink() string {
	return fmt.Sprintf(`<a href="/users/%s">%s</a>`, u.ID, u.Username)
}
