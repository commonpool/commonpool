package model

import (
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"time"
)

type HistoryEntry struct {
	Timestamp         time.Time
	FromUserID        usermodel.UserKey
	ToUserID          usermodel.UserKey
	ResourceID        *resourcemodel.ResourceKey
	TimeAmountSeconds *int64
}
