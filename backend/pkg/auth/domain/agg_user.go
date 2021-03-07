package domain

import (
	"fmt"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type User struct {
	key     keys.UserKey
	version int
	changes []eventsource.Event
	isNew   bool
	info    UserInfo
}

func NewUser(key keys.UserKey) *User {
	return &User{
		key:     key,
		version: 0,
		changes: []eventsource.Event{},
		isNew:   true,
		info:    UserInfo{},
	}
}

func NewFromEvents(key keys.UserKey, events []eventsource.Event) *User {
	user := NewUser(key)
	for _, event := range events {
		user.on(event, false)
	}
	return user
}

func (u *User) DiscoverUser(userInfo UserInfo) error {
	if err := u.assertIsNew(); err != nil {
		return err
	}
	u.raise(NewUserDiscovered(userInfo))
	return nil
}

func (u *User) handleUserDiscovered(e UserDiscovered) {
	u.info = e.UserInfo
	u.isNew = false
}

func (u *User) ChangeUserInfo(userInfo UserInfo) error {
	if err := u.assertIsNotNew(); err != nil {
		return err
	}
	if userInfo == u.info {
		return nil
	}
	u.raise(NewUserInfoChanged(u.info, userInfo))
	return nil
}

func (u *User) handleUserInfoChanged(e UserInfoChanged) {
	u.info = e.NewUserInfo
}

func (u *User) GetKey() keys.UserKey {
	return u.key
}

func (u *User) GetVersion() int {
	return u.version
}

func (u *User) GetChanges() []eventsource.Event {
	return u.changes
}

func (u *User) assertIsNew() error {
	if !u.isNew {
		return fmt.Errorf("user is new")
	}
	return nil
}

func (u *User) assertIsNotNew() error {
	if u.isNew {
		return fmt.Errorf("user is not new")
	}
	return nil
}

func (u *User) MarkAsCommitted() {
	u.version = u.version + len(u.changes)
	u.changes = []eventsource.Event{}
}

func (u *User) StreamKey() keys.StreamKey {
	return u.key.StreamKey()
}

func (u *User) raise(event eventsource.Event) {
	u.changes = append(u.changes, event)
	u.on(event, true)
}

func (u *User) on(evt eventsource.Event, isNew bool) {
	switch e := evt.(type) {
	case UserDiscovered:
		u.handleUserDiscovered(e)
	case UserInfoChanged:
		u.handleUserInfoChanged(e)
	}
	if !isNew {
		u.version++
	}
}
