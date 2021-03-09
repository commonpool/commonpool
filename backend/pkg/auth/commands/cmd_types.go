package commands

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/commands"
)

const (
	ChangeUserInfoCmd = "change_user_info"
	DiscoverUserCmd   = "discover_user"
)

func RegisterCommands(mapper *commands.CommandMapper) {
	mapper.RegisterMapper(DiscoverUserCmd, MapCommand)
	mapper.RegisterMapper(ChangeUserInfoCmd, MapCommand)
}

func MapCommand(commandType string, bytes []byte) (commands.Command, error) {
	var dest commands.Command
	switch commandType {
	case DiscoverUserCmd:
		dest = DiscoverUser{}
	case ChangeUserInfoCmd:
		dest = ChangeUserInfo{}
	default:
		return nil, fmt.Errorf("invalid command type")
	}
	err := json.Unmarshal(bytes, &dest)
	return dest, err
}
