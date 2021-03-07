package keys

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

type StreamKeyer interface {
	StreamKey() StreamKey
}
