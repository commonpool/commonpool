package domain

import "github.com/commonpool/backend/pkg/commands"

// DiscoverUserPayload command to discover a new user
type DiscoverUserPayload struct {
	UserInfo UserInfo `json:"user_info"`
}

type DiscoverUser struct {
	commands.CommandEnvelope
	DiscoverUserPayload `json:"payload"`
}

var _ commands.Command = &DiscoverUser{}

func (d DiscoverUser) GetPayload() interface{} {
	return d.DiscoverUserPayload
}
