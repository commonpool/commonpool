package chatmodel

import (
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

type MessageSender struct {
	Type     MessageSenderType
	UserKey  usermodel.UserKey
	Username string
}
