package model

import (
	usermodel "github.com/commonpool/backend/pkg/user/model"
)

type MessageSender struct {
	Type     MessageSenderType
	UserKey  usermodel.UserKey
	Username string
}
