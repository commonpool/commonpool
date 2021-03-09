package commands

import (
	"github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/commands"
)

// ChangeUserInfoPayload command to change the user info
type ChangeUserInfoPayload struct {
	UserInfo domain.UserInfo `json:"user_info"`
}

type ChangeUserInfo struct {
	commands.CommandEnvelope
	ChangeUserInfoPayload `json:"payload"`
}

var _ commands.Command = &ChangeUserInfo{}

func (c ChangeUserInfo) GetPayload() interface{} {
	return c.ChangeUserInfoPayload
}
