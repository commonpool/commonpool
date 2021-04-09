package api

import (
	"fmt"
	"time"
)

type Group struct {
	ID           string
	Name         string
	Memberships  []*Membership
	Posts        []*Post
	CreatedAt    time.Time
	MyMembership *Membership `gorm:"-"`
}

func (g Group) HTMLLink() string {
	return fmt.Sprintf(`<a href="/groups/%s">%s</a>`, g.ID, g.Name)
}
