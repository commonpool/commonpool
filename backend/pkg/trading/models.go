package trading

import (
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"time"
)

type HistoryEntry struct {
	Timestamp         time.Time
	FromUserID        usermodel.UserKey
	ToUserID          usermodel.UserKey
	ResourceID        *resourcemodel.ResourceKey
	TimeAmountSeconds *int64
}
