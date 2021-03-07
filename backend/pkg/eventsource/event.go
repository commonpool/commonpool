package eventsource

import "time"

type Event interface {
	GetEventType() string
	GetEventTime() time.Time
	GetSequenceNo() int
	GetCorrelationID() string
	GetEventID() string
	GetAggregateID() string
	GetAggregateType() string
	GetEventVersion() int
	// SetEventType(eventType string)
	// SetEventTime(eventTime time.Time)
	// SetSequenceNo(sequenceNo int)
	// SetCorrelationID(correlationID string)
	// SetEventID(eventID string)
	// SetAggregateID(aggregateID string)
	// SetAggregateType(aggregateType string)
	// SetEventVersion(eventVersion int)
}

type EventEnvelope struct {
	EventTime     time.Time `json:"event_time,omitempty"`
	EventType     string    `json:"event_type,omitempty"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	EventID       string    `json:"event_id,omitempty"`
	AggregateID   string    `json:"aggregate_id,omitempty"`
	AggregateType string    `json:"aggregate_type,omitempty"`
	EventVersion  int       `json:"event_version,omitempty"`
	SequenceNo    int       `json:"sequence_no,omitempty"`
}

func NewEventEnvelope(eventType string, version int) EventEnvelope {
	return EventEnvelope{
		EventTime:     time.Now().UTC(),
		EventType:     eventType,
		CorrelationID: "",
		EventID:       "",
		AggregateID:   "",
		AggregateType: "",
		EventVersion:  version,
		SequenceNo:    0,
	}
}

func (c EventEnvelope) GetCorrelationID() string {
	return c.CorrelationID
}

func (c *EventEnvelope) SetCorrelationID(correlationID string) {
	c.CorrelationID = correlationID
}

func (c EventEnvelope) GetEventID() string {
	return c.EventID
}

func (c *EventEnvelope) SetEventID(eventID string) {
	c.EventID = eventID
}

func (c EventEnvelope) GetEventTime() time.Time {
	return c.EventTime
}

func (c *EventEnvelope) SetEventTime(eventTime time.Time) {
	c.EventTime = eventTime
}

func (c EventEnvelope) GetEventType() string {
	return c.EventType
}

func (c *EventEnvelope) SetEventType(eventType string) {
	c.EventType = eventType
}

func (c EventEnvelope) GetAggregateID() string {
	return c.AggregateID
}

func (c *EventEnvelope) SetAggregateID(aggregateID string) {
	c.AggregateID = aggregateID
}

func (c EventEnvelope) GetAggregateType() string {
	return c.AggregateType
}

func (c *EventEnvelope) SetAggregateType(aggregateType string) {
	c.AggregateType = aggregateType
}

func (c EventEnvelope) GetEventVersion() int {
	return c.EventVersion
}

func (c *EventEnvelope) SetEventVersion(eventVersion int) {
	c.EventVersion = eventVersion
}

func (c EventEnvelope) GetSequenceNo() int {
	return c.SequenceNo
}

func (c *EventEnvelope) SetSequenceNo(sequenceNo int) {
	c.SequenceNo = sequenceNo
}

type ChangeGetter interface {
	GetChanges() []Event
}

type RevisionGetter interface {
	GetVersion() int
}
