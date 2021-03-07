package eventstore

import (
	"context"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type ReplayEventsByTypeOptions struct {
	BatchSize int
}

type EventStore interface {
	Load(ctx context.Context, streamKey keys.StreamKey) ([]eventsource.Event, error)
	Save(ctx context.Context, streamKey keys.StreamKey, expectedRevision int, events []eventsource.Event) ([]eventsource.Event, error)
	ReplayEventsByType(ctx context.Context, eventTypes []string, timestamp time.Time, replayFunc func(events []eventsource.Event) error, options ...ReplayEventsByTypeOptions) error
}

type StreamEvent struct {
	SequenceNo    int       `gorm:"index:idx_stream_event,sort:asc,priority:3;not null;primaryKey;check:sequence_no >= 0"`
	EventTime     time.Time `gorm:"not null"`
	CorrelationID string    `gorm:"not null;type:varchar(128)"`
	StreamID      string    `gorm:"index:idx_stream_event,priority:1;not null;type:varchar(128);primaryKey"`
	StreamType    string    `gorm:"index:idx_stream_event,priority:2;not null;type:varchar(128);primaryKey"`
	EventID       string    `gorm:"not null;type:varchar(128);primaryKey"`
	EventType     string    `gorm:"not null;type:varchar(128)"`
	EventVersion  int       `gorm:"not null"`
	Body          string    `gorm:"not null;type:jsonb"`
}

func (s *StreamEvent) StreamKey() keys.StreamKey {
	return keys.StreamKey{
		StreamID:   s.StreamID,
		StreamType: s.StreamType,
	}
}

func (s *StreamEvent) StreamEventKey() keys.StreamEventKey {
	return keys.StreamEventKey{
		EventID:   s.EventID,
		EventType: s.EventType,
	}
}

type NewStreamEventOptions struct {
	EventTime     time.Time
	CorrelationID string
	Version       int
}

func NewStreamEvent(streamKey keys.StreamKey, streamEventKey keys.StreamEventKey, payload string, options ...NewStreamEventOptions) *StreamEvent {
	streamEvent := &StreamEvent{
		EventID:      streamEventKey.EventID,
		EventType:    streamEventKey.EventType,
		Body:         payload,
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

func (s *Stream) StreamKey() keys.StreamKey {
	return keys.StreamKey{
		StreamID:   s.StreamID,
		StreamType: s.StreamType,
	}
}
