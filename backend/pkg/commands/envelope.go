package commands

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/exceptions"
	uuid "github.com/satori/go.uuid"
	"time"
)

type CommandEnvelope struct {
	CommandTime   time.Time `json:"command_time"`
	CommandType   string    `json:"command_type"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	CorrelationID string    `json:"correlation_id"`
	CommandID     string    `json:"command_id"`
	User          string    `json:"user"`
}

func (c CommandEnvelope) GetCorrelationID() string {
	return c.CorrelationID
}

func (c CommandEnvelope) GetCommandID() string {
	return c.CommandID
}

func (c CommandEnvelope) GetCommandTime() time.Time {
	return c.CommandTime
}

func (c CommandEnvelope) GetCommandType() string {
	return c.CommandType
}

func (c CommandEnvelope) GetAggregateID() string {
	return c.AggregateID
}

func (c CommandEnvelope) GetAggregateType() string {
	return c.AggregateType
}

type NewCommandEnvelopeOptions struct {
	CommandTime *time.Time
	CommandID   *string
}

func NewCommandEnvelope(ctx context.Context, commandType string, aggregateType string, aggregateId string, options ...NewCommandEnvelopeOptions) CommandEnvelope {

	correlationID := uuid.NewV4().String()
	correlationIDFromCtx := ctx.Value("correlationID")
	if correlationIDFromCtx != nil {
		if correlationIDStr, ok := correlationIDFromCtx.(string); ok {
			correlationID = correlationIDStr
		}
	}

	commandTime := time.Now().UTC()
	commandID := uuid.NewV4().String()
	user := ""

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err == nil {
		user = loggedInUser.Subject
	}

	if len(options) > 0 {
		option := options[0]
		if option.CommandTime != nil {
			commandTime = option.CommandTime.UTC()
		}
		if option.CommandID != nil && *option.CommandID != "" {
			commandID = *option.CommandID
		}
	}

	return CommandEnvelope{
		CommandTime:   commandTime,
		CommandType:   commandType,
		AggregateID:   aggregateId,
		AggregateType: aggregateType,
		CorrelationID: correlationID,
		CommandID:     commandID,
		User:          user,
	}

}

type Command interface {
	GetCommandTime() time.Time
	GetCommandType() string
	GetAggregateID() string
	GetAggregateType() string
	GetCorrelationID() string
	GetCommandID() string
	GetPayload() interface{}
}

type CommandResponse interface {
	GetResponse() interface{}
	GetCommand() Command
	GetStatusCode() int
	GetError() error
}

type CommandResponseEnvelope struct {
	Command    Command
	Error      error
	StatusCode int
	Response   interface{}
}

func NewCommandResponseEnvelope(ctx context.Context, command Command, statusCode int, response interface{}, error error) CommandResponseEnvelope {
	return CommandResponseEnvelope{
		Command:    command,
		Error:      error,
		StatusCode: statusCode,
		Response:   response,
	}
}

func NewCommandSuccessResponse(ctx context.Context, command Command, statusCode int, response interface{}) CommandResponseEnvelope {
	return NewCommandResponseEnvelope(ctx, command, statusCode, response, nil)
}

func NewCommandErrResponse(ctx context.Context, command Command, statusCode int, error error) CommandResponseEnvelope {
	return NewCommandResponseEnvelope(ctx, command, statusCode, nil, error)
}

func NewCommandErrFrom(ctx context.Context, command Command, err error) CommandResponseEnvelope {
	return NewCommandErrResponse(ctx, command, exceptions.GetStatusCode(err), err)
}

func (c CommandResponseEnvelope) GetResponse() interface{} {
	return c.Response
}

func (c CommandResponseEnvelope) GetCommand() Command {
	return c.Command
}

func (c CommandResponseEnvelope) GetError() error {
	return c.Error
}

func (c CommandResponseEnvelope) GetStatusCode() int {
	return c.StatusCode
}

var _ CommandResponse = &CommandResponseEnvelope{}
