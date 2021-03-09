package commands

import (
	"context"
	"fmt"
	"net/http"
)

type CommandBus interface {
	RegisterHandler(handler CommandHandler) error
	Send(ctx context.Context, command Command) CommandResponse
}

type LocalCommandBus struct {
	handlers map[string]CommandHandler
}

func NewLocalCommandBus() *LocalCommandBus {
	return &LocalCommandBus{
		handlers: map[string]CommandHandler{},
	}
}

func (l *LocalCommandBus) Send(ctx context.Context, command Command) CommandResponse {
	handler, ok := l.handlers[command.GetCommandType()]
	if !ok {
		return NewCommandErrResponse(ctx, command, http.StatusInternalServerError, fmt.Errorf("no handler registered for command : %s", command.GetCommandType()))
	}
	return handler.HandleCommand(ctx, command)
}

func (l *LocalCommandBus) RegisterHandler(handler CommandHandler) error {
	l.handlers[handler.GetName()] = handler
	return nil
}

var _ CommandBus = &LocalCommandBus{}
