package model

import (
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"time"
)

type Group struct {
	Key         GroupKey
	CreatedBy   usermodel.UserKey
	CreatedAt   time.Time
	Name        string
	Description string
}

func (o *Group) GetKey() GroupKey {
	return o.Key
}

func (o *Group) GetCreatedByKey() usermodel.UserKey {
	return o.CreatedBy
}
