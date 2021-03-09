package commands

import "context"

type CommandHandler interface {
	GetName() string
	GetCommandTypes() []string
	HandleCommand(ctx context.Context, cmd Command) CommandResponse
}
