package group

import (
	"github.com/commonpool/backend/model"
	"time"
)

type Group struct {
	Key         model.GroupKey
	CreatedBy   model.UserKey
	CreatedAt   time.Time
	Name        string
	Description string
}

func (o *Group) GetKey() model.GroupKey {
	return o.Key
}

func (o *Group) GetCreatedByKey() model.UserKey {
	return o.CreatedBy
}
