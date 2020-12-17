package group

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type Group struct {
	Key         keys.GroupKey
	CreatedBy   keys.UserKey
	CreatedAt   time.Time
	Name        string
	Description string
}

func (o *Group) GetKey() keys.GroupKey {
	return o.Key
}

func (o *Group) GetCreatedByKey() keys.UserKey {
	return o.CreatedBy
}
