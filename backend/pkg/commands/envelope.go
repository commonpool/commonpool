package commands

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
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

type CommandResponse struct {
	Error      error
	StatusCode int
	Payload    interface{}
}

type CommandResponseEnvelope struct {
	CommandEnvelope
	ResponseDuration time.Duration
	ResponseTime     time.Time
	Error            error
	StatusCode       int
	Response         interface{}
}
