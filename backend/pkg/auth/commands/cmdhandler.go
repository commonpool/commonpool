package commands

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
	"net/http"
)

type UserCommandHandler struct {
	repo domain.UserRepository
}

func (h *UserCommandHandler) GetName() string {
	return "commandhandler.UserCommandHandler"
}

func NewUserCommandHandler(repo domain.UserRepository) *UserCommandHandler {
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

func (h *UserCommandHandler) HandleCommand(ctx context.Context, cmd commands.Command) commands.CommandResponse {
	switch c := cmd.(type) {
	case *ChangeUserInfo:
		user, err := h.loadUser(ctx, c)
		if err != nil {
			return commands.NewCommandErrFrom(ctx, cmd, err)
		}
		if err := user.ChangeUserInfo(c.UserInfo); err != nil {
			return commands.NewCommandErrFrom(ctx, cmd, err)
		}
		if err := h.repo.Save(ctx, user); err != nil {
			return commands.NewCommandErrFrom(ctx, cmd, err)
		}
	case *DiscoverUser:
		user, err := h.loadUser(ctx, c)
		if err != nil {
			return commands.NewCommandErrFrom(ctx, cmd, err)
		}
		if err := user.DiscoverUser(c.UserInfo); err != nil {
			return commands.NewCommandErrFrom(ctx, cmd, err)
		}
		if err := h.repo.Save(ctx, user); err != nil {
			return commands.NewCommandErrFrom(ctx, cmd, err)
		}

	default:
		err := fmt.Errorf("unexpected command type: %s", cmd.GetCommandType())
		return commands.NewCommandErrFrom(ctx, cmd, err)
	}
	return commands.NewCommandSuccessResponse(ctx, cmd, http.StatusOK, nil)
}

func (h *UserCommandHandler) loadUser(ctx context.Context, c commands.Command) (*domain.User, error) {
	user, err := h.repo.Load(ctx, keys.NewUserKey(c.GetAggregateID()))
	if err != nil {
		return nil, err
	}
	return user, nil
}

var _ commands.CommandHandler = &UserCommandHandler{}
