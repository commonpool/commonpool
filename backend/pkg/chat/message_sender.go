package chat

import (
	"github.com/commonpool/backend/pkg/keys"
)

type MessageSender struct {
	Type     MessageSenderType
	UserKey  keys.UserKey
	Username string
}
