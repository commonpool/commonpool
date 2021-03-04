package eventstore

import (
	"context"
	"time"
)

type ReplayEventsByTypeOptions struct {
	BatchSize int
}

type EventStore interface {
	Load(ctx context.Context, streamKey StreamKey) ([]*StreamEvent, error)
	Save(ctx context.Context, streamKey StreamKey, expectedRevision int, events []*StreamEvent) error
	ReplayEventsByType(ctx context.Context, eventTypes []string, timestamp time.Time, replayFunc func(events []*StreamEvent) error, options ...ReplayEventsByTypeOptions) error
}

type StreamEventKey struct {
	EventID   string `gorm:"not null;type:varchar(128);primaryKey"`
	EventType string `gorm:"not null;type:varchar(128);primaryKey"`
}

func NewStreamEventKey(eventType string, eventID string) StreamEventKey {
	return StreamEventKey{
		EventID:   eventID,
		EventType: eventType,
	}
}

type StreamEvent struct {
	SequenceNo    int       `gorm:"index:idx_stream_event,sort:asc,priority:3;not null;primaryKey;check:sequence_no >= 0"`
	EventTime     time.Time `gorm:"not null"`
	CorrelationID string    `gorm:"not null;type:varchar(128)"`
	Payload       string    `gorm:"not null;type:jsonb"`
	StreamID      string    `gorm:"index:idx_stream_event,priority:1;not null;type:varchar(128);primaryKey"`
	StreamType    string    `gorm:"index:idx_stream_event,priority:2;not null;type:varchar(128);primaryKey"`
	EventID       string    `gorm:"not null;type:varchar(128);primaryKey"`
	EventType     string    `gorm:"not null;type:varchar(128)"`
	EventVersion  int       `gorm:"not null"`
}

func (s *StreamEvent) StreamKey() StreamKey {
	return StreamKey{
		StreamID:   s.StreamID,
		StreamType: s.StreamType,
	}
}

func (s *StreamEvent) StreamEventKey() StreamEventKey {
	return StreamEventKey{
		EventID:   s.EventID,
		EventType: s.EventType,
	}
}

type NewStreamEventOptions struct {
	EventTime     time.Time
	CorrelationID string
	Version       int
}

func NewStreamEvent(streamKey StreamKey, streamEventKey StreamEventKey, payload string, options ...NewStreamEventOptions) *StreamEvent {
	streamEvent := &StreamEvent{
		EventID:      streamEventKey.EventID,
		EventType:    streamEventKey.EventType,
		Payload:      payload,
		StreamID:     streamKey.StreamID,
		StreamType:   streamKey.StreamType,
		EventVersion: 1,
	}
	if len(options) > 0 {
		streamEvent.EventTime = options[0].EventTime
		streamEvent.CorrelationID = options[0].CorrelationID
		if options[0].Version > 1 {
			streamEvent.EventVersion = options[0].Version
		}
	}
	return streamEvent
}

type Stream struct {
	StreamID      string `gorm:"index:;not null;type:varchar(128);primaryKey"`
	StreamType    string `gorm:"index:;not null;type:varchar(128);primaryKey"`
	LatestVersion int
}

func (s *Stream) StreamKey() StreamKey {
	return StreamKey{
		StreamID:   s.StreamID,
		StreamType: s.StreamType,
	}
}

type StreamKey struct {
	StreamID   string
	StreamType string
}

func NewStreamKey(streamType string, id string) StreamKey {
	return StreamKey{
		StreamID:   id,
		StreamType: streamType,
	}
}
