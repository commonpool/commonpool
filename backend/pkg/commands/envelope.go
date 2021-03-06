package commands

import (
	"time"
)

type CommandEnvelope struct {
	CommandTime   time.Time `json:"command_time"`
	CommandType   string    `json:"command_type"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	CorrelationID string    `json:"correlation_id"`
	CommandID     string    `json:"command_id"`
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

type Command interface {
	GetCommandTime() time.Time
	GetCommandType() string
	GetAggregateID() string
	GetAggregateType() string
	GetCorrelationID() string
	GetCommandID() string
	GetPayload() interface{}
}
