package api

type TargetType string

const (
	UserTarget  TargetType = "user"
	GroupTarget TargetType = "group"
)

type Target struct {
	UserID  *string
	User    *User
	GroupID *string
	Group   *Group
	Type    TargetType
}

func (t *Target) HTMLLink() string {
	if t.Type == UserTarget {
		return t.User.HTMLLink()
	} else if t.Type == GroupTarget {
		return t.Group.HTMLLink()
	} else {
		return ""
	}
}

func (t *Target) IsGroup() bool {
	return t.Type == GroupTarget
}

func (t *Target) IsUser() bool {
	return t.Type == UserTarget
}

func (t *Target) GetGroupID() string {
	return *t.GroupID
}

func (t *Target) GetUserID() string {
	return *t.UserID
}

func (t *Target) GetUser() *User {
	return t.User
}

func (t *Target) getGroup() *Group {
	return t.Group
}
