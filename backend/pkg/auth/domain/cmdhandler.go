package domain

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

type UserCommandHandler struct {
	repo UserRepository
}

func (h *UserCommandHandler) GetName() string {
	return "commandhandler.UserCommandHandler"
}

func NewUserCommandHandler(repo UserRepository) *UserCommandHandler {
	return &UserCommandHandler{
		repo: repo,
	}
}

func (h *UserCommandHandler) GetCommandTypes() []string {
	return []string{
		ChangeUserInfoCmd,
		DiscoverUserCmd,
	}
}

func (h *UserCommandHandler) HandleCommand(ctx context.Context, cmd commands.Command) error {

	switch c := cmd.(type) {
	case *ChangeUserInfo:

		user, err := h.loadUser(ctx, c)
		if err != nil {
			return err
		}

		if err := user.ChangeUserInfo(c.UserInfo); err != nil {
			return err
		}

		return h.repo.Save(ctx, user)

	case *DiscoverUser:
		user, err := h.loadUser(ctx, c)
		if err != nil {
			return err
		}

		if err := user.DiscoverUser(c.UserInfo); err != nil {
			return err
		}

		return h.repo.Save(ctx, user)
	default:
		return fmt.Errorf("unexpected command type: %s", cmd.GetCommandType())
	}

}

func (h *UserCommandHandler) loadUser(ctx context.Context, c commands.Command) (*User, error) {
	user, err := h.repo.Load(ctx, keys.NewUserKey(c.GetAggregateID()))
	if err != nil {
		return nil, err
	}
	return user, nil
}

var _ commands.CommandHandler = &UserCommandHandler{}
